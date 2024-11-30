package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/repository"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg"
	"github.com/andibalo/meowhasiswa-be/pkg/apperr"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/google/uuid"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type userService struct {
	cfg      config.Config
	userRepo repository.UserRepository
	uniRepo  repository.UniversityRepository
}

func NewUserService(cfg config.Config, userRepo repository.UserRepository, uniRepo repository.UniversityRepository) UserService {

	return &userService{
		cfg:      cfg,
		userRepo: userRepo,
		uniRepo:  uniRepo,
	}
}

func (s *userService) GetUserProfile(ctx context.Context, req request.GetUserProfileReq) (*model.User, error) {
	//ctx, endFunc := trace.Start(ctx, "UserService.GetUserProfile", "service")
	//defer endFunc()

	user, err := s.userRepo.GetUserProfileByEmail(req.UserEmail)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[GetUserProfile] User profile not found", zap.Error(err))
			return nil, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User profile not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[GetUserProfile] Failed to get user profile", zap.Error(err))
		return nil, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user profile")
	}

	if user.UniversityID != nil && user.HasRateUniversity {
		unir, err := s.uniRepo.GetUniversityRatingByUserIDAndUniversityID(req.UserID, *user.UniversityID)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[GetUserProfile] Failed to get user university rating", zap.Error(err))
			return nil, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user university rating")
		}

		user.UniversityRatingID = &unir.ID
	}

	return user, nil
}

func (s *userService) GetUserDevices(ctx context.Context, req request.GetUserDevicesReq) ([]model.UserDevice, error) {
	//ctx, endFunc := trace.Start(ctx, "UserService.GetUserDevices", "service")
	//defer endFunc()

	userDevices, err := s.userRepo.GetUserDevices(req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[GetUserDevices] User devices not found", zap.Error(err))
			return nil, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User devices not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[GetUserDevices] Failed to get user devices", zap.Error(err))
		return nil, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user devices")
	}

	return userDevices, nil
}

func (s *userService) CreateUserDevice(ctx context.Context, req request.CreateUserDeviceReq) error {
	//ctx, endFunc := trace.Start(ctx, "UserService.CreateUserDevice", "service")
	//defer endFunc()

	_, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[CreateUserDevice] User not found", zap.Error(err))
			return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUserDevice] Failed to get user by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	ud := &model.UserDevice{
		ID:                   uuid.NewString(),
		UserID:               req.UserID,
		NotificationToken:    req.NotificationToken,
		IsNotificationActive: true,
		CreatedBy:            req.UserEmail,
	}

	if req.Brand != "" {
		ud.Brand = pkg.ToPointer(req.Brand)
	}

	if req.Type != "" {
		ud.Type = pkg.ToPointer(req.Type)
	}

	if req.Model != "" {
		ud.Model = pkg.ToPointer(req.Model)
	}

	err = s.userRepo.SaveUserDevice(ud)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateUserDevice] Failed to create user device", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to create user device")
	}

	return nil
}
