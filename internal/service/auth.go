package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/constants"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/repository"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg"
	"github.com/andibalo/meowhasiswa-be/pkg/apperr"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type authService struct {
	cfg      config.Config
	userRepo repository.UserRepository
	db       *bun.DB
}

func NewAuthService(cfg config.Config, userRepo repository.UserRepository, db *bun.DB) AuthService {

	return &authService{
		cfg:      cfg,
		userRepo: userRepo,
		db:       db,
	}
}

func (s *authService) Register(ctx context.Context, req request.RegisterUserReq) error {
	//ctx, endFunc := trace.Start(ctx, "AuthService.Register", "service")
	//defer endFunc()

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to get user by email", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Invalid Email/Password")
	}

	if existingUser != nil && existingUser.ID != "" {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] User already exists")
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User already exists")
	}

	user, err := s.mapCreateUserReqToUserModel(ctx, req)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to map payload to user model", zap.Error(err))
		return oops.Wrapf(err, "[Register] Failed to map payload to user model")
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	err = s.userRepo.SaveTx(user, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to insert user to database", zap.Error(err))
		tx.Rollback()

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	userVerifyEmail := &model.UserVerifyEmail{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Code:      pkg.GenRandomString(10),
		Email:     user.Email,
		IsUsed:    false,
		ExpiredAt: time.Now().Add(time.Minute * time.Duration(s.cfg.GetAuthCfg().UserSecretCodeExpiryMins)),
		CreatedBy: s.cfg.AppName(),
	}

	err = s.userRepo.SaveUserVerifyEmailTx(userVerifyEmail, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to insert user verify email to database", zap.Error(err))
		tx.Rollback()

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	// SEND EMAIL
	//verifyUrl := s.cfg.AppURL() + constants.UserVerifyEmailPath + fmt.Sprintf("?secret_code=%s&id=%s", userVerifyEmail.SecretCode, userVerifyEmail.ID)
	//
	//msg := pubsubCommons.CoreNewRegisteredUserPayload{
	//	FirstName: data.FirstName,
	//	LastName:  data.LastName,
	//	Email:     data.Email,
	//	VerifyURL: verifyUrl,
	//}
	//
	//go func() {
	//	s.pb.PublishNewUserRegistered(msg)
	//}()

	return nil
}

func (s *authService) mapCreateUserReqToUserModel(ctx context.Context, data request.RegisterUserReq) (*model.User, error) {

	hasedPassword, err := pkg.HashPassword(data.Password)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[mapCreateUserReqToUserModel] Failed to hash password", zap.Error(err))

		return nil, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return &model.User{
		ID:               uuid.NewString(),
		Username:         data.Username,
		Email:            data.Email,
		Role:             constants.USER_ROLE,
		Password:         hasedPassword,
		IsBanned:         false,
		IsEmailVerified:  false,
		ReputationPoints: 50,
		CreatedBy:        data.Email,
	}, nil
}

func (s *authService) Login(ctx context.Context, req request.LoginUserReq) (token string, err error) {
	//ctx, endFunc := trace.Start(ctx, "AuthService.Login", "service")
	//defer endFunc()

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[Login] Invalid email/password", zap.Error(err))
			return "", oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Invalid Email/Password")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[Login] Failed to get user by email", zap.Error(err))
		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	isMatch := pkg.CheckPasswordHash(req.Password, existingUser.Password)
	if !isMatch {
		s.cfg.Logger().ErrorWithContext(ctx, "[Login] Invalid password for user", zap.String("email", req.Email))
		return "", oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Invalid Email/Password")
	}

	token, err = pkg.GenerateToken(existingUser)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Login] Failed to generate JWT Token for user", zap.String("email", req.Email))
		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return token, nil
}
