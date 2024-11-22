package repository

import (
	"context"
	"fmt"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/pkg"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

type universityRepository struct {
	db *bun.DB
}

func NewUniversityRepository(db *bun.DB) UniversityRepository {
	return &universityRepository{
		db: db,
	}
}

func (r *universityRepository) GetList(req request.GetUniversityRatingListReq) ([]model.UniversityRating, pkg.Pagination, error) {

	var (
		uniRatings []model.UniversityRating
		nextCursor string
	)

	pagination := pkg.Pagination{}

	query := r.db.NewSelect().
		Column("unir.*").
		Model(&uniRatings).
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "username")
		}).
		Relation("University", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "name", "abbreviated_name", "image_url")
		}).
		Relation("UniversityRatingPoints").
		Limit(req.Limit + 1)

	if req.Search != "" {
		query.Join("JOIN university AS uni ON uni.id = unir.university_id")

		searchCols := []string{
			"uni.name",
			"uni.abbreviated_name",
			"unir.title",
			"unir.content",
			"unir.university_major",
		}

		query.Where("CONCAT("+strings.Join(searchCols, ", ")+") ILIKE ?", "%"+req.Search+"%")
	}

	if req.Cursor != "" {
		createdAt, id := pkg.GetCursorData(req.Cursor)
		query.Where("(unir.created_at, unir.id) <= (?, ?)", createdAt, id)

		query.Order("unir.created_at desc", "unir.id desc")

	} else {
		query.Order("unir.created_at desc")
	}

	err := query.Scan(context.Background())
	if err != nil {
		return uniRatings, pagination, err
	}

	if len(uniRatings) > req.Limit {
		lastUniRating := uniRatings[len(uniRatings)-1]
		nextCursor = fmt.Sprintf("%s_%s", lastUniRating.CreatedAt.Format(time.RFC3339Nano), lastUniRating.ID)
		uniRatings = uniRatings[:req.Limit] // Trim to the requested limit
	}

	pagination.CurrentCursor = req.Cursor
	pagination.NextCursor = nextCursor

	return uniRatings, pagination, nil
}

func (r *universityRepository) Save(university *model.University) error {

	_, err := r.db.NewInsert().Model(university).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *universityRepository) GetUniversityRatingByUserIDAndUniversityID(userID string, universityID string) (*model.UniversityRating, error) {
	universityRating := &model.UniversityRating{}

	err := r.db.NewSelect().Model(universityRating).Where("user_id = ? AND university_id = ?", userID, universityID).Scan(context.Background())
	if err != nil {
		return universityRating, err
	}

	return universityRating, nil
}

func (r *universityRepository) SaveUniversityRatingTx(universityRating *model.UniversityRating, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(universityRating).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *universityRepository) BulkSaveUniversityRatingPointsTx(urp []model.UniversityRatingPoints, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(&urp).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
