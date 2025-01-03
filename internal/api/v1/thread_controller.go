package v1

import (
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/middleware"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/internal/service"
	"github.com/andibalo/meowhasiswa-be/pkg"
	"github.com/andibalo/meowhasiswa-be/pkg/apperr"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type ThreadController struct {
	cfg       config.Config
	threadSvc service.ThreadService
}

func NewThreadController(cfg config.Config, threadSvc service.ThreadService) *ThreadController {

	return &ThreadController{
		cfg:       cfg,
		threadSvc: threadSvc,
	}
}

func (h *ThreadController) AddRoutes(r *gin.Engine) {
	tr := r.Group("/api/v1/thread")

	tr.POST("", middleware.JwtMiddleware(h.cfg), h.CreateThread)
	tr.GET("", middleware.JwtMiddleware(h.cfg), h.GetThreadList)
	tr.GET("/:thread_id", middleware.JwtMiddleware(h.cfg), h.GetThreadDetail)
	tr.DELETE("/:thread_id", middleware.JwtMiddleware(h.cfg), h.DeleteThread)
	tr.PATCH("/:thread_id", middleware.JwtMiddleware(h.cfg), h.UpdateThread)
	tr.POST("/subscribe/:thread_id", middleware.JwtMiddleware(h.cfg), h.SubscribeThread)
	tr.PATCH("/unsubscribe/:thread_id", middleware.JwtMiddleware(h.cfg), h.UnSubscribeThread)
	tr.PATCH("/like/:thread_id", middleware.JwtMiddleware(h.cfg), h.LikeThread)
	tr.PATCH("/dislike/:thread_id", middleware.JwtMiddleware(h.cfg), h.DislikeThread)
	tr.GET("/comment/:thread_id", middleware.JwtMiddleware(h.cfg), h.GetThreadComments)
	tr.POST("/comment/:thread_id", middleware.JwtMiddleware(h.cfg), h.CommentThread)
	tr.DELETE("/comment/:comment_id", middleware.JwtMiddleware(h.cfg), h.DeleteThreadComment)
	tr.PATCH("/comment/:comment_id", middleware.JwtMiddleware(h.cfg), h.UpdateThreadComment)
	tr.POST("/comment/reply/:comment_id", middleware.JwtMiddleware(h.cfg), h.ReplyComment)
	tr.DELETE("/comment/reply/:comment_id", middleware.JwtMiddleware(h.cfg), h.DeleteThreadCommentReply)
	tr.PATCH("/comment/reply/:comment_id", middleware.JwtMiddleware(h.cfg), h.UpdateThreadCommentReply)
	tr.PATCH("/comment/like/:comment_id", middleware.JwtMiddleware(h.cfg), h.LikeComment)
	tr.PATCH("/comment/dislike/:comment_id", middleware.JwtMiddleware(h.cfg), h.DislikeComment)
}

func (h *ThreadController) GetThreadList(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.GetThreadList", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.GetThreadListReq

	limit, err := pkg.GetIntQueryParams(c, 10, "limit")
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	isTrending, err := pkg.GetBoolQueryParams(c, "is_trending")
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	isUserFollowing, err := pkg.GetBoolQueryParams(c, "is_user_following")
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	includeUserActivity, err := pkg.GetBoolQueryParams(c, "include_user_activity")
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	data.Limit = limit
	data.IsTrending = isTrending
	data.IsUserFollowing = isUserFollowing
	data.Cursor = c.Query("cursor")
	data.UserIDParam = c.Query("user_id")
	data.Search = c.Query("_q")
	data.IncludeUserActivity = includeUserActivity

	data.UserID = claims.ID
	data.UserEmail = claims.Email

	resp, err := h.threadSvc.GetThreadList(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetThreadList] Failed to get thread list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, resp, nil)
	return
}

func (h *ThreadController) GetThreadDetail(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.GetThreadDetail", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.GetThreadDetailReq

	data.ThreadID = c.Param("thread_id")

	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	resp, err := h.threadSvc.GetThreadDetail(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetThreadDetail] Failed to get thread detail", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, resp, nil)
	return
}

func (h *ThreadController) UpdateThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.UpdateThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.UpdateThreadReq
	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateThread] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.ThreadID = c.Param("thread_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.UpdateThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateThread] Failed to update thread by id", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) CreateThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.CreateThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.CreateThreadReq
	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateThread] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UserID = claims.ID
	data.UserEmail = claims.Email

	err := h.threadSvc.CreateThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateThread] Failed to create thread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) DeleteThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.DeleteThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.DeleteThreadReq

	data.ThreadID = c.Param("thread_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.DeleteThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[DeleteThread] Failed to delete thread by id", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) LikeThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.LikeThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.LikeThreadReq

	data.ThreadID = c.Param("thread_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.LikeThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[LikeThread] Failed to like thread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) DislikeThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.DislikeThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.DislikeThreadReq

	data.ThreadID = c.Param("thread_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.DislikeThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[DislikeThread] Failed to dislike thread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) GetThreadComments(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.GetThreadComments", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.GetThreadCommentsReq

	data.ThreadID = c.Param("thread_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	resp, err := h.threadSvc.GetThreadComments(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetThreadComments] Failed to get thread comments", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, resp, nil)
	return
}

func (h *ThreadController) CommentThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.CommentThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.CommentThreadReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CommentThread] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.ThreadID = c.Param("thread_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.CommentThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CommentThread] Failed to comment thread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) UpdateThreadComment(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.UpdateThreadComment", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.UpdateThreadCommentReq
	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateThreadComment] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.CommentID = c.Param("comment_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.UpdateThreadComment(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateThreadComment] Failed to update thread comment by id", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) ReplyComment(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.ReplyComment", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.ReplyCommentReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[ReplyComment] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.CommentID = c.Param("comment_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.ReplyComment(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[ReplyComment] Failed to reply comment", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) LikeComment(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.LikeComment", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.LikeCommentReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[LikeComment] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.CommentID = c.Param("comment_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.LikeComment(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[LikeComment] Failed to like comment", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) DislikeComment(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.DislikeComment", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.DislikeCommentReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[DislikeCommentReq] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.CommentID = c.Param("comment_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.DislikeComment(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[DislikeComment] Failed to dislike comment", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) DeleteThreadComment(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.DeleteThreadComment", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.DeleteThreadCommentReq

	data.CommentID = c.Param("comment_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.DeleteThreadComment(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[DeleteThreadComment] Failed to delete thread comment by id", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) DeleteThreadCommentReply(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.DeleteThreadCommentReply", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.DeleteThreadCommentReplyReq

	data.CommentID = c.Param("comment_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.DeleteThreadCommentReply(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[DeleteThreadCommentReply] Failed to delete thread comment reply by id", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) UpdateThreadCommentReply(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.UpdateThreadCommentReply", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.UpdateThreadCommentReplyReq
	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateThreadCommentReply] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.CommentID = c.Param("comment_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.UpdateThreadCommentReply(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateThreadCommentReply] Failed to update thread comment reply by id", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) SubscribeThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.SubscribeThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.SubscribeThreadReq

	data.ThreadID = c.Param("thread_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.SubscribeThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[SubscribeThread] Failed to subscribe thread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *ThreadController) UnSubscribeThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ThreadController.UnSubscribeThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.UnSubscribeThreadReq

	data.ThreadID = c.Param("thread_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.threadSvc.UnSubscribeThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UnSubscribeThread] Failed to unsubscribe thread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}
