package request

type RateUniversityReq struct {
	UniversityID              string   `json:"university_id"`
	Title                     string   `json:"title" binding:"required"`
	Content                   string   `json:"content" binding:"required"`
	UniversityMajor           string   `json:"university_major" binding:"required"`
	FacilityRating            int      `json:"facility_rating"  binding:"required"`
	StudentOrganizationRating int      `json:"student_organization_rating"  binding:"required"`
	SocialEnvironmentRating   int      `json:"social_environment_rating"  binding:"required"`
	EducationQualityRating    int      `json:"education_quality_rating"  binding:"required"`
	PriceToValueRating        int      `json:"price_to_value_rating"  binding:"required"`
	Pros                      []string `json:"pros"  binding:"required,min=1,max=3"`
	Cons                      []string `json:"cons"  binding:"required,min=1,max=3"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}

type GetUniversityRatingListReq struct {
	Search string `json:"_q"`
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`

	UserID    string `json:"-"`
	UserEmail string `json:"-"`
}
