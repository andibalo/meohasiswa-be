package model

import (
	"github.com/uptrace/bun"
	"time"
)

type Thread struct {
	bun.BaseModel `bun:"table:thread,alias:th"`

	ID             string       `bun:",pk" json:"id"`
	UserID         string       `bun:"user_id" json:"user_id"`
	SubThreadID    string       `bun:"subthread_id" json:"subthread_id"`
	Title          string       `bun:"title" json:"title"`
	Content        string       `bun:"content" json:"content"`
	ContentSummary string       `bun:"content_summary" json:"content_summary"`
	IsActive       bool         `bun:"is_active" json:"is_active"`
	LikeCount      int64        `bun:"like_count" json:"like_count"`
	DislikeCount   int64        `bun:"dislike_count" json:"dislike_count"`
	CommentCount   int64        `bun:"comment_count" json:"comment_count"`
	CreatedBy      string       `bun:"created_by" json:"created_by"`
	CreatedAt      time.Time    `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy      *string      `json:"updated_by"`
	UpdatedAt      bun.NullTime `json:"updated_at"`
	DeletedBy      *string      `json:"-"`
	DeletedAt      time.Time    `bun:",nullzero,soft_delete" json:"-"`
}
