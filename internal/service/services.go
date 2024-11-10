package service

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
)

type UserService interface {
	TestCallNotifService(ctx context.Context, req request.TestCallNotifServiceReq) error
}

type AuthService interface {
	Register(ctx context.Context, req request.RegisterUserReq) error
	Login(ctx context.Context, req request.LoginUserReq) (token string, err error)
}

type SubThreadService interface {
	GetSubThreadList(ctx context.Context, req request.GetSubThreadListReq) (response.GetSubThreadListResponse, error)
	CreateSubThread(ctx context.Context, req request.CreateSubThreadReq) error
	FollowSubThread(ctx context.Context, req request.FollowSubThreadReq) error
	UnFollowSubThread(ctx context.Context, req request.UnFollowSubThreadReq) error
}

type ThreadService interface {
	CreateThread(ctx context.Context, req request.CreateThreadReq) error
	GetThreadList(ctx context.Context, req request.GetThreadListReq) (response.GetThreadListResponse, error)
	LikeThread(ctx context.Context, req request.LikeThreadReq) error
	DislikeThread(ctx context.Context, req request.DislikeThreadReq) error
}
