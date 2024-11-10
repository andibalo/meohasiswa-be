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

	resp.Data = s.mapThreadListData(threads)

	return resp, nil
}

func (s *threadService) mapThreadListData(threads []model.Thread) []response.ThreadListData {

	threadData := []response.ThreadListData{}

	for _, t := range threads {

		tld := response.ThreadListData{
			ID:             t.ID,
			UserID:         t.UserID,
			UserName:       t.User.Username,
			SubThreadID:    t.SubThreadID,
			SubThreadName:  t.SubThread.Name,
			SubThreadColor: t.SubThread.LabelColor,
			Title:          t.Title,
			Content:        t.Content,
			ContentSummary: t.ContentSummary,
			IsActive:       t.IsActive,
			LikeCount:      t.LikeCount,
			DislikeCount:   t.DislikeCount,
			CommentCount:   t.CommentCount,
			CreatedBy:      t.CreatedBy,
			CreatedAt:      t.CreatedAt,
			UpdatedBy:      t.UpdatedBy,
			UpdatedAt:      t.UpdatedAt,
		}

		if t.User.University != nil {
			tld.UniversityAbbreviatedName = pkg.ToPointer(t.User.University.AbbreviatedName)
			tld.UniversityImageURL = pkg.ToPointer(t.User.University.ImageURL)
		}

		threadData = append(threadData, tld)
	}

	return threadData
}

func (s *threadService) GetThreadDetail(ctx context.Context, req request.GetThreadDetailReq) (response.GetThreadDetailResponse, error) {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.GetThreadDetail", "service")
	//defer endFunc()

	var resp response.GetThreadDetailResponse

	thread, err := s.threadRepo.GetByID(req.ThreadID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[GetThreadDetail] Thread does not exist", zap.Error(err))

			return resp, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Thread does not exist")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[GetThreadDetail] Failed to get thread detail", zap.Error(err))

		return resp, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get thread detail")
	}

	resp.Data = s.mapThreadDetailData(thread)

	return resp, nil
}

func (s *threadService) mapThreadDetailData(thread model.Thread) response.ThreadDetailData {

	threadComments := []response.ThreadComment{}

	td := response.ThreadDetailData{
		ID:             thread.ID,
		UserID:         thread.UserID,
		UserName:       thread.User.Username,
		SubThreadID:    thread.SubThreadID,
		SubThreadName:  thread.SubThread.Name,
		SubThreadColor: thread.SubThread.LabelColor,
		Title:          thread.Title,
		Content:        thread.Content,
		ContentSummary: thread.ContentSummary,
		IsActive:       thread.IsActive,
		LikeCount:      thread.LikeCount,
		DislikeCount:   thread.DislikeCount,
		CommentCount:   thread.CommentCount,
		CreatedBy:      thread.CreatedBy,
		CreatedAt:      thread.CreatedAt,
		UpdatedBy:      thread.UpdatedBy,
		UpdatedAt:      thread.UpdatedAt,
	}

	if thread.User.University != nil {
		td.UniversityAbbreviatedName = pkg.ToPointer(thread.User.University.AbbreviatedName)
		td.UniversityImageURL = pkg.ToPointer(thread.User.University.ImageURL)
	}

	if thread.Comments != nil && len(thread.Comments) > 0 {
		for _, c := range thread.Comments {
			tc := response.ThreadComment{
				ID:           c.ID,
				UserID:       c.UserID,
				UserName:     c.User.Username,
				Content:      c.Content,
				LikeCount:    c.LikeCount,
				DislikeCount: c.DislikeCount,
				CreatedBy:    c.CreatedBy,
				CreatedAt:    c.CreatedAt,
				UpdatedBy:    c.UpdatedBy,
				UpdatedAt:    c.UpdatedAt,
			}

			if c.User.University != nil {
				tc.UniversityAbbreviatedName = pkg.ToPointer(c.User.University.AbbreviatedName)
				tc.UniversityImageURL = pkg.ToPointer(c.User.University.ImageURL)
			}

			threadComments = append(threadComments, tc)
		}

		td.Comments = threadComments
	}

	return td
}

