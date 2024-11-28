package model

import (
	"github.com/uptrace/bun"
	"time"
)

type Thread struct {
	bun.BaseModel `bun:"table:thread,alias:th"`

	ID             string           `bun:",pk" json:"id"`
	UserID         string           `bun:"user_id" json:"user_id"`
	User           User             `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	SubThreadID    string           `bun:"subthread_id" json:"subthread_id"`
	SubThread      SubThread        `bun:"rel:belongs-to,join:subthread_id=id" json:"subthread"`
	Title          string           `bun:"title" json:"title"`
	Content        string           `bun:"content" json:"content"`
	ContentSummary string           `bun:"content_summary" json:"content_summary"`
	IsActive       bool             `bun:"is_active" json:"is_active"`
	LikeCount      int64            `bun:"like_count" json:"like_count"`
	DislikeCount   int64            `bun:"dislike_count" json:"dislike_count"`
	CommentCount   int64            `bun:"comment_count" json:"comment_count"`
	TrendingScore  float64          `bun:"trending_score,scanonly"`
	ThreadAction   string           `bun:"thread_action,scanonly"`
	Comments       []*ThreadComment `bun:"rel:has-many,join:id=thread_id"`
	CreatedBy      string           `bun:"created_by" json:"created_by"`
	CreatedAt      time.Time        `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy      *string          `json:"updated_by"`
	UpdatedAt      bun.NullTime     `json:"updated_at"`
	DeletedBy      *string          `json:"-"`
	DeletedAt      time.Time        `bun:",nullzero,soft_delete" json:"-"`
}

type ThreadActivity struct {
	bun.BaseModel `bun:"table:thread_activity,alias:tha"`

	ID            string    `bun:",pk" json:"id"`
	ThreadID      string    `bun:"thread_id" json:"thread_id"`
	Thread        Thread    `bun:"rel:belongs-to,join:thread_id=id" json:"thread"`
	ActorID       string    `bun:"actor_id" json:"actor_id"`
	ActorEmail    string    `bun:"actor_email" json:"actor_email"`
	ActorUsername string    `bun:"actor_username" json:"actor_username"`
	Action        string    `bun:"action" json:"action"`
	CreatedBy     string    `bun:"created_by" json:"created_by"`
	CreatedAt     time.Time `bun:",nullzero,default:now()" json:"created_at"`
}

type ThreadComment struct {
	bun.BaseModel `bun:"table:thread_comment,alias:thc"`

	ID           string       `bun:",pk" json:"id"`
	UserID       string       `bun:"user_id" json:"user_id"`
	User         User         `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	ThreadID     string       `bun:"thread_id" json:"thread_id"`
	Thread       SubThread    `bun:"rel:belongs-to,join:thread_id=id" json:"thread"`
	Content      string       `bun:"content" json:"content"`
	LikeCount    int64        `bun:"like_count" json:"like_count"`
	DislikeCount int64        `bun:"dislike_count" json:"dislike_count"`
	ReplyCount   int64        `bun:"reply_count" json:"reply_count"`
	CreatedBy    string       `bun:"created_by" json:"created_by"`
	CreatedAt    time.Time    `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy    *string      `json:"updated_by"`
	UpdatedAt    bun.NullTime `json:"updated_at"`
	DeletedBy    *string      `json:"-"`
	DeletedAt    time.Time    `bun:",nullzero,soft_delete" json:"-"`
}

type ThreadCommentActivity struct {
	bun.BaseModel `bun:"table:thread_comment_activity,alias:thca"`

	ID                   string              `bun:",pk" json:"id"`
	ThreadID             string              `bun:"thread_id" json:"thread_id"`
	Thread               Thread              `bun:"rel:belongs-to,join:thread_id=id" json:"thread"`
	ThreadCommentID      string              `bun:"thread_comment_id" json:"thread_comment_id"`
	ThreadComment        ThreadComment       `bun:"rel:belongs-to,join:thread_comment_id=id" json:"thread_comment"`
	ThreadCommentReplyID *string             `bun:"thread_comment_reply_id" json:"thread_comment_reply_id"`
	ThreadCommentReply   *ThreadCommentReply `bun:"rel:belongs-to,join:thread_comment_reply_id=id" json:"thread_comment_reply"`
	ActorID              string              `bun:"actor_id" json:"actor_id"`
	ActorEmail           string              `bun:"actor_email" json:"actor_email"`
	ActorUsername        string              `bun:"actor_username" json:"actor_username"`
	Action               string              `bun:"action" json:"action"`
	CreatedBy            string              `bun:"created_by" json:"created_by"`
	CreatedAt            time.Time           `bun:",nullzero,default:now()" json:"created_at"`
}

type ThreadCommentReply struct {
	bun.BaseModel `bun:"table:thread_comment_reply,alias:thcr"`

	ID              string       `bun:",pk" json:"id"`
	UserID          string       `bun:"user_id" json:"user_id"`
	User            User         `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	ThreadID        string       `bun:"thread_id" json:"thread_id"`
	Thread          SubThread    `bun:"rel:belongs-to,join:thread_id=id" json:"thread"`
	ThreadCommentID string       `bun:"thread_comment_id" json:"thread_comment_id"`
	Content         string       `bun:"content" json:"content"`
	LikeCount       int64        `bun:"like_count" json:"like_count"`
	DislikeCount    int64        `bun:"dislike_count" json:"dislike_count"`
	CreatedBy       string       `bun:"created_by" json:"created_by"`
	CreatedAt       time.Time    `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy       *string      `json:"updated_by"`
	UpdatedAt       bun.NullTime `json:"updated_at"`
	DeletedBy       *string      `json:"-"`
	DeletedAt       time.Time    `bun:",nullzero,soft_delete" json:"-"`
}
