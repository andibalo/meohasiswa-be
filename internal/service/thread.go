package service

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/repository"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"net/http"
)

type threadService struct {
	cfg        config.Config
	threadRepo repository.ThreadRepository
	db         *bun.DB
}

func NewThreadService(cfg config.Config, threadRepo repository.ThreadRepository, db *bun.DB) ThreadService {

	return &threadService{
		cfg:        cfg,
		threadRepo: threadRepo,
		db:         db,
	}
}

func (s *threadService) CreateThread(ctx context.Context, req request.CreateThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.CreateThread", "service")
	//defer endFunc()

	thread := &model.Thread{
		ID:             uuid.NewString(),
		UserID:         req.UserID,
		SubThreadID:    req.SubThreadID,
		Title:          req.Title,
		Content:        req.Content,
		ContentSummary: req.ContentSummary,
		IsActive:       true,
		LikeCount:      0,
		DislikeCount:   0,
		CommentCount:   0,
		CreatedBy:      req.UserEmail,
	}

	err := s.threadRepo.Save(thread)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CreateThread] Failed to insert thread to database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to create thread")
	}

	return nil
}

func (s *threadService) GetThreadList(ctx context.Context, req request.GetThreadListReq) (response.GetThreadListResponse, error) {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.GetThreadList", "service")
	//defer endFunc()

	var resp response.GetThreadListResponse

	threads, pagination, err := s.threadRepo.GetList(req)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetThreadList] Failed to get thread list", zap.Error(err))

		return resp, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get thread list")
	}

	resp.Meta = response.PaginationMeta{
		CurrentCursor: pagination.CurrentCursor,
		NextCursor:    pagination.NextCursor,
	}

	resp.Data = threads

	return resp, nil
}
