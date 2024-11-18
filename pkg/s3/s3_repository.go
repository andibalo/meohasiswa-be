package s3

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/pkg"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	cfg      config.Config
	client   *s3.Client
	uploader *manager.Uploader
}

func NewS3Repository(cfg config.Config, client *s3.Client) *S3Repository {

	uploader := manager.NewUploader(client)

	return &S3Repository{
		cfg:      cfg,
		client:   client,
		uploader: uploader,
	}
}

func (r *S3Repository) Upload(ctx context.Context, uploadFileData model.UploadFileDTO) (model.UploadFileOutputDTO, error) {

	var resp model.UploadFileOutputDTO

	uploadResp, err := r.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(uploadFileData.Bucket),
		Key:    aws.String(uploadFileData.Name),
		Body:   uploadFileData.File,
	})

	if err != nil {
		return resp, err
	}

	resp = model.UploadFileOutputDTO{
		Name:        uploadFileData.Name,
		URL:         uploadResp.Location,
		ETag:        pkg.NullStrToStr(uploadResp.ETag),
		VersionID:   pkg.NullStrToStr(uploadResp.VersionID),
		MegaBytes:   uploadFileData.MegaBytes,
		ContentType: uploadFileData.ContentType,
		Bucket:      uploadFileData.Bucket,
	}

	return resp, nil
}
