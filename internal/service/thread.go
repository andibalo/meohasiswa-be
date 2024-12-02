package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/constants"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/repository"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg"
	"github.com/andibalo/meowhasiswa-be/pkg/apperr"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/andibalo/meowhasiswa-be/pkg/integration/notifsvc"
	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type threadService struct {
	cfg        config.Config
	threadRepo repository.ThreadRepository
	userRepo   repository.UserRepository
	notifCl    notifsvc.INotifSvc
	db         *bun.DB
}

func NewThreadService(cfg config.Config, threadRepo repository.ThreadRepository, userRepo repository.UserRepository, notifCl notifsvc.INotifSvc, db *bun.DB) ThreadService {

	return &threadService{
		cfg:        cfg,
		threadRepo: threadRepo,
		userRepo:   userRepo,
		notifCl:    notifCl,
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

func (s *threadService) UpdateThread(ctx context.Context, req request.UpdateThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.UpdateThread", "service")
	//defer endFunc()

	_, err := s.threadRepo.GetByID(req.ThreadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[UpdateThread] Thread not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateThread] Failed to get thread by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread by id")
	}

	updateValues := map[string]interface{}{
		"title":           req.Title,
		"content":         req.Content,
		"content_summary": req.ContentSummary,
		"updated_by":      req.UserEmail,
		"updated_at":      time.Now(),
	}

	err = s.threadRepo.UpdateByID(req.ThreadID, updateValues)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateThread] Failed to update thread in database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread")
	}

	return nil
}

func (s *threadService) DeleteThread(ctx context.Context, req request.DeleteThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.DeleteThread", "service")
	//defer endFunc()

	_, err := s.threadRepo.GetByID(req.ThreadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[DeleteThread] Thread not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteThread] Failed to get thread by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get thread by id")
	}

	updateValues := map[string]interface{}{
		"deleted_by": req.UserEmail,
		"deleted_at": time.Now(),
	}

	err = s.threadRepo.DeleteByID(req.ThreadID, updateValues)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteThread] Failed to delete thread in database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to delete thread")
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

		if t.ThreadAction != "" {
			if t.ThreadAction == constants.LIKE_ACTION {
				tld.IsLiked = true
			}

			if t.ThreadAction == constants.DISLIKE_ACTION {
				tld.IsDisliked = true
			}
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

	ta, err := s.threadRepo.GetLastThreadActivityByUserID(req.ThreadID, req.UserID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetThreadDetail] Failed to get user last thread activity", zap.Error(err))
		return resp, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user last thread activity")
	}

	resp.Data = s.mapThreadDetailData(thread, ta)

	return resp, nil
}

func (s *threadService) mapThreadDetailData(thread model.Thread, threadActivity *model.ThreadActivity) response.ThreadDetailData {

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

	if threadActivity != nil {
		if threadActivity.Action == constants.LIKE_ACTION {
			td.IsLiked = true
		}

		if threadActivity.Action == constants.DISLIKE_ACTION {
			td.IsDisliked = true
		}
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

	if lastThreadActivity != "" {

		updateValues := map[string]interface{}{
			"action":     constants.LIKE_ACTION,
			"updated_by": req.UserEmail,
			"updated_at": time.Now(),
		}

		err = s.threadRepo.UpdateThreadActivityTx(req.ThreadID, req.UserID, updateValues, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to update thread activity", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread activity")
		}

		// TODO: Save to thread activity history

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeThread] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
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

	updateValues := map[string]interface{}{
		"action":     constants.UNLIKE_ACTION,
		"updated_by": req.UserEmail,
		"updated_at": time.Now(),
	}

	err = s.threadRepo.UpdateThreadActivityTx(req.ThreadID, req.UserID, updateValues, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unlikeThread] Failed to update thread activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread activity")
	}

	// TODO: Save to thread activity history

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

	if lastThreadActivity != "" {

		updateValues := map[string]interface{}{
			"action":     constants.DISLIKE_ACTION,
			"updated_by": req.UserEmail,
			"updated_at": time.Now(),
		}

		err = s.threadRepo.UpdateThreadActivityTx(req.ThreadID, req.UserID, updateValues, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to update thread activity", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread activity")
		}

		// TODO: Save to thread activity history

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeThread] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
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

	updateValues := map[string]interface{}{
		"action":     constants.UNDISLIKE_ACTION,
		"updated_by": req.UserEmail,
		"updated_at": time.Now(),
	}

	err = s.threadRepo.UpdateThreadActivityTx(req.ThreadID, req.UserID, updateValues, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unDislikeThread] Failed to update thread activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread activity")
	}

	// TODO: Save to thread activity history

	return nil
}

