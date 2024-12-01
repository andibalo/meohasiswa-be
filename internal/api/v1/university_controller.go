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

type UniversityController struct {
	cfg           config.Config
	universitySvc service.UniversityService
}

func NewUniversityController(cfg config.Config, universitySvc service.UniversityService) *UniversityController {

	return &UniversityController{
		cfg:           cfg,
		universitySvc: universitySvc,
	}
}

func (h *UniversityController) AddRoutes(r *gin.Engine) {
	ur := r.Group("/api/v1/university")

	ur.GET("/ratings", middleware.JwtMiddleware(h.cfg), h.GetUniversityRatingList)
	ur.GET("/rating/:rating_id", middleware.JwtMiddleware(h.cfg), h.GetUniversityRatingDetail)
	ur.POST("/rate/:university_id", middleware.JwtMiddleware(h.cfg), h.RateUniversity)
	ur.PATCH("/rating/:rating_id", middleware.JwtMiddleware(h.cfg), h.UpdateUniversityRating)
}

func (h *UniversityController) GetUniversityRatingList(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UniversityController.GetUniversityRatingList", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.GetUniversityRatingListReq

	limit, err := pkg.GetIntQueryParams(c, 10, "limit")
	if err != nil {
		httpresp.HttpRespError(c, err)
		return
	}

	data.Limit = limit
	data.Cursor = c.Query("cursor")
	data.Search = c.Query("_q")

	data.UserID = claims.ID
	data.UserEmail = claims.Email

	resp, err := h.universitySvc.GetUniversityRatingList(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetUniversityRatingList] Failed to get university rating list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, resp, nil)
	return
}

func (h *UniversityController) GetUniversityRatingDetail(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UniversityController.GetUniversityRatingDetail", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.GetUniversityRatingDetailReq

	data.UniversityRatingID = c.Param("rating_id")

	data.UserID = claims.ID
	data.UserEmail = claims.Email
	data.Username = claims.UserName

	resp, err := h.universitySvc.GetUniversityRatingDetail(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[GetUniversityRatingDetail] Failed to get university rating detail", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, resp, nil)
	return
}

func (h *UniversityController) RateUniversity(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UniversityController.RateUniversity", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.RateUniversityReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[RateUniversity] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UniversityID = c.Param("university_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email

	err := h.universitySvc.CreateUniversityRating(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[RateUniversity] Failed to rate university", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}

func (h *UniversityController) UpdateUniversityRating(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "UniversityController.UpdateUniversityRating", "controller")
	//defer endFunc()

	claims := middleware.ParseToken(c)
	if len(claims.Token) == 0 {
		httpresp.HttpRespError(c, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(apperr.ErrUnauthorized))
		return
	}

	var data request.UpdateUniversityRatingReq

	if err := c.ShouldBindJSON(&data); err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UpdateUniversityRating] Failed to bind json", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	data.UniversityRatingID = c.Param("rating_id")
	data.UserID = claims.ID
	data.UserEmail = claims.Email

	err := h.universitySvc.UpdateUniversityRating(c.Request.Context(), data)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[RateUniversity] Failed to rate university", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
	return
}
