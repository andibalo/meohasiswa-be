package model

import (
	"github.com/uptrace/bun"
	"time"
)

type SubThread struct {
	bun.BaseModel `bun:"table:subthread,alias:st"`

	ID                    string    `bun:",pk"`
	Name                  string    `bun:"name"`
	FollowersCount        int64     `bun:"followers_count"`
	Description           string    `bun:"description"`
	ImageUrl              string    `bun:"image_url"`
	UniversityID          *string   `bun:"university_id"`
	IsUniversitySubThread bool      `bun:"is_university_subthread"`
	CreatedBy             string    `bun:"created_by"`
	CreatedAt             time.Time `bun:",nullzero,default:now()"`
	UpdatedBy             *string
	UpdatedAt             bun.NullTime
	DeletedBy             *string
	DeletedAt             time.Time `bun:",nullzero,soft_delete"`
}

type SubThreadFollower struct {
	bun.BaseModel `bun:"table:subthread_follower,alias:stf"`

	ID          string    `bun:",pk"`
	UserID      string    `bun:"user_id"`
	SubThreadID string    `bun:"subthread_id"`
	IsFollowing bool      `bun:"is_following"`
	CreatedBy   string    `bun:"created_by"`
	CreatedAt   time.Time `bun:",nullzero,default:now()"`
	UpdatedBy   *string
	UpdatedAt   bun.NullTime
	DeletedBy   *string
	DeletedAt   time.Time `bun:",nullzero,soft_delete"`
}