func (s *threadService) CommentThread(ctx context.Context, req request.CommentThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.CommentThread", "service")
	//defer endFunc()

	thread, err := s.threadRepo.GetByID(req.ThreadID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[CommentThread] Thread does not exist", zap.Error(err))

			return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Thread does not exist")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[CommentThread] Failed to get thread detail", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get thread detail")
	}

	userDevices, err := s.userRepo.GetUserDevices(request.GetUserDevicesReq{
		UserID: thread.UserID,
	})

	if len(userDevices) == 0 {
		s.cfg.Logger().ErrorWithContext(ctx, "[CommentThread] User devices not found", zap.Error(err))
		return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.NotFound).Errorf("User devices not found")
	}

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[CommentThread] Failed to get user devices", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user devices")
	}

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

	if thread.UserID != req.UserID {
		var notificationTokens []string

		for _, ud := range userDevices {
			if ud.IsNotificationActive {
				notificationTokens = append(notificationTokens, ud.NotificationToken)
			}
		}

		if len(notificationTokens) > 0 {

			notifData := map[string]string{
				constants.APP_ROUTE_KEY:  "/thread/" + thread.ID,
				constants.EVENT_TYPE_KEY: constants.COMMENT_ON_THREAD_EVENT,
			}

			_, err = s.notifCl.SendPushNotification(ctx, notifsvc.SendPushNotificationReq{
				NotificationTokens: notificationTokens,
				Title:              "Someone commented on your thread!",
				Content:            fmt.Sprintf("%s: %s", req.Username, pkg.TruncateWithEllipsis(req.Content, 50)),
				Data:               notifData,
			})

			if err != nil {
				s.cfg.Logger().ErrorWithContext(ctx, "[CommentThread] Failed to send push notification", zap.Error(err))
			}
		}
	}

	return nil
}

func (s *threadService) ReplyComment(ctx context.Context, req request.ReplyCommentReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.ReplyComment", "service")
	//defer endFunc()

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[ReplyComment] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	err = s.threadRepo.IncrementCommentReplyCountTx(req.CommentID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[ReplyComment] Failed to increment comment reply count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment comment reply count")
	}

	threadCommentReply := &model.ThreadCommentReply{
		ID:              uuid.NewString(),
		ThreadID:        req.ThreadID,
		UserID:          req.UserID,
		ThreadCommentID: req.CommentID,
		Content:         req.Content,
		CreatedBy:       req.UserEmail,
	}

	err = s.threadRepo.SaveCommentReplyTx(threadCommentReply, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[ReplyComment] Failed to save thread comment reply", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread comment reply")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[ReplyComment] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}

