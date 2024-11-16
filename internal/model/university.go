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

type UniversityRating struct {
	bun.BaseModel `bun:"table:university_rating,alias:unir"`

	ID                        string                   `bun:",pk" json:"id"`
	UserID                    string                   `bun:"user_id" json:"user_id"`
	User                      User                     `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	UniversityID              string                   `bun:"university_id" json:"university_id"`
	University                University               `bun:"rel:belongs-to,join:university_id=id" json:"university"`
	Title                     string                   `bun:"title" json:"title"`
	Content                   string                   `bun:"content" json:"content"`
	UniversityMajor           string                   `bun:"university_major" json:"university_major"`
	FacilityRating            int                      `bun:"facility_rating" json:"facility_rating"`
	StudentOrganizationRating int                      `bun:"student_organization_rating" json:"student_organization_rating"`
	SocialEnvironmentRating   int                      `bun:"social_environment_rating" json:"social_environment_rating"`
	EducationQualityRating    int                      `bun:"education_quality_rating" json:"education_quality_rating"`
	PriceToValueRating        int                      `bun:"price_to_value_rating" json:"price_to_value_rating"`
	OverallRating             float64                  `bun:"overall_rating" json:"overall_rating"`
	UniversityRatingPoints    []UniversityRatingPoints `bun:"rel:has-many,join:id=university_rating_id" json:"university_rating_points"`
	CreatedBy                 string                   `bun:"created_by" json:"created_by"`
	CreatedAt                 time.Time                `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy                 string                   `bun:"updated_by" json:"updated_by"`
	UpdatedAt                 time.Time                `bun:",nullzero,default:now()" json:"updated_at"`
	DeletedBy                 *string                  `bun:"deleted_by" json:"-"`
	DeletedAt                 time.Time                `bun:",nullzero,soft_delete" json:"-"`
}

type UniversityRatingPoints struct {
	bun.BaseModel `bun:"table:university_rating_point,alias:unirp"`

	ID                 string    `bun:",pk" json:"id"`
	UniversityRatingID string    `bun:"university_rating_id" json:"university_rating_id"`
	Type               string    `bun:"type" json:"type"`
	Content            string    `bun:"content" json:"content"`
	CreatedBy          string    `bun:"created_by" json:"created_by"`
	CreatedAt          time.Time `bun:",nullzero,default:now()" json:"created_at"`
	UpdatedBy          string    `json:"updated_by"`
	UpdatedAt          time.Time `bun:",nullzero,default:now()" json:"updated_at"`
}
