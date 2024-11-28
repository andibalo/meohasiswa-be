package service

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
)

type UserService interface {
	GetUserProfile(ctx context.Context, req request.GetUserProfileReq) (*model.User, error)
	CreateUserDevice(ctx context.Context, req request.CreateUserDeviceReq) error
	GetUserDevices(ctx context.Context, req request.GetUserDevicesReq) ([]model.UserDevice, error)
}

type AuthService interface {
	Register(ctx context.Context, req request.RegisterUserReq) error
	Login(ctx context.Context, req request.LoginUserReq) (token string, err error)
	VerifyEmail(ctx context.Context, req request.VerifyEmailReq) (err error)
}

type SubThreadService interface {
	GetSubThreadList(ctx context.Context, req request.GetSubThreadListReq) (response.GetSubThreadListResponse, error)
	GetSubThreadByID(ctx context.Context, req request.GetSubThreadByIDReq) (response.GetSubThreadByIDResponse, error)
	UpdateSubThread(ctx context.Context, req request.UpdateSubThreadReq) error
	CreateSubThread(ctx context.Context, req request.CreateSubThreadReq) error
	FollowSubThread(ctx context.Context, req request.FollowSubThreadReq) error
	UnFollowSubThread(ctx context.Context, req request.UnFollowSubThreadReq) error
	DeleteSubThread(ctx context.Context, req request.DeleteSubThreadReq) error
}

type ThreadService interface {
	CreateThread(ctx context.Context, req request.CreateThreadReq) error
	UpdateThread(ctx context.Context, req request.UpdateThreadReq) error
	DeleteThread(ctx context.Context, req request.DeleteThreadReq) error
	GetThreadList(ctx context.Context, req request.GetThreadListReq) (response.GetThreadListResponse, error)
	GetThreadDetail(ctx context.Context, req request.GetThreadDetailReq) (response.GetThreadDetailResponse, error)
	LikeThread(ctx context.Context, req request.LikeThreadReq) error
	DislikeThread(ctx context.Context, req request.DislikeThreadReq) error
	CommentThread(ctx context.Context, req request.CommentThreadReq) error
	ReplyComment(ctx context.Context, req request.ReplyCommentReq) error
	LikeComment(ctx context.Context, req request.LikeCommentReq) error
	DislikeComment(ctx context.Context, req request.DislikeCommentReq) error
	GetThreadComments(ctx context.Context, req request.GetThreadCommentsReq) (response.GetThreadCommentsResponse, error)
}

type UniversityService interface {
	GetUniversityRatingList(ctx context.Context, req request.GetUniversityRatingListReq) (response.GetUniversityRatingListResponse, error)
	CreateUniversityRating(ctx context.Context, req request.RateUniversityReq) error
}

type ImageService interface {
	UploadImage(ctx context.Context, fileData model.File) (response.UploadImageResp, error)
}

type NotificationService interface {
	SendPushNotification(ctx context.Context, req request.SendPushNotificationReq) error
}