func (s *threadService) LikeComment(ctx context.Context, req request.LikeCommentReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.LikeComment", "service")
	//defer endFunc()

	// Handle is comment reply
	if req.IsReply {
		err := s.likeCommentReply(ctx, req)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to like comment reply", zap.Error(err))
			return err
		}

		return nil
	}

	lastThreadCommentActivity, err := s.getUserLastThreadCommentAction(ctx, req.ThreadID, req.CommentID, req.UserID)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to get user last thread comment action", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	// Unlike if previously liked
	if lastThreadCommentActivity == constants.LIKE_ACTION {
		err = s.unlikeComment(ctx, req, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to unlike comment", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to unlike comment")
		}

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	// Decrement dislike count if previously disliked
	if lastThreadCommentActivity == constants.DISLIKE_ACTION {
		err = s.threadRepo.DecrementCommentDislikesCountTx(req.CommentID, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to decrement comment dislikes count", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement comment dislikes count")
		}
	}

	err = s.threadRepo.IncrementCommentLikesCountTx(req.CommentID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to increment comment likes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment comment likes count")
	}

	if lastThreadCommentActivity != "" {

		updateValues := map[string]interface{}{
			"action":     constants.LIKE_ACTION,
			"updated_by": req.UserEmail,
			"updated_at": time.Now(),
		}

		err = s.threadRepo.UpdateThreadCommentActivityTx(req.CommentID, req.UserID, updateValues, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to update thread comment activity", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment activity")
		}

		// TODO: Save to thread activity history

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	threadActivity := &model.ThreadCommentActivity{
		ID:              uuid.NewString(),
		ThreadID:        req.ThreadID,
		ThreadCommentID: req.CommentID,
		ActorID:         req.UserID,
		ActorEmail:      req.UserEmail,
		ActorUsername:   req.Username,
		Action:          constants.LIKE_ACTION,
		CreatedBy:       req.UserEmail,
	}

	err = s.threadRepo.SaveThreadCommentActivityTx(threadActivity, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to save thread comment activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread comment activity")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[LikeComment] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}

func (s *threadService) likeCommentReply(ctx context.Context, req request.LikeCommentReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.likeCommentReply", "service")
	//defer endFunc()

	lastThreadCommentReplyActivity, err := s.getUserLastThreadCommentReplyAction(ctx, req.ThreadID, req.CommentID, req.UserID)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to get user last thread comment reply action", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	threadCommentReply, err := s.threadRepo.GetThreadCommentReplyByID(req.CommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Thread comment reply does not exist", zap.Error(err))

			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread comment reply does not exist")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to get thread comment reply by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get thread comment reply by id")
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	// Unlike if previously liked
	if lastThreadCommentReplyActivity == constants.LIKE_ACTION {
		err = s.unlikeCommentReply(ctx, req, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to unlike comment reply", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to unlike comment reply")
		}

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	// Decrement dislike count if previously disliked
	if lastThreadCommentReplyActivity == constants.DISLIKE_ACTION {
		err = s.threadRepo.DecrementCommentReplyDislikesCountTx(req.CommentID, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to decrement comment reply dislikes count", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement comment reply dislikes count")
		}
	}

	err = s.threadRepo.IncrementCommentReplyLikesCountTx(req.CommentID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to increment comment reply likes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment comment reply likes count")
	}

	if lastThreadCommentReplyActivity != "" {

		updateValues := map[string]interface{}{
			"action":     constants.LIKE_ACTION,
			"updated_by": req.UserEmail,
			"updated_at": time.Now(),
		}

		err = s.threadRepo.UpdateThreadCommentActivityReplyTx(req.CommentID, req.UserID, updateValues, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to update thread comment reply activity", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment reply activity")
		}

		// TODO: Save to thread activity history

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	threadActivity := &model.ThreadCommentActivity{
		ID:                   uuid.NewString(),
		ThreadID:             req.ThreadID,
		ThreadCommentID:      threadCommentReply.ThreadCommentID,
		ThreadCommentReplyID: pkg.ToPointer(req.CommentID),
		ActorID:              req.UserID,
		ActorEmail:           req.UserEmail,
		ActorUsername:        req.Username,
		Action:               constants.LIKE_ACTION,
		CreatedBy:            req.UserEmail,
	}

	err = s.threadRepo.SaveThreadCommentActivityTx(threadActivity, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to save thread comment reply activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread comment reply activity")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[likeCommentReply] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	// TODO: Save to thread comment activity history

	return nil
}

func (s *threadService) unlikeComment(ctx context.Context, req request.LikeCommentReq, tx bun.Tx) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.unlikeComment", "service")
	//defer endFunc()

	err := s.threadRepo.DecrementCommentLikesCountTx(req.CommentID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unlikeComment] Failed to decrement thread comment likes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread comment likes count")
	}

	updateValues := map[string]interface{}{
		"action":     constants.UNLIKE_ACTION,
		"updated_by": req.UserEmail,
		"updated_at": time.Now(),
	}

	err = s.threadRepo.UpdateThreadCommentActivityTx(req.CommentID, req.UserID, updateValues, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unlikeComment] Failed to update thread comment activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment activity")
	}

	// TODO: Save to thread comment activity history

	return nil
}

