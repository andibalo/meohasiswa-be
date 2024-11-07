package service

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/repository"
	"github.com/andibalo/meowhasiswa-be/internal/request"
)

type authService struct {
	cfg      config.Config
	userRepo repository.UserRepository
}

func NewAuthService(cfg config.Config, userRepo repository.UserRepository) AuthService {

	return &authService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req request.RegisterUserReq) error {
	//ctx, endFunc := trace.Start(ctx, "AuthService.Register", "service")
	//defer endFunc()

	return nil
}
