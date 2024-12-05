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
	uniRepo   repository.UniversityRepository
	db        *bun.DB
	mailerSvc mailer.MailService
}

func NewAuthService(cfg config.Config, userRepo repository.UserRepository, uniRepo repository.UniversityRepository, db *bun.DB, mailerSvc mailer.MailService) AuthService {

	return &authService{
		cfg:       cfg,
		userRepo:  userRepo,
		uniRepo:   uniRepo,
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

	emailDomain, err := pkg.ExtractDomainFromEmail(req.Email)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to extract domain from email", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Failed to extract domain from email")
	}

	uni, err := s.uniRepo.GetByDomain(emailDomain)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to get university by domain", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get university by domain")
	}

	if uni.ID == "" {
		s.cfg.Logger().WarnWithContext(ctx, "[Register] Could not detect any university from email")
	}

	if uni.ID != "" {
		user.UniversityID = pkg.ToPointer(uni.ID)
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

	userVerifyCode := &model.UserVerifyCode{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Type:      constants.TYPE_VERIFY_EMAIL,
		Code:      pkg.GenRandNumber(6),
		Email:     user.Email,
		IsUsed:    false,
		ExpiredAt: time.Now().Add(time.Minute * time.Duration(s.cfg.GetAuthCfg().UserSecretCodeExpiryMins)),
		CreatedBy: s.cfg.AppName(),
	}

	err = s.userRepo.SaveUserVerifyCodeTx(userVerifyCode, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[Register] Failed to insert user verify code to database", zap.Error(err))
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
				"code": userVerifyCode.Code,
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

	if !existingUser.IsEmailVerified {
		s.cfg.Logger().ErrorWithContext(ctx, "[Login] User email is not yet verified", zap.String("email", req.Email))
		return "", oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User email is not yet verified")
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

	userVerifyCode, err := s.userRepo.GetUserVerifyCodeByEmail(req.Email, constants.TYPE_VERIFY_EMAIL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] User verify code not found", zap.Error(err))
			return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User verify code not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Failed to get user verify code by email", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	if userVerifyCode.ExpiredAt.Before(time.Now()) {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Code has expired", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Code has expired")
	}

	if userVerifyCode.Code != req.Code {
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

	err = s.userRepo.UpdateUserVerifyCodeByIDTx(userVerifyCode.ID, updateValues, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyEmail] Failed to update user verify code", zap.Error(err))
		tx.Rollback()

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update user verify code")
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

func (s *authService) ResetPassword(ctx context.Context, req request.ResetPasswordReq) (err error) {
	//ctx, endFunc := trace.Start(ctx, "AuthService.ResetPassword", "service")
	//defer endFunc()

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[ResetPassword] User not found", zap.Error(err))
			return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[ResetPassword] Failed to get user by email", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	hasedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[ResetPassword] Failed to hash password", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	updateValues := map[string]interface{}{
		"password":   hasedPassword,
		"updated_by": req.Email,
		"updated_at": time.Now(),
	}

	err = s.userRepo.UpdateUserPasswordByUserID(existingUser.ID, updateValues)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[ResetPassword] Failed to update user password", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update user password")
	}

	return nil
}

func (s *authService) VerifyResetPassword(ctx context.Context, req request.VerifyResetPasswordReq) (err error) {
	//ctx, endFunc := trace.Start(ctx, "AuthService.VerifyResetPassword", "service")
	//defer endFunc()

	_, err = s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[VerifyResetPassword] User not found", zap.Error(err))
			return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyResetPassword] Failed to get user by email", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	userVerifyCode, err := s.userRepo.GetUserVerifyCodeByEmail(req.Email, constants.TYPE_RESET_PASSWORD)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[VerifyResetPassword] User verify code not found", zap.Error(err))
			return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User verify code not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyResetPassword] Failed to get user verify code by email", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	if userVerifyCode.ExpiredAt.Before(time.Now()) {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyResetPassword] Code has expired", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Code has expired")
	}

	if userVerifyCode.Code != req.Code {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyResetPassword] Code does not match", zap.Error(err))
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Code does not match")
	}

	updateValues := map[string]interface{}{
		"is_used":    true,
		"updated_by": req.Email,
		"updated_at": time.Now(),
	}

	err = s.userRepo.UpdateUserVerifyCodeByID(userVerifyCode.ID, updateValues)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[VerifyResetPassword] Failed to update user verify code", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update user verify code")
	}

	return nil
}

func (s *authService) SendResetPasswordLink(ctx context.Context, req request.SendResetPasswordLinkReq) (err error) {
	//ctx, endFunc := trace.Start(ctx, "AuthService.SendResetPasswordLink", "service")
	//defer endFunc()

	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[SendResetPasswordLink] User not found", zap.Error(err))
			return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("User not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[SendResetPasswordLink] Failed to get user by email", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	userVerifyCode := &model.UserVerifyCode{
		ID:        uuid.NewString(),
		UserID:    existingUser.ID,
		Type:      constants.TYPE_RESET_PASSWORD,
		Code:      pkg.GenRandNumber(6),
		Email:     existingUser.Email,
		IsUsed:    false,
		ExpiredAt: time.Now().Add(time.Minute * time.Duration(s.cfg.GetAuthCfg().UserSecretCodeExpiryMins)),
		CreatedBy: s.cfg.AppName(),
	}

	err = s.userRepo.SaveUserVerifyCode(userVerifyCode)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[SendResetPasswordLink] Failed to insert user verify code to database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	if s.cfg.GetFlags().EnableSendEmail {
		err = s.mailerSvc.SendMail(ctx, mailer.Mail{
			To: []string{
				req.Email,
			},
			Name:       mailer.SEND_RESET_PASSWORD_EMAIL,
			Subject:    mailer.SEND_RESET_PASSWORD_EMAIL_SUBJECT,
			TemplateID: s.cfg.GetBrevoSvcCfg().SendResetPasswordTemplateId,
			Data: map[string]interface{}{
				"code": userVerifyCode.Code,
			},
		})

		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[SendResetPasswordLink] Failed to send reset password link email", zap.Error(err))
		}
	}

	return nil
}