func (s *threadService) unlikeCommentReply(ctx context.Context, req request.LikeCommentReq, tx bun.Tx) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.unlikeCommentReply", "service")
	//defer endFunc()

	err := s.threadRepo.DecrementCommentReplyLikesCountTx(req.CommentID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unlikeCommentReply] Failed to decrement thread comment reply likes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread comment reply likes count")
	}

	updateValues := map[string]interface{}{
		"action":     constants.UNLIKE_ACTION,
		"updated_by": req.UserEmail,
		"updated_at": time.Now(),
	}

	err = s.threadRepo.UpdateThreadCommentActivityReplyTx(req.CommentID, req.UserID, updateValues, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unlikeCommentReply] Failed to update thread comment reply activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment reply activity")
	}

	// TODO: Save to thread comment activity history

	return nil
}

func (s *threadService) getUserLastThreadCommentAction(ctx context.Context, threadId string, threadCommentId string, userId string) (string, error) {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.getUserLastThreadCommentAction", "service")
	//defer endFunc()

	tca, err := s.threadRepo.GetLastThreadCommentActivityByUserID(threadId, threadCommentId, userId)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[getUserLastThreadCommentAction] Failed to get user last thread comment activity", zap.Error(err))
		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user last thread comment activity")
	}

	if tca != nil {
		return tca.Action, nil
	}

	return "", nil
}

func (s *threadService) getUserLastThreadCommentReplyAction(ctx context.Context, threadId string, threadCommentReplyId string, userId string) (string, error) {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.getUserLastThreadCommentReplyAction", "service")
	//defer endFunc()

	tca, err := s.threadRepo.GetLastThreadCommentActivityReplyByUserID(threadId, threadCommentReplyId, userId)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[getUserLastThreadCommentAction] Failed to get user last thread comment reply activity", zap.Error(err))
		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get user last thread comment reply activity")
	}

	if tca != nil {
		return tca.Action, nil
	}

	return "", nil
}

func (s *threadService) DislikeComment(ctx context.Context, req request.DislikeCommentReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.DislikeComment", "service")
	//defer endFunc()

	// Handle is comment reply
	if req.IsReply {
		err := s.dislikeCommentReply(ctx, req)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to dislike comment reply", zap.Error(err))
			return err
		}

		return nil
	}

	lastThreadCommentActivity, err := s.getUserLastThreadCommentAction(ctx, req.ThreadID, req.CommentID, req.UserID)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to get user last thread comment action", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	// Unlike if previously disliked
	if lastThreadCommentActivity == constants.DISLIKE_ACTION {
		err = s.unDislikeComment(ctx, req, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to undislike comment", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to undislike comment")
		}

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	// Decrement like count if previously liked
	if lastThreadCommentActivity == constants.LIKE_ACTION {
		err = s.threadRepo.DecrementCommentLikesCountTx(req.CommentID, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to decrement thread comment likes count", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread comment likes count")
		}
	}

	err = s.threadRepo.IncrementCommentDislikesCountTx(req.CommentID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to increment thread comment dislikes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment thread comment dislikes count")
	}

	if lastThreadCommentActivity != "" {

		updateValues := map[string]interface{}{
			"action":     constants.DISLIKE_ACTION,
			"updated_by": req.UserEmail,
			"updated_at": time.Now(),
		}

		err = s.threadRepo.UpdateThreadCommentActivityTx(req.CommentID, req.UserID, updateValues, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to update thread comment activity", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment activity")
		}

		// TODO: Save to thread activity history

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	threadCommentActivity := &model.ThreadCommentActivity{
		ID:              uuid.NewString(),
		ThreadID:        req.ThreadID,
		ThreadCommentID: req.CommentID,
		ActorID:         req.UserID,
		ActorEmail:      req.UserEmail,
		ActorUsername:   req.Username,
		Action:          constants.DISLIKE_ACTION,
		CreatedBy:       req.UserEmail,
	}

	err = s.threadRepo.SaveThreadCommentActivityTx(threadCommentActivity, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to save thread comment activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread comment activity")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DislikeComment] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	return nil
}

func (s *threadService) unDislikeComment(ctx context.Context, req request.DislikeCommentReq, tx bun.Tx) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.unDislikeComment", "service")
	//defer endFunc()

	err := s.threadRepo.DecrementCommentDislikesCountTx(req.CommentID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unDislikeComment] Failed to decrement thread comment dislikes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread comment dislikes count")
	}

	updateValues := map[string]interface{}{
		"action":     constants.UNDISLIKE_ACTION,
		"updated_by": req.UserEmail,
		"updated_at": time.Now(),
	}

	err = s.threadRepo.UpdateThreadCommentActivityTx(req.CommentID, req.UserID, updateValues, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unDislikeComment] Failed to update thread comment activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment activity")
	}

	// TODO: Save to thread activity history

	return nil
}

