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
	ar := r.Group("/api/v1/thread")

	ar.POST("", middleware.JwtMiddleware(h.cfg), h.CreateThread)
	ar.GET("", middleware.JwtMiddleware(h.cfg), h.GetThreadList)
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

	data.Limit = limit
	data.IsTrending = isTrending
	data.Cursor = c.Query("cursor")

	resp, err := h.threadSvc.GetThreadList(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetThreadList] Failed to get thread list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, resp, nil)
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
