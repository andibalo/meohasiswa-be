package request

type CreateThreadReq struct {
	SubThreadID    string `json:"subthread_id" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Content        string `json:"content" binding:"required"`
	ContentSummary string `json:"content_summary" binding:"required"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}

type GetThreadListReq struct {
	IsTrending      bool   `json:"is_trending"`
	IsUserFollowing bool   `json:"is_user_following"`
	UserIDParam     string `json:"user_id"`
	Limit           int    `json:"limit"`
	Cursor          string `json:"cursor"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}

type GetThreadDetailReq struct {
	ThreadID string `json:"thread_id"`
}

type LikeThreadReq struct {
	ThreadID string `json:"thread_id"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type DislikeThreadReq struct {
	ThreadID string `json:"thread_id"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type CommentThreadReq struct {
	Content string `json:"content" binding:"required"`

	ThreadID  string `json:"-"`
	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type ReplyCommentReq struct {
	Content  string `json:"content" binding:"required"`
	ThreadID string `json:"thread_id" binding:"required"`

	CommentID string `json:"-"`
	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type UpdateThreadReq struct {
	ThreadID       string `json:"thread_id"`
	Title          string `json:"title" binding:"required"`
	Content        string `json:"content" binding:"required"`
	ContentSummary string `json:"content_summary" binding:"required"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}
