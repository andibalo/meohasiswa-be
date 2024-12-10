package response

import (
	"github.com/uptrace/bun"
	"time"
)

type GetUniversityRatingListResponse struct {
	Data []UniversityRatingListData `json:"university_ratings"`
	Meta PaginationMeta             `json:"meta"`
}

type UniversityRatingListData struct {
	ID                        string       `json:"id"`
	UserID                    string       `json:"user_id"`
	UserName                  string       `json:"username"`
	UniversityID              string       `bun:"university_id" json:"university_id"`
	UniversityName            string       `json:"university_name"`
	UniversityAbbreviatedName string       `json:"university_abbreviated_name"`
	UniversityImageURL        string       `json:"university_image_url"`
	Title                     string       `json:"title"`
	Content                   string       `json:"content"`
	UniversityMajor           string       `json:"university_major"`
	FacilityRating            int          `json:"facility_rating"`
	StudentOrganizationRating int          `json:"student_organization_rating"`
	SocialEnvironmentRating   int          `json:"social_environment_rating"`
	EducationQualityRating    int          `json:"education_quality_rating"`
	PriceToValueRating        int          `json:"price_to_value_rating"`
	OverallRating             float64      `json:"overall_rating"`
	Pros                      []string     `json:"pros"`
	Cons                      []string     `json:"cons"`
	CreatedBy                 string       `json:"created_by"`
	CreatedAt                 time.Time    `json:"created_at"`
	UpdatedBy                 *string      `json:"updated_by"`
	UpdatedAt                 bun.NullTime `json:"updated_at"`
}

type UniversityRatingDetailData struct {
	ID                        string       `json:"id"`
	UserID                    string       `json:"user_id"`
	UserName                  string       `json:"username"`
	UniversityID              string       `bun:"university_id" json:"university_id"`
	UniversityName            string       `json:"university_name"`
	UniversityAbbreviatedName string       `json:"university_abbreviated_name"`
	UniversityImageURL        string       `json:"university_image_url"`
	Title                     string       `json:"title"`
	Content                   string       `json:"content"`
	UniversityMajor           string       `json:"university_major"`
	FacilityRating            int          `json:"facility_rating"`
	StudentOrganizationRating int          `json:"student_organization_rating"`
	SocialEnvironmentRating   int          `json:"social_environment_rating"`
	EducationQualityRating    int          `json:"education_quality_rating"`
	PriceToValueRating        int          `json:"price_to_value_rating"`
	OverallRating             float64      `json:"overall_rating"`
	Pros                      []string     `json:"pros"`
	Cons                      []string     `json:"cons"`
	CreatedBy                 string       `json:"created_by"`
	CreatedAt                 time.Time    `json:"created_at"`
	UpdatedBy                 *string      `json:"updated_by"`
	UpdatedAt                 bun.NullTime `json:"updated_at"`
}
