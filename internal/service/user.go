package service

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/repository"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/pkg/integration/notifsvc"
)

type userService struct {
	cfg      config.Config
	notifSvc notifsvc.INotifSvc
	userRepo repository.UserRepository
}

func NewUserService(cfg config.Config, notifSvc notifsvc.INotifSvc, userRepo repository.UserRepository) UserService {

	return &userService{
		cfg:      cfg,
		notifSvc: notifSvc,
		userRepo: userRepo,
	}
}

func (s *userService) GetUserProfile(ctx context.Context, req request.GetUserProfileReq) (*model.User, error) {
	//ctx, endFunc := trace.Start(ctx, "UserService.GetUserProfile", "service")
	//defer endFunc()

	user, err := s.userRepo.GetUserProfileByEmail(req.UserEmail)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) TestCallNotifService(ctx context.Context, req request.TestCallNotifServiceReq) error {
	//ctx, endFunc := trace.Start(ctx, "UserService.TestCallNotifService", "service")
	//defer endFunc()

	_, err := s.notifSvc.CreateNotifTemplate(ctx, notifsvc.CreateNotifTemplateReq{TemplateName: req.TemplateName})
	if err != nil {
		return err
	}

	return nil
}
