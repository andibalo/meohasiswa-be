package model

import (
	"github.com/uptrace/bun"
	"time"
)

type SubThread struct {
	bun.BaseModel `bun:"table:subthread,alias:st"`

	ID                    string       `bun:",pk" json:"id"`
	Name                  string       `bun:"name" json:"name"`
	FollowersCount        int64        `bun:"followers_count" json:"followers_count"`
	Description           string       `bun:"description" json:"description"`
	ImageUrl              string       `bun:"image_url" json:"image_url"`
	LabelColor            string       `bun:"label_color" json:"label_color"`
	UniversityID          *string      `bun:"university_id" json:"university_id"`
	IsUniversitySubThread bool         `bun:"is_university_subthread" json:"is_university_subthread"`
	CreatedBy             string       `bun:"created_by" json:"created_by"`
	CreatedAt             time.Time    `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy             *string      `json:"updated_by"`
	UpdatedAt             bun.NullTime `json:"updated_at"`
	DeletedBy             *string      `json:"-"`
	DeletedAt             time.Time    `bun:",nullzero,soft_delete" json:"-"`
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
