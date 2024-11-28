package response

import (
	"github.com/uptrace/bun"
	"time"
)

type ThreadListData struct {
	ID                        string       `json:"id"`
	UserID                    string       `json:"user_id"`
	UserName                  string       `json:"username"`
	UniversityAbbreviatedName *string      `json:"university_abbreviated_name"`
	UniversityImageURL        *string      `json:"university_image_url"`
	SubThreadID               string       `json:"subthread_id"`
	SubThreadName             string       `json:"subthread_name"`
	SubThreadColor            string       `json:"subthread_color"`
	Title                     string       `json:"title"`
	Content                   string       `json:"content"`
	ContentSummary            string       `json:"content_summary"`
	IsActive                  bool         `json:"is_active"`
	LikeCount                 int64        `json:"like_count"`
	DislikeCount              int64        `json:"dislike_count"`
	CommentCount              int64        `json:"comment_count"`
	IsLiked                   bool         `json:"is_liked"`
	IsDisliked                bool         `json:"is_disliked"`
	CreatedBy                 string       `json:"created_by"`
	CreatedAt                 time.Time    `json:"created_at"`
	UpdatedBy                 *string      `json:"updated_by"`
	UpdatedAt                 bun.NullTime `json:"updated_at"`
}

type GetThreadListResponse struct {
	Data []ThreadListData `json:"threads"`
	Meta PaginationMeta   `json:"meta"`
}

type ThreadDetailData struct {
	ID                        string       `json:"id"`
	UserID                    string       `json:"user_id"`
	UserName                  string       `json:"username"`
	UniversityAbbreviatedName *string      `json:"university_abbreviated_name"`
	UniversityImageURL        *string      `json:"university_image_url"`
	SubThreadID               string       `json:"subthread_id"`
	SubThreadName             string       `json:"subthread_name"`
	SubThreadColor            string       `json:"subthread_color"`
	Title                     string       `json:"title"`
	Content                   string       `json:"content"`
	ContentSummary            string       `json:"content_summary"`
	IsActive                  bool         `json:"is_active"`
	LikeCount                 int64        `json:"like_count"`
	DislikeCount              int64        `json:"dislike_count"`
	CommentCount              int64        `json:"comment_count"`
	IsLiked                   bool         `json:"is_liked"`
	IsDisliked                bool         `json:"is_disliked"`
	CreatedBy                 string       `json:"created_by"`
	CreatedAt                 time.Time    `json:"created_at"`
	UpdatedBy                 *string      `json:"updated_by"`
	UpdatedAt                 bun.NullTime `json:"updated_at"`
}

type ThreadComment struct {
	ID                        string       `json:"id"`
	UserID                    string       `json:"user_id"`
	UserName                  string       `json:"username"`
	UniversityAbbreviatedName *string      `json:"university_abbreviated_name"`
	UniversityImageURL        *string      `json:"university_image_url"`
	Content                   string       `json:"content"`
	LikeCount                 int64        `json:"like_count"`
	DislikeCount              int64        `json:"dislike_count"`
	CreatedBy                 string       `json:"created_by"`
	CreatedAt                 time.Time    `json:"created_at"`
	UpdatedBy                 *string      `json:"updated_by"`
	UpdatedAt                 bun.NullTime `json:"updated_at"`
}

type GetThreadDetailResponse struct {
	Data ThreadDetailData `json:"thread"`
}

type GetThreadCommentsData struct {
	ID                        string               `json:"id"`
	ThreadID                  string               `json:"thread_id"`
	UserID                    string               `json:"user_id"`
	UserName                  string               `json:"username"`
	UniversityAbbreviatedName *string              `json:"university_abbreviated_name"`
	UniversityImageURL        *string              `json:"university_image_url"`
	Content                   string               `json:"content"`
	LikeCount                 int64                `json:"like_count"`
	DislikeCount              int64                `json:"dislike_count"`
	IsLiked                   bool                 `json:"is_liked"`
	IsDisliked                bool                 `json:"is_disliked"`
	Replies                   []ThreadCommentReply `json:"replies"`
	CreatedBy                 string               `json:"created_by"`
	CreatedAt                 time.Time            `json:"created_at"`
	UpdatedBy                 *string              `json:"updated_by"`
	UpdatedAt                 bun.NullTime         `json:"updated_at"`
}

type ThreadCommentReply struct {
	ID                        string       `json:"id"`
	ThreadID                  string       `json:"thread_id"`
	ThreadCommentID           string       `json:"thread_comment_id"`
	UserID                    string       `json:"user_id"`
	UserName                  string       `json:"username"`
	UniversityAbbreviatedName *string      `json:"university_abbreviated_name"`
	UniversityImageURL        *string      `json:"university_image_url"`
	Content                   string       `json:"content"`
	LikeCount                 int64        `json:"like_count"`
	DislikeCount              int64        `json:"dislike_count"`
	IsLiked                   bool         `json:"is_liked"`
	IsDisliked                bool         `json:"is_disliked"`
	CreatedBy                 string       `json:"created_by"`
	CreatedAt                 time.Time    `json:"created_at"`
	UpdatedBy                 *string      `json:"updated_by"`
	UpdatedAt                 bun.NullTime `json:"updated_at"`
}

type GetThreadCommentsResponse struct {
	Data []GetThreadCommentsData `json:"thread_comments"`
}
