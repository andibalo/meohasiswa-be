package core

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/api"
	v1 "github.com/andibalo/meowhasiswa-be/internal/api/v1"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/middleware"
	"github.com/andibalo/meowhasiswa-be/internal/repository"
	"github.com/andibalo/meowhasiswa-be/internal/service"
	"github.com/andibalo/meowhasiswa-be/pkg/httpclient"
	"github.com/andibalo/meowhasiswa-be/pkg/integration/notifsvc"
	s3Repository "github.com/andibalo/meowhasiswa-be/pkg/s3"
	"github.com/andibalo/meowhasiswa-be/pkg/trace"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"

	"net/http"
)

type Server struct {
	gin *gin.Engine
	srv *http.Server
}

func NewServer(cfg config.Config, tracer *trace.Tracer, db *bun.DB, s3Client *s3.Client) *Server {

	router := gin.New()

	router.Use(middleware.LogPreReq(cfg.Logger()))

	if cfg.GetFlags().EnableTracer {
		tracer.SetGinMiddleware(router, cfg.AppName())

		router.Use(trace.TracerLogger())
	}

	router.Use(cors.Default())
	router.Use(gin.Recovery())

	hc := httpclient.Init(httpclient.Options{Config: cfg})

	s3Repo := s3Repository.NewS3Repository(cfg, s3Client)
	universityRepo := repository.NewUniversityRepository(db)
	subThreadRepo := repository.NewSubThreadRepository(db)
	userRepo := repository.NewUserRepository(db)
	threadRepo := repository.NewThreadRepository(db)

	notifSvc := notifsvc.NewNotificationService(cfg, hc)
	imageSvc := service.NewImageService(cfg, s3Repo)
	universitySvc := service.NewUniversityService(cfg, universityRepo, userRepo, db)
	authSvc := service.NewAuthService(cfg, userRepo, db)
	userSvc := service.NewUserService(cfg, notifSvc, userRepo)
	subThreadSvc := service.NewSubThreadService(cfg, subThreadRepo, db)
	threadSvc := service.NewThreadService(cfg, threadRepo, db)

	ic := v1.NewImageController(cfg, imageSvc)
	uc := v1.NewUserController(cfg, userSvc)
	ac := v1.NewAuthController(cfg, authSvc)
	stc := v1.NewSubThreadController(cfg, subThreadSvc)
	tc := v1.NewThreadController(cfg, threadSvc)
	unc := v1.NewUniversityController(cfg, universitySvc)

	registerHandlers(router, &api.HealthCheck{}, uc, ac, stc, tc, unc, ic)

	return &Server{
		gin: router,
	}
}

func (s *Server) Start(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.gin,
	}

	s.srv = srv

	return srv.ListenAndServe()
}

func (s *Server) GetGin() *gin.Engine {

	return s.gin
}

func (s *Server) Shutdown(ctx context.Context) error {

	return s.srv.Shutdown(ctx)
}

func registerHandlers(g *gin.Engine, handlers ...api.Handler) {
	for _, handler := range handlers {
		handler.AddRoutes(g)
	}
}
