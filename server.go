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
	"github.com/andibalo/meowhasiswa-be/pkg/trace"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"

	"net/http"
)

type Server struct {
	gin *gin.Engine
	srv *http.Server
}

func NewServer(cfg config.Config, tracer *trace.Tracer, db *bun.DB) *Server {

	router := gin.New()

	router.Use(middleware.LogPreReq(cfg.Logger()))

	if cfg.GetFlags().EnableTracer {
		tracer.SetGinMiddleware(router, cfg.AppName())

		router.Use(trace.TracerLogger())
	}

	router.Use(gin.Recovery())

	hc := httpclient.Init(httpclient.Options{Config: cfg})

	userRepo := repository.NewUserRepository(db)

	notifSvc := notifsvc.NewNotificationService(cfg, hc)

	authSvc := service.NewAuthService(cfg, userRepo, db)
	userSvc := service.NewUserService(cfg, notifSvc)

	uc := v1.NewUserController(cfg, userSvc)
	ac := v1.NewAuthController(cfg, authSvc)

	registerHandlers(router, &api.HealthCheck{}, uc, ac)

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
