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
	"github.com/andibalo/meowhasiswa-be/pkg/apperr"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"net/http"
)

type subThreadService struct {
	cfg           config.Config
	subThreadRepo repository.SubThreadRepository
	db            *bun.DB
}

func NewSubThreadService(cfg config.Config, subThreadRepo repository.SubThreadRepository, db *bun.DB) SubThreadService {

	return &subThreadService{
		cfg:           cfg,
		subThreadRepo: subThreadRepo,
		db:            db,
	}
}

func (s *subThreadService) CreateSubThread(ctx context.Context, req request.CreateSubThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "SubThreadService.CreateSubThread", "service")
	//defer endFunc()

	existingUser, err := s.subThreadRepo.GetByName(req.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateSubThread] Failed to get subthread by name", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get subthread by name")
	}

	if existingUser != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateSubThread] Subthread already exists")
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Subthread already exists")
	}

	subThread := &model.SubThread{
		ID:                    uuid.NewString(),
		Name:                  req.Name,
		ImageUrl:              req.ImageUrl,
		Description:           req.Description,
		UniversityID:          req.UniversityID,
		IsUniversitySubThread: req.IsUniversitySubThread,
		CreatedBy:             "SYSTEM", // TODO: update to using token
	}

	err = s.subThreadRepo.Save(subThread)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateSubThread] Failed to insert subthread to database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to create subthread")
	}

	return nil
}

func (s *subThreadService) FollowSubThread(ctx context.Context, req request.FollowSubThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "SubThreadService.FollowSubThread", "service")
	//defer endFunc()

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	subThreadFollower := &model.SubThreadFollower{
		ID:          uuid.NewString(),
		UserID:      req.UserID,
		SubThreadID: req.SubThreadID,
		CreatedBy:   "SYSTEM", // TODO: update to using token
	}

	err = s.subThreadRepo.IncrementFollowersCountTx(req.SubThreadID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to increment subthread followers count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment subthread followers count")
	}

	err = s.subThreadRepo.SaveSubThreadFollowerTx(subThreadFollower, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to save subthread follower", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save subthread follower")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}
