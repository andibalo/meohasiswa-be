package request

type CreateThreadReq struct {
	SubThreadID    string `json:"subthread_id" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Content        string `json:"content" binding:"required"`
	ContentSummary string `json:"content_summary" binding:"required"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}
