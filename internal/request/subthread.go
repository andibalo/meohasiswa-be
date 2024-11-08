package request

type CreateSubThreadReq struct {
	Name                  string  `json:"name" binding:"required"`
	Description           string  `json:"description" binding:"required"`
	ImageUrl              string  `json:"image_url" binding:"required"`
	UniversityID          *string `json:"university_id"`
	IsUniversitySubThread bool    `json:"is_university_subthread"`
}

type FollowSubThreadReq struct {
	SubThreadID string `json:"subthread_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
}

type UnFollowSubThreadReq struct {
	SubThreadID string `json:"subthread_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
}
