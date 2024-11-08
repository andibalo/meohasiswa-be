package v1

import (
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/internal/service"
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

	str.POST("", h.CreateSubThread)
	str.POST("/follow", h.FollowSubThread)
	str.PATCH("/unfollow", h.UnfollowSubThread)
}

func (h *SubThreadController) CreateSubThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "SubThreadController.CreateSubThread", "controller")
	//defer endFunc()

	var data request.CreateSubThreadReq

	if err := c.ShouldBindJSON(&data); err != nil {
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	err := h.subThreadSvc.CreateSubThread(c.Request.Context(), data)

	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[CreateSubThread] Failed to create subthread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *SubThreadController) FollowSubThread(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "SubThreadController.FollowSubThread", "controller")
	//defer endFunc()

	var data request.FollowSubThreadReq

	if err := c.ShouldBindJSON(&data); err != nil {
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

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

	var data request.UnFollowSubThreadReq

	if err := c.ShouldBindJSON(&data); err != nil {
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	err := h.subThreadSvc.UnFollowSubThread(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UnfollowSubThread] Failed to unfollow subthread", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}
