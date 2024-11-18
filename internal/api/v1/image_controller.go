package v1

import (
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/internal/service"
	"github.com/andibalo/meowhasiswa-be/pkg/apperr"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type ImageController struct {
	cfg      config.Config
	imageSvc service.ImageService
}

func NewImageController(cfg config.Config, imageSvc service.ImageService) *ImageController {

	return &ImageController{
		cfg:      cfg,
		imageSvc: imageSvc,
	}
}

func (h *ImageController) AddRoutes(r *gin.Engine) {
	ir := r.Group("/api/v1/image")

	ir.POST("/upload", h.UploadImage)
}

func (h *ImageController) UploadImage(c *gin.Context) {
	//_, endFunc := trace.Start(c.Copy().Request.Context(), "ImageController.UploadImage", "controller")
	//defer endFunc()

	fileNameParam := c.PostForm("file_name")

	file, err := c.FormFile("image")

	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UploadImage] Failed to read payload", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	// max 1 mb
	maxFileSize := int64(h.cfg.HttpMaxUploadSizeMB() * 1024 * 1024)
	if file.Size > maxFileSize {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UploadImage] File size exceeds 1 mb", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf(apperr.ErrBadRequest))
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UploadImage] Failed to read file", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrBadRequest))
		return
	}

	defer srcFile.Close()

	contentType, err := mimetype.DetectReader(srcFile)
	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UploadImage] Failed to detect mime type", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrBadRequest))
		return
	}

	var fileName string

	fileName = file.Filename

	if fileNameParam != "" {
		fileName = fileNameParam + contentType.Extension()
	}

	uniqueFileName := uuid.New().String() + "_" + fileName

	resp, err := h.imageSvc.UploadImage(c.Request.Context(), model.File{
		Reader:           srcFile,
		OriginalFileName: fileName,
		FileName:         uniqueFileName,
		Size:             file.Size,
		ContentType:      contentType.String(),
	})

	if err != nil {
		h.cfg.Logger().ErrorWithContext(c.Request.Context(), "[UploadImage] Failed to upload image", zap.Error(err))
		httpresp.HttpRespError(c, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf(apperr.ErrBadRequest))
		return
	}

	httpresp.HttpRespSuccess(c, resp, nil)

	return
}
