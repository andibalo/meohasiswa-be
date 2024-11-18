package service

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/repository"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
)

type imageService struct {
	cfg      config.Config
	fileRepo repository.FileRepository
}

func NewImageService(cfg config.Config, fileRepo repository.FileRepository) ImageService {

	return &imageService{
		cfg:      cfg,
		fileRepo: fileRepo,
	}
}

func (s *imageService) UploadImage(ctx context.Context, fileData model.File) (response.UploadImageResp, error) {
	//ctx, endFunc := trace.Start(ctx, "ImageService.UploadImage", "service")
	//defer endFunc()

	var resp response.UploadImageResp

	uploadData := model.UploadFileDTO{
		File:        fileData.Reader,
		Name:        fileData.FileName,
		Bucket:      s.cfg.GetAWSCfg().DefaultBucket,
		IsPrivate:   false,
		MegaBytes:   float64(fileData.Size),
		ContentType: fileData.ContentType,
	}

	uploadResp, err := s.fileRepo.Upload(ctx, uploadData)
	if err != nil {
		s.cfg.Logger().ErrorWithContext(ctx, "[UploadImage] Failed to upload image to remote file storage", zap.Error(err))

		return resp, oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to upload image to remote file storage")

	}

	resp = response.UploadImageResp{
		URL:         uploadResp.URL,
		Name:        uploadResp.Name,
		ETag:        uploadResp.ETag,
		VersionID:   uploadResp.VersionID,
		MegaBytes:   uploadResp.MegaBytes,
		ContentType: uploadResp.ContentType,
		Bucket:      uploadResp.Bucket,
	}

	return resp, nil
}