func (s *threadService) LikeThread(ctx context.Context, req request.LikeThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.LikeThread", "service")
	//defer endFunc()

	lastThreadActivity, err := s.getUserLastThreadAction(ctx, req.ThreadID, req.UserID)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to get user last thread action", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	// Unlike if previously liked
	if lastThreadActivity == constants.LIKE_ACTION {
		err = s.unlikeThread(ctx, req, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to unlike thread", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to unlike thread")
		}

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	// Decrement dislike count if previously disliked
	if lastThreadActivity == constants.DISLIKE_ACTION {
		err = s.threadRepo.DecrementDislikesCountTx(req.ThreadID, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to decrement thread dislikes count", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread dislikes count")
		}
	}

	err = s.threadRepo.IncrementLikesCountTx(req.ThreadID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to increment thread likes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment thread likes count")
	}

	threadActivity := &model.ThreadActivity{
		ID:            uuid.NewString(),
		ThreadID:      req.ThreadID,
		ActorID:       req.UserID,
		ActorEmail:    req.UserEmail,
		ActorUsername: req.Username,
		Action:        constants.LIKE_ACTION,
		CreatedBy:     req.UserEmail,
	}

	err = s.threadRepo.SaveThreadActivityTx(threadActivity, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to save thread activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread activity")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}

func (s *threadService) getUserLastThreadAction(ctx context.Context, threadId string, userId string) (string, error) {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.getUserLastThreadAction", "service")
	//defer endFunc()

	ta, err := s.threadRepo.GetLastThreadActivityByUserID(threadId, userId)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[getUserLastThreadAction] Failed to get user last thread activity", zap.Error(err))
		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user last thread activity")
	}

	if ta != nil {

		return ta.Action, nil
	}

	return "", nil
}

func (s *threadService) unlikeThread(ctx context.Context, req request.LikeThreadReq, tx bun.Tx) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.unlikeThread", "service")
	//defer endFunc()

	err := s.threadRepo.DecrementLikesCountTx(req.ThreadID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unlikeThread] Failed to decrement thread likes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread likes count")
	}

	threadActivity := &model.ThreadActivity{
		ID:            uuid.NewString(),
		ThreadID:      req.ThreadID,
		ActorID:       req.UserID,
		ActorEmail:    req.UserEmail,
		ActorUsername: req.Username,
		Action:        constants.UNLIKE_ACTION,
		CreatedBy:     req.UserEmail,
	}

	err = s.threadRepo.SaveThreadActivityTx(threadActivity, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unlikeThread] Failed to save thread activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread activity")
	}

	return nil
}

func (s *threadService) DislikeThread(ctx context.Context, req request.DislikeThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.DislikeThread", "service")
	//defer endFunc()

	lastThreadActivity, err := s.getUserLastThreadAction(ctx, req.ThreadID, req.UserID)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to get user last thread action", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	// Unlike if previously disliked
	if lastThreadActivity == constants.DISLIKE_ACTION {
		err = s.unDislikeThread(ctx, req, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to undislike thread", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to undislike thread")
		}

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	// Decrement like count if previously liked
	if lastThreadActivity == constants.LIKE_ACTION {
		err = s.threadRepo.DecrementLikesCountTx(req.ThreadID, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to decrement thread likes count", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread likes count")
		}
	}

	err = s.threadRepo.IncrementDislikesCountTx(req.ThreadID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to increment thread dislikes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment thread dislikes count")
	}

	threadActivity := &model.ThreadActivity{
		ID:            uuid.NewString(),
		ThreadID:      req.ThreadID,
		ActorID:       req.UserID,
		ActorEmail:    req.UserEmail,
		ActorUsername: req.Username,
		Action:        constants.DISLIKE_ACTION,
		CreatedBy:     req.UserEmail,
	}

	err = s.threadRepo.SaveThreadActivityTx(threadActivity, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to save thread activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread activity")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}

func (s *threadService) unDislikeThread(ctx context.Context, req request.DislikeThreadReq, tx bun.Tx) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.unDislikeThread", "service")
	//defer endFunc()

	err := s.threadRepo.DecrementDislikesCountTx(req.ThreadID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unDislikeThread] Failed to decrement thread dislikes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread dislikes count")
	}

	threadActivity := &model.ThreadActivity{
		ID:            uuid.NewString(),
		ThreadID:      req.ThreadID,
		ActorID:       req.UserID,
		ActorEmail:    req.UserEmail,
		ActorUsername: req.Username,
		Action:        constants.UNDISLIKE_ACTION,
		CreatedBy:     req.UserEmail,
	}

	err = s.threadRepo.SaveThreadActivityTx(threadActivity, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unDislikeThread] Failed to save thread activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread activity")
	}

	return nil
}

func (s *threadService) CommentThread(ctx context.Context, req request.CommentThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.CommentThread", "service")
	//defer endFunc()

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CommentThread] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	err = s.threadRepo.IncrementCommentsCountTx(req.ThreadID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CommentThread] Failed to increment thread comments count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment thread comments count")
	}

	threadComment := &model.ThreadComment{
		ID:        uuid.NewString(),
		ThreadID:  req.ThreadID,
		UserID:    req.UserID,
		Content:   req.Content,
		CreatedBy: req.UserEmail,
	}

	err = s.threadRepo.SaveThreadCommentTx(threadComment, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CommentThread] Failed to save thread comment", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread comment")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CommentThread] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}
