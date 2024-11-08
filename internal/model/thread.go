package model

import (
	"github.com/uptrace/bun"
	"time"
)

type Thread struct {
	bun.BaseModel `bun:"table:thread,alias:th"`

	ID             string    `bun:",pk"`
	UserID         string    `bun:"user_id"`
	SubThreadID    string    `bun:"subthread_id"`
	Title          string    `bun:"title"`
	Content        string    `bun:"content"`
	ContentSummary string    `bun:"content_summary"`
	IsActive       bool      `bun:"is_active"`
	LikeCount      int64     `bun:"like_count"`
	DislikeCount   int64     `bun:"dislike_count"`
	CommentCount   int64     `bun:"comment_count"`
	CreatedBy      string    `bun:"created_by"`
	CreatedAt      time.Time `bun:",nullzero,default:now()"`
	UpdatedBy      *string
	UpdatedAt      bun.NullTime
	DeletedBy      *string
	DeletedAt      time.Time `bun:",nullzero,soft_delete"`
}
