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
	"time"
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

func (s *subThreadService) GetSubThreadList(ctx context.Context, req request.GetSubThreadListReq) (response.GetSubThreadListResponse, error) {
	//ctx, endFunc := trace.Start(ctx, "SubThreadService.GetSubThreadList", "service")
	//defer endFunc()

	var resp response.GetSubThreadListResponse

	subThreads, pagination, err := s.subThreadRepo.GetList(req)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetSubThreadList] Failed to get subthread list", zap.Error(err))

		return resp, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get subthread list")
	}

	resp.Meta = response.PaginationMeta{
		CurrentCursor: pagination.CurrentCursor,
		NextCursor:    pagination.NextCursor,
	}

	resp.Data = subThreads

	return resp, nil
}

func (s *subThreadService) GetSubThreadByID(ctx context.Context, req request.GetSubThreadByIDReq) (response.GetSubThreadByIDResponse, error) {
	//ctx, endFunc := trace.Start(ctx, "SubThreadService.GetSubThreadByID", "service")
	//defer endFunc()

	var resp response.GetSubThreadByIDResponse

	subthread, err := s.subThreadRepo.GetByID(req.SubThreadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[UpdateSubThread] SubThread not found", zap.Error(err))
			return resp, oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("SubThread not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateSubThread] Failed to get subthread by id", zap.Error(err))
		return resp, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update subthread by id")
	}

	resp.Data = subthread

	return resp, nil
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
		LabelColor:            req.LabelColor,
		UniversityID:          req.UniversityID,
		IsUniversitySubThread: req.IsUniversitySubThread,
		CreatedBy:             req.UserEmail,
	}

	err = s.subThreadRepo.Save(subThread)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateSubThread] Failed to insert subthread to database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to create subthread")
	}

	return nil
}

func (s *subThreadService) UpdateSubThread(ctx context.Context, req request.UpdateSubThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "SubThreadService.UpdateSubThread", "service")
	//defer endFunc()

	_, err := s.subThreadRepo.GetByID(req.SubThreadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[UpdateSubThread] SubThread not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("SubThread not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateSubThread] Failed to get subthread by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update subthread by id")
	}

	updateValues := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"image_url":   req.ImageUrl,
		"label_color": req.LabelColor,
		"updated_by":  req.UserEmail,
		"updated_at":  time.Now(),
	}

	if req.UniversityID != nil && *req.UniversityID != "" {
		updateValues["university_id"] = req.UniversityID
	}

	if req.IsUniversitySubThread != nil {
		updateValues["is_university_subthread"] = req.IsUniversitySubThread
	}

	err = s.subThreadRepo.UpdateByID(req.SubThreadID, updateValues)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateSubThread] Failed to update subthread in database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update subthread")
	}

	return nil
}

func (s *subThreadService) DeleteSubThread(ctx context.Context, req request.DeleteSubThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "SubThreadService.DeleteSubThread", "service")
	//defer endFunc()

	_, err := s.subThreadRepo.GetByID(req.SubThreadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[DeleteSubThread] SubThread not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("SubThread not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteSubThread] Failed to delete subthread by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to delete subthread by id")
	}

	updateValues := map[string]interface{}{
		"deleted_by": req.UserEmail,
		"deleted_at": time.Now(),
	}

	err = s.subThreadRepo.DeleteByID(req.SubThreadID, updateValues)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteSubThread] Failed to delete subthread in database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to delete subthread")
	}

	return nil
}

func (s *subThreadService) FollowSubThread(ctx context.Context, req request.FollowSubThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "SubThreadService.FollowSubThread", "service")
	//defer endFunc()

	stf, err := s.subThreadRepo.GetSubThreadFollowerByUserIDAndSubThreadID(req.UserID, req.SubThreadID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to fetch subthread follower by user and subthread id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	if stf != nil && stf.IsFollowing {
		s.cfg.Logger().WarnWithContext(ctx, "[FollowSubThread] Subthread already followed")
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	subThreadFollower := &model.SubThreadFollower{
		ID:          uuid.NewString(),
		UserID:      req.UserID,
		SubThreadID: req.SubThreadID,
		IsFollowing: true,
		CreatedBy:   "SYSTEM", // TODO: update to using token
	}

	err = s.subThreadRepo.IncrementFollowersCountTx(req.SubThreadID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to increment subthread followers count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment subthread followers count")
	}

	if stf != nil {

		err = s.subThreadRepo.UpdateSubThreadFollowerIsFollowingTx(stf.ID, true, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to update subthread follower", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update subthread follower")
		}

	} else {

		err = s.subThreadRepo.SaveSubThreadFollowerTx(subThreadFollower, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to save subthread follower", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save subthread follower")
		}
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[FollowSubThread] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}

func (s *subThreadService) UnFollowSubThread(ctx context.Context, req request.UnFollowSubThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "SubThreadService.UnFollowSubThread", "service")
	//defer endFunc()

	stf, err := s.subThreadRepo.GetSubThreadFollowerByUserIDAndSubThreadID(req.UserID, req.SubThreadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[UnFollowSubThread] Subthread follower not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Subthread follower not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[UnFollowSubThread] Failed to fetch subthread follower by user and subthread id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	if !stf.IsFollowing {
		s.cfg.Logger().WarnWithContext(ctx, "[UnFollowSubThread] Subthread already unfollowed")
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UnFollowSubThread] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	err = s.subThreadRepo.DecrementFollowersCountTx(req.SubThreadID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UnFollowSubThread] Failed to decrement subthread followers count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement subthread followers count")
	}

	err = s.subThreadRepo.UpdateSubThreadFollowerIsFollowingTx(stf.ID, false, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UnFollowSubThread] Failed to update subthread follower", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update subthread follower")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UnFollowSubThread] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}
