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

type SubThreadController struct {
	cfg          config.Config
	subThreadSvc service.SubThreadService
}

func NewSubThreadController(cfg config.Config, subThreadSvc service.SubThreadService) *SubThreadController {

	return &SubThreadController{
		cfg:          cfg,
		subThreadSvc: subThreadSvc,
	}
}

func (h *SubThreadController) AddRoutes(r *gin.Engine) {
	str := r.Group("/api/v1/subthread")

	str.GET("", middleware.JwtMiddleware(h.cfg), h.GetListSubThread)
	str.POST("", middleware.JwtMiddleware(h.cfg), h.CreateSubThread)
	str.PATCH("/:subthread_id", middleware.JwtMiddleware(h.cfg), h.UpdateSubThread)
	str.DELETE("/:subthread_id", middleware.JwtMiddleware(h.cfg), h.DeleteSubThread)
	str.POST("/follow", middleware.JwtMiddleware(h.cfg), h.FollowSubThread)
	str.PATCH("/unfollow", middleware.JwtMiddleware(h.cfg), h.UnfollowSubThread)
}

func (h *SubThreadController) GetListSubThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "SubThreadController.GetListSubThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.GetSubThreadListReq

	limit, err := pkg.GetIntQueryParams(c, 10, "limit")
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	isFollowing, err := pkg.GetBoolQueryParams(c, "is_following")
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	includeUniversitySubThread, err := pkg.GetBoolQueryParams(c, "include_university_subthread")
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	data.Limit = limit
	data.IsFollowing = isFollowing
	data.Cursor = c.Query("cursor")
	data.Search = c.Query("_q")
	data.IncludeUniversitySubThread = includeUniversitySubThread

	data.UserID = claims.ID
	data.UserEmail = claims.Email

	resp, err := h.subThreadSvc.GetSubThreadList(c.Request.Context(), data)

	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetListSubThread] Failed to get subthread list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, resp, nil)
	return
}

func (h *SubThreadController) CreateSubThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "SubThreadController.CreateSubThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.CreateSubThreadReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateSubThread] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UserEmail = claims.Email

	err := h.subThreadSvc.CreateSubThread(c.Request.Context(), data)

	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateSubThread] Failed to create subthread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *SubThreadController) UpdateSubThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "SubThreadController.UpdateSubThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.UpdateSubThreadReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateSubThread] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.SubThreadID = c.Param("subthread_id")
	data.UserEmail = claims.Email

	err := h.subThreadSvc.UpdateSubThread(c.Request.Context(), data)

	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateSubThread] Failed to update subthread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *SubThreadController) DeleteSubThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "SubThreadController.DeleteSubThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.DeleteSubThreadReq

	data.SubThreadID = c.Param("subthread_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	err := h.subThreadSvc.DeleteSubThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[DeleteSubThread] Failed to delete subthread by id", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *SubThreadController) FollowSubThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "SubThreadController.FollowSubThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.FollowSubThreadReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[FollowSubThread] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UserID = claims.ID

	err := h.subThreadSvc.FollowSubThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[FollowSubThread] Failed to follow subthread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *SubThreadController) UnfollowSubThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "SubThreadController.UnfollowSubThread", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.UnFollowSubThreadReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UnfollowSubThread] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UserID = claims.ID

	err := h.subThreadSvc.UnFollowSubThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UnfollowSubThread] Failed to unfollow subthread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}
