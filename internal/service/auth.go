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
	"github.com/andibalo/meowhasiswa-be/pkg/mailer"
	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type authService struct {
	cfg       config.Config
	userRepo  repository.UserRepository
	db        *bun.DB
	mailerSvc mailer.MailService
}

func NewAuthService(cfg config.Config, userRepo repository.UserRepository, db *bun.DB, mailerSvc mailer.MailService) AuthService {

	return &authService{
		cfg:       cfg,
		userRepo:  userRepo,
		db:        db,
		mailerSvc: mailerSvc,
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
		Code:      pkg.GenRandNumber(6),
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

	if s.cfg.GetFlags().EnableSendEmail {
		err = s.mailerSvc.SendMail(ctx, mailer.Mail{
			To: []string{
				req.Email,
			},
			Name:       mailer.SEND_VERIFICATION_CODE_EMAIL,
			Subject:    mailer.SEND_VERIFICATION_CODE_EMAIL_SUBJECT,
			TemplateID: s.cfg.GetBrevoSvcCfg().SendVerificationCodeTemplateId,
			Data: map[string]interface{}{
				"code": userVerifyEmail.Code,
			},
		})

		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to send mail", zap.Error(err))
		}
	}

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

func (s *authService) VerifyEmail(ctx context.Context, req request.VerifyEmailReq) (err error) {
	//ctx, endFunc := trace.Start(ctx, "AuthService.VerifyEmail", "service")
	//defer endFunc()

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] User not found", zap.Error(err))
			return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Failed to get user by email", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	if existingUser.IsEmailVerified {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] User is already verified", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User is already verified")
	}

	userVerifyEmail, err := s.userRepo.GetUserVerifyEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] User verify email not found", zap.Error(err))
			return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User verify email not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Failed to get user verify email by email", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	if userVerifyEmail.ExpiredAt.Before(time.Now()) {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Code has expired", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Code has expired")
	}

	if userVerifyEmail.Code != req.Code {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Code does not match", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Code does not match")
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	updateValues := map[string]interface{}{
		"is_used":    true,
		"updated_by": req.Email,
		"updated_at": time.Now(),
	}

	err = s.userRepo.UpdateUserVerifyEmailByIDTx(userVerifyEmail.ID, updateValues, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Failed to update user verify email", zap.Error(err))
		tx.Rollback()

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update user verify email")
	}

	err = s.userRepo.SetUserVerifyEmailToUsedTx(userVerifyEmail.ID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Failed to set user verify email to used", zap.Error(err))
		tx.Rollback()

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to set user verify email to used")
	}

	err = s.userRepo.SetUserToEmailVerifiedTx(existingUser.ID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Failed to set user to email verified", zap.Error(err))
		tx.Rollback()

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to set user to email verified")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}