func (s *threadService) dislikeCommentReply(ctx context.Context, req request.DislikeCommentReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.dislikeCommentReply", "service")
	//defer endFunc()

	lastThreadCommentReplyActivity, err := s.getUserLastThreadCommentReplyAction(ctx, req.ThreadID, req.CommentID, req.UserID)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to get user last thread comment reply action", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	threadCommentReply, err := s.threadRepo.GetThreadCommentReplyByID(req.CommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Thread comment reply does not exist", zap.Error(err))

			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread comment reply does not exist")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to get thread comment reply by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get thread comment reply by id")
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to begin transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	// Unlike if previously disliked
	if lastThreadCommentReplyActivity == constants.DISLIKE_ACTION {
		err = s.unDislikeCommentReply(ctx, req, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to undislike comment reply", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to undislike comment reply")
		}

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	// Decrement like count if previously liked
	if lastThreadCommentReplyActivity == constants.LIKE_ACTION {
		err = s.threadRepo.DecrementCommentReplyLikesCountTx(req.CommentID, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to decrement thread comment reply likes count", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread comment reply likes count")
		}
	}

	err = s.threadRepo.IncrementCommentReplyDislikesCountTx(req.CommentID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to increment thread comment reply dislikes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to increment thread comment reply dislikes count")
	}

	if lastThreadCommentReplyActivity != "" {

		updateValues := map[string]interface{}{
			"action":     constants.DISLIKE_ACTION,
			"updated_by": req.UserEmail,
			"updated_at": time.Now(),
		}

		err = s.threadRepo.UpdateThreadCommentActivityReplyTx(req.CommentID, req.UserID, updateValues, tx)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to update thread comment reply activity", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment reply activity")
		}

		// TODO: Save to thread activity history

		err = tx.Commit()
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to commit transaction", zap.Error(err))
			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
		}

		return nil
	}

	threadCommentActivity := &model.ThreadCommentActivity{
		ID:                   uuid.NewString(),
		ThreadID:             req.ThreadID,
		ThreadCommentID:      threadCommentReply.ThreadCommentID,
		ThreadCommentReplyID: pkg.ToPointer(req.CommentID),
		ActorID:              req.UserID,
		ActorEmail:           req.UserEmail,
		ActorUsername:        req.Username,
		Action:               constants.DISLIKE_ACTION,
		CreatedBy:            req.UserEmail,
	}

	err = s.threadRepo.SaveThreadCommentActivityTx(threadCommentActivity, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to save thread comment reply activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread comment reply activity")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[dislikeCommentReply] Failed to commit transaction", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	// TODO: Save to thread comment activity history

	return nil
}

func (s *threadService) unDislikeCommentReply(ctx context.Context, req request.DislikeCommentReq, tx bun.Tx) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.unDislikeCommentReply", "service")
	//defer endFunc()

	err := s.threadRepo.DecrementCommentReplyDislikesCountTx(req.CommentID, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unDislikeCommentReply] Failed to decrement thread comment reply dislikes count", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to decrement thread comment reply dislikes count")
	}

	updateValues := map[string]interface{}{
		"action":     constants.UNDISLIKE_ACTION,
		"updated_by": req.UserEmail,
		"updated_at": time.Now(),
	}

	err = s.threadRepo.UpdateThreadCommentActivityReplyTx(req.CommentID, req.UserID, updateValues, tx)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[unDislikeComment] Failed to update thread comment reply activity", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment reply activity")
	}

	// TODO: Save to thread activity history

	return nil
}

func (s *threadService) GetThreadComments(ctx context.Context, req request.GetThreadCommentsReq) (response.GetThreadCommentsResponse, error) {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.GetThreadComments", "service")
	//defer endFunc()

	var resp response.GetThreadCommentsResponse

	threadComments, err := s.threadRepo.GetThreadCommentsByThreadID(req.ThreadID, req.UserID)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[GetThreadComments] Failed to get thread comments", zap.Error(err))
		return resp, err
	}

	resp.Data = s.mapThreadCommentsData(threadComments)

	return resp, nil
}

