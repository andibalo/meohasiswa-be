package response

type UploadImageResp struct {
	URL         string
	Name        string
	ETag        string
	VersionID   string
	MegaBytes   float64
	ContentType string
	Tags        string
	Bucket      string
}
