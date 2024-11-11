package repository

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/uptrace/bun"
)

type universityRepository struct {
	db *bun.DB
}

func NewUniversityRepository(db *bun.DB) UniversityRepository {
	return &universityRepository{
		db: db,
	}
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
