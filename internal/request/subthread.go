package request

type CreateSubThreadReq struct {
	Name                  string  `json:"name" binding:"required"`
	Description           string  `json:"description" binding:"required"`
	ImageUrl              string  `json:"image_url" binding:"required"`
	LabelColor            string  `json:"label_color" binding:"required"`
	UniversityID          *string `json:"university_id"`
	IsUniversitySubThread bool    `json:"is_university_subthread"`

	UserEmail string `json:"-"`
}

type FollowSubThreadReq struct {
	SubThreadID string `json:"subthread_id" binding:"required"`

	UserID string `json:"-"`
}

type UnFollowSubThreadReq struct {
	SubThreadID string `json:"subthread_id" binding:"required"`

	UserID string `json:"-"`
}

type GetSubThreadListReq struct {
	IsFollowing bool   `json:"is_following"`
	Limit       int    `json:"limit"`
	Cursor      string `json:"cursor"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}