func (s *threadService) mapThreadCommentsData(threadComments []model.ThreadComment) []response.GetThreadCommentsData {

	threadCommentsData := []response.GetThreadCommentsData{}

	for _, tc := range threadComments {

		threadCommentReplies := []response.ThreadCommentReply{}

		tcd := response.GetThreadCommentsData{
			ID:           tc.ID,
			ThreadID:     tc.ThreadID,
			UserID:       tc.UserID,
			UserName:     tc.User.Username,
			Content:      tc.Content,
			LikeCount:    tc.LikeCount,
			DislikeCount: tc.DislikeCount,
			CreatedBy:    tc.CreatedBy,
			CreatedAt:    tc.CreatedAt,
			UpdatedBy:    tc.UpdatedBy,
			UpdatedAt:    tc.UpdatedAt,
		}

		if tc.User.University != nil {
			tcd.UniversityAbbreviatedName = pkg.ToPointer(tc.User.University.AbbreviatedName)
			tcd.UniversityImageURL = pkg.ToPointer(tc.User.University.ImageURL)
		}

		if tc.CommentAction != "" {
			if tc.CommentAction == constants.LIKE_ACTION {
				tcd.IsLiked = true
			}

			if tc.CommentAction == constants.DISLIKE_ACTION {
				tcd.IsDisliked = true
			}
		}

		for _, tcr := range tc.Replies {
			tcrd := response.ThreadCommentReply{
				ID:              tcr.ID,
				ThreadID:        tcr.ThreadID,
				ThreadCommentID: tcr.ThreadCommentID,
				UserID:          tcr.UserID,
				UserName:        tcr.User.Username,
				Content:         tcr.Content,
				LikeCount:       tcr.LikeCount,
				DislikeCount:    tcr.DislikeCount,
				CreatedBy:       tcr.CreatedBy,
				CreatedAt:       tcr.CreatedAt,
				UpdatedBy:       tcr.UpdatedBy,
				UpdatedAt:       tcr.UpdatedAt,
			}

			if tcr.User.University != nil {
				tcrd.UniversityAbbreviatedName = pkg.ToPointer(tcr.User.University.AbbreviatedName)
				tcrd.UniversityImageURL = pkg.ToPointer(tcr.User.University.ImageURL)
			}

			if tcr.CommentReplyAction != "" {
				if tcr.CommentReplyAction == constants.LIKE_ACTION {
					tcrd.IsLiked = true
				}

				if tcr.CommentReplyAction == constants.DISLIKE_ACTION {
					tcrd.IsDisliked = true
				}
			}

			threadCommentReplies = append(threadCommentReplies, tcrd)
		}

		tcd.Replies = threadCommentReplies

		threadCommentsData = append(threadCommentsData, tcd)
	}

	return threadCommentsData
}

func (s *threadService) DeleteThreadComment(ctx context.Context, req request.DeleteThreadCommentReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.DeleteThreadComment", "service")
	//defer endFunc()

	_, err := s.threadRepo.GetThreadCommentByID(req.CommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[DeleteThreadComment] Thread comment not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread comment not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteThreadComment] Failed to delete thread comment by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to delete thread comment by id")
	}

	updateValues := map[string]interface{}{
		"deleted_by": req.UserEmail,
		"deleted_at": time.Now(),
	}

	err = s.threadRepo.DeleteThreadCommentByID(req.CommentID, updateValues)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteThreadComment] Failed to delete thread comment in database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to delete thread comment")
	}

	return nil
}

func (s *threadService) UpdateThreadComment(ctx context.Context, req request.UpdateThreadCommentReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.UpdateThreadComment", "service")
	//defer endFunc()

	_, err := s.threadRepo.GetThreadCommentByID(req.CommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[UpdateThreadComment] Thread comment not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread comment not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateThreadComment] Failed to get thread comment by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment by id")
	}

	updateValues := map[string]interface{}{
		"content":    req.Content,
		"updated_by": req.UserEmail,
		"updated_at": time.Now(),
	}

	err = s.threadRepo.UpdateThreadCommentByID(req.CommentID, updateValues)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateThreadComment] Failed to update thread comment in database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment")
	}

	return nil
}

