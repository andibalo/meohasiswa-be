package model

import (
	"github.com/uptrace/bun"
	"time"
)

type University struct {
	bun.BaseModel `bun:"table:university,alias:uni"`

	ID              string       `bun:",pk" json:"id"`
	Name            string       `bun:"name" json:"name"`
	AbbreviatedName string       `bun:"abbreviated_name" json:"abbreviated_name"`
	ImageURL        string       `bun:"image_url" json:"image_url"`
	CreatedBy       string       `bun:"created_by" json:"created_by"`
	CreatedAt       time.Time    `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy       *string      `json:"updated_by"`
	UpdatedAt       bun.NullTime `json:"updated_at"`
	DeletedBy       *string      `json:"-"`
	DeletedAt       time.Time    `bun:",nullzero,soft_delete" json:"-"`
}
