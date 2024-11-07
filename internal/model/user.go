package model

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:user,alias:u"`

	ID               string    `bun:",pk"`
	Username         string    `bun:"username"`
	Email            string    `bun:"email"`
	Password         string    `bun:"password"`
	Role             string    `bun:"role"`
	IsBanned         bool      `bun:"is_banned"`
	IsEmailVerified  bool      `bun:"is_email_verified"`
	ReputationPoints int64     `bun:"reputation_points"`
	CreatedBy        string    `bun:"created_by"`
	CreatedAt        time.Time `bun:",nullzero,default:now()"`
	UpdatedBy        *string
	UpdatedAt        bun.NullTime
	DeletedBy        *string
	DeletedAt        time.Time `bun:",nullzero,soft_delete"`
}

type UserVerifyEmail struct {
	bun.BaseModel `bun:"table:user_verify_email,alias:uve"`

	ID         string
	UserID     string
	SecretCode string
	Email      string
	IsUsed     bool
	ExpiredAt  time.Time `bun:",nullzero"`
	CreatedBy  string
	CreatedAt  time.Time `bun:",nullzero,default:now()"`
	UpdatedBy  *string
	UpdatedAt  bun.NullTime
	DeletedBy  *string
	DeletedAt  time.Time `bun:",nullzero,soft_delete"`
}

type UserDevice struct {
	bun.BaseModel `bun:"table:user_device,alias:ud"`

	ID                   string
	DeviceType           string
	DeviceID             string
	UserID               string
	NotificationToken    string
	IsNotificationActive bool
	CreatedBy            string
	CreatedAt            time.Time `bun:",nullzero,default:now()"`
	UpdatedBy            *string
	UpdatedAt            bun.NullTime
	DeletedBy            *string
	DeletedAt            time.Time `bun:",nullzero,soft_delete"`
}
