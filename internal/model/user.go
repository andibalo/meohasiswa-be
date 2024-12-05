package model

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:user,alias:u"`

	ID                 string       `bun:",pk" json:"id"`
	Username           string       `bun:"username" json:"username"`
	Email              string       `bun:"email" json:"email"`
	Password           string       `bun:"password" json:"password"`
	Role               string       `bun:"role" json:"role"`
	UniversityID       *string      `bun:"university_id" json:"university_id"`
	University         *University  `bun:"rel:belongs-to,join:university_id=id" json:"university"`
	IsBanned           bool         `bun:"is_banned" json:"is_banned"`
	IsEmailVerified    bool         `bun:"is_email_verified" json:"is_email_verified"`
	HasRateUniversity  bool         `bun:"has_rate_university" json:"has_rate_university"`
	ReputationPoints   int64        `bun:"reputation_points" json:"reputation_points"`
	UniversityRatingID *string      `bun:"-" json:"university_rating_id"`
	CreatedBy          string       `bun:"created_by" json:"created_by"`
	CreatedAt          time.Time    `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy          *string      `json:"updated_by"`
	UpdatedAt          bun.NullTime `json:"updated_at"`
	DeletedBy          *string      `json:"-"`
	DeletedAt          time.Time    `bun:",nullzero,soft_delete" json:"-"`
}

type UserVerifyCode struct {
	bun.BaseModel `bun:"table:user_verify_code,alias:uvc"`

	ID        string
	UserID    string
	Type      string
	Code      string
	Email     string
	IsUsed    bool
	ExpiredAt time.Time `bun:",nullzero"`
	CreatedBy string
	CreatedAt time.Time `bun:",nullzero,default:now()"`
	UpdatedBy *string
	UpdatedAt bun.NullTime
	DeletedBy *string
	DeletedAt time.Time `bun:",nullzero,soft_delete"`
}

type UserDevice struct {
	bun.BaseModel `bun:"table:user_device,alias:ud"`

	ID                   string       `bun:",pk" json:"id"`
	UserID               string       `bun:"user_id" json:"user_id"`
	Brand                *string      `bun:"brand" json:"brand"`
	Type                 *string      `bun:"type" json:"type"`
	Model                *string      `bun:"model" json:"model"`
	NotificationToken    string       `bun:"notification_token" json:"notification_token"`
	IsNotificationActive bool         `bun:"is_notification_active" json:"is_notification_active"`
	CreatedBy            string       `bun:"created_by" json:"created_by"`
	CreatedAt            time.Time    `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy            *string      `bun:"updated_by" json:"updated_by"`
	UpdatedAt            bun.NullTime `bun:"updated_at" json:"updated_at"`
	DeletedBy            *string      `bun:"deleted_by" json:"-"`
	DeletedAt            time.Time    `bun:",nullzero,soft_delete" json:"-"`
}