func (s *threadService) DeleteThreadCommentReply(ctx context.Context, req request.DeleteThreadCommentReplyReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.DeleteThreadCommentReply", "service")
	//defer endFunc()

	_, err := s.threadRepo.GetThreadCommentReplyByID(req.CommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[DeleteThreadCommentReply] Thread comment reply not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread comment reply not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteThreadCommentReply] Failed to delete thread comment reply by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to delete thread comment reply by id")
	}

	updateValues := map[string]interface{}{
		"deleted_by": req.UserEmail,
		"deleted_at": time.Now(),
	}

	err = s.threadRepo.DeleteThreadCommentReplyByID(req.CommentID, updateValues)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[DeleteThreadCommentReply] Failed to delete thread comment reply in database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to delete thread reply comment")
	}

	return nil
}

func (s *threadService) UpdateThreadCommentReply(ctx context.Context, req request.UpdateThreadCommentReplyReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.UpdateThreadCommentReply", "service")
	//defer endFunc()

	_, err := s.threadRepo.GetThreadCommentReplyByID(req.CommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[UpdateThreadCommentReply] Thread comment reply not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread comment reply not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateThreadCommentReply] Failed to get thread comment reply by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment reply by id")
	}

	updateValues := map[string]interface{}{
		"content":    req.Content,
		"updated_by": req.UserEmail,
		"updated_at": time.Now(),
	}

	err = s.threadRepo.UpdateThreadCommentReplyByID(req.CommentID, updateValues)

	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UpdateThreadCommentReply] Failed to update thread comment reply in database", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread comment reply")
	}

	return nil
}

func (s *threadService) SubscribeThread(ctx context.Context, req request.SubscribeThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.SubscribeThread", "service")
	//defer endFunc()

	t, err := s.threadRepo.GetByID(req.ThreadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[SubscribeThread] Thread not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[SubscribeThread] Failed to get thread by id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to get thread by id")
	}

	if t.UserID == req.UserID {
		s.cfg.Logger().WarnWithContext(ctx, "[SubscribeThread] Cannot subscribe to the user's own thread")
		return oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Cannot subscribe to the user's own thread")
	}

	ts, err := s.threadRepo.GetThreadSubscriptionByUserAndThreadID(req.UserID, req.ThreadID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().ErrorWithContext(ctx, "[SubscribeThread] Failed to get thread subscription by user and thread id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	if ts.ID != "" {

		if ts.IsSubscribed {
			s.cfg.Logger().WarnWithContext(ctx, "[SubscribeThread] Thread already subscribed")
			return nil
		}

		err = s.threadRepo.UpdateThreadSubscriptionIsSubscribed(ts.ID, true)
		if err != nil {
			s.cfg.Logger().ErrorWithContext(ctx, "[SubscribeThread] Failed to update thread subscription", zap.Error(err))

			return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread subscription")
		}

		return nil
	}

	threadSubscription := &model.ThreadSubscription{
		ID:           uuid.NewString(),
		UserID:       req.UserID,
		ThreadID:     req.ThreadID,
		IsSubscribed: true,
		CreatedBy:    req.UserEmail,
	}

	err = s.threadRepo.SaveThreadSubscription(threadSubscription)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[SubscribeThread] Failed to save thread subscription", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to save thread subscription")
	}

	return nil
}

func (s *threadService) UnSubscribeThread(ctx context.Context, req request.UnSubscribeThreadReq) error {
	//ctx, endFunc := trace.Start(ctx, "ThreadService.UnSubscribeThread", "service")
	//defer endFunc()

	ts, err := s.threadRepo.GetThreadSubscriptionByUserAndThreadID(req.UserID, req.ThreadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().ErrorWithContext(ctx, "[UnSubscribeThread] Thread subscription not found", zap.Error(err))
			return oops.Code(response.NotFound.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusNotFound).Errorf("Thread subscription not found")
		}

		s.cfg.Logger().ErrorWithContext(ctx, "[UnSubscribeThread] Failed to fetch thread subscription by user and thread id", zap.Error(err))
		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrInternalServerError)
	}

	if !ts.IsSubscribed {
		s.cfg.Logger().WarnWithContext(ctx, "[UnSubscribeThread] Thread already unsubscribed")
		return nil
	}

	err = s.threadRepo.UpdateThreadSubscriptionIsSubscribed(ts.ID, false)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UnSubscribeThread] Failed to update thread subscription", zap.Error(err))

		return oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to update thread subscription")
	}

	return nil
}
