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
	Title                     string       `json:"title"`
	Content                   string       `json:"content"`
	ContentSummary            string       `json:"content_summary"`
	IsActive                  bool         `json:"is_active"`
	LikeCount                 int64        `json:"like_count"`
	DislikeCount              int64        `json:"dislike_count"`
	CommentCount              int64        `json:"comment_count"`
	CreatedBy                 string       `json:"created_by"`
	CreatedAt                 time.Time    `json:"created_at"`
	UpdatedBy                 *string      `json:"updated_by"`
	UpdatedAt                 bun.NullTime `json:"updated_at"`
}

type GetThreadListResponse struct {
	Data []ThreadListData `json:"threads"`
	Meta PaginationMeta   `json:"meta"`
}
