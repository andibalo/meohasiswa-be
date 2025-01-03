package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/andibalo/meowhasiswa-be"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/pkg/db"
	"github.com/andibalo/meowhasiswa-be/pkg/trace"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.InitConfig()

	database := db.InitDB(cfg)

	var tracer *trace.Tracer

	if cfg.GetFlags().EnableTracer {
		tracer = initTracer(cfg)
	}

	awsCreds := credentials.NewStaticCredentialsProvider(cfg.GetAWSCfg().ACCESS_KEY_ID, cfg.GetAWSCfg().SECRET_ACCESS_KEY, "")

	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(cfg.GetAWSCfg().Region), awsConfig.WithCredentialsProvider(awsCreds))
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(awsCfg)

	server := core.NewServer(cfg, tracer, database, client)

	cfg.Logger().Info(fmt.Sprintf("Server starting at port %s", cfg.AppAddress()))

	go func() {
		if err := server.Start(cfg.AppAddress()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cfg.Logger().Fatal("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cfg.Logger().Info("shutting down gracefully, press Ctrl+C again to force")

	if err := tracer.Close(ctx); err != nil {
		log.Fatalln("error shutdown tracer: ", err)
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		cfg.Logger().Fatal("Server force to shutdown")
	}

	_ = database.Close()

	cfg.Logger().Info("Server exiting")
}

func initTracer(cfg config.Config) *trace.Tracer {

	traceConfig := cfg.TraceConfig()

	// init tracer type
	tracer, err := trace.Init(context.Background(), traceConfig)
	if err != nil {
		log.Fatal("error init tracer: ", err)
	}

	return tracer
}
