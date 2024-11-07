package service

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/request"
)

type UserService interface {
	TestCallNotifService(ctx context.Context, req request.TestCallNotifServiceReq) error
}

type AuthService interface {
	Register(ctx context.Context, req request.RegisterUserReq) error
	Login(ctx context.Context, req request.LoginUserReq) (token string, err error)
}
