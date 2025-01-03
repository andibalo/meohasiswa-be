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
	Search              string `json:"_q"`
	IsTrending          bool   `json:"is_trending"`
	IsUserFollowing     bool   `json:"is_user_following"`
	UserIDParam         string `json:"user_id"`
	Limit               int    `json:"limit"`
	Cursor              string `json:"cursor"`
	IncludeUserActivity bool   `json:"include_user_activity"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}

type GetThreadDetailReq struct {
	ThreadID string `json:"thread_id"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
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

type GetThreadCommentsReq struct {
	ThreadID string `json:"-"`

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

type DeleteThreadReq struct {
	ThreadID string `json:"thread_id"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type LikeCommentReq struct {
	ThreadID  string `json:"thread_id"`
	IsReply   bool   `json:"is_reply"`
	CommentID string `json:"-"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type DislikeCommentReq struct {
	ThreadID  string `json:"thread_id" binding:"required"`
	IsReply   bool   `json:"is_reply"`
	CommentID string `json:"-"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type DeleteThreadCommentReq struct {
	CommentID string `json:"comment_id"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type UpdateThreadCommentReq struct {
	Content string `json:"content" binding:"required"`

	CommentID string `json:"-"`
	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type DeleteThreadCommentReplyReq struct {
	CommentID string `json:"comment_id"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type UpdateThreadCommentReplyReq struct {
	Content string `json:"content" binding:"required"`

	CommentID string `json:"-"`
	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type SubscribeThreadReq struct {
	ThreadID string `json:"-"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}

type UnSubscribeThreadReq struct {
	ThreadID string `json:"-"`

	UserID    string `json:"-"`
	Username  string `json:"-"`
	UserEmail string `json:"-"`
}
