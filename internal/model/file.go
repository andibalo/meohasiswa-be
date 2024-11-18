package model

import (
	"image"
	"io"
)

type File struct {
	ETag             string      `json:"e_tag"`
	VersionId        string      `json:"version_id"`
	FileDir          string      `json:"file_dir"`
	OriginalFileName string      `json:"original_file_name"`
	FileName         string      `json:"file_name"`
	Slug             string      `json:"slug"`
	Buffer           []byte      `json:"buffer"`
	Reader           io.Reader   `json:"reader"`
	Size             int64       `json:"size"`
	ContentType      string      `json:"content_type"`
	Img              image.Image `json:"img"`
}

type UploadFileDTO struct {
	File          io.Reader
	Name          string
	Album         string
	Bucket        string
	IsPrivate     bool
	MegaBytes     float64
	Region        string
	ContentType   string
	WithWatermark bool
	AuthInternal  bool //
}

type UploadFileOutputDTO struct {
	URL         string
	Name        string
	ETag        string
	VersionID   string
	MegaBytes   float64
	ContentType string
	Tags        string
	Bucket      string
}
