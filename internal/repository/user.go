package repository

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/uptrace/bun"
)

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) SaveTx(user *model.User, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(user).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Save(user *model.User) error {

	_, err := r.db.NewInsert().Model(user).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserVerifyEmailByIDTx(id string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewUpdate().
		Model(&updateValues).
		TableExpr("user_verify_email").
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserVerifyEmail(email string) (*model.UserVerifyEmail, error) {
	userVerifyEmail := &model.UserVerifyEmail{}

	err := r.db.NewSelect().
		Model(userVerifyEmail).
		Where("email = ?", email).
		Order("created_at desc").
		Limit(1).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return userVerifyEmail, nil
}

func (r *userRepository) GetUserProfileByEmail(email string) (*model.User, error) {
	user := &model.User{}

	err := r.db.NewSelect().
		Model(user).
		ExcludeColumn("password").
		Relation("University").
		Where("email = ?", email).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}

	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) SaveUserVerifyEmailTx(userVerifyEmail *model.UserVerifyEmail, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(userVerifyEmail).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserVerifyEmailByID(id string) (*model.UserVerifyEmail, error) {
	userVerifyEmail := &model.UserVerifyEmail{}

	err := r.db.NewSelect().Model(userVerifyEmail).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return userVerifyEmail, nil
}

func (r *userRepository) SetUserToEmailVerifiedTx(id string, tx bun.Tx) error {
	user := &model.User{}
	user.IsEmailVerified = true

	_, err := tx.NewUpdate().
		Model(user).
		Column("is_email_verified").
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) SetUserVerifyEmailToUsedTx(id string, tx bun.Tx) error {
	userVerifyEmail := &model.UserVerifyEmail{}
	userVerifyEmail.IsUsed = true

	_, err := tx.NewUpdate().
		Model(userVerifyEmail).
		Column("is_used").
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) SetUserHasRateUniversityTx(id string, hru bool, tx bun.Tx) error {
	user := &model.User{}
	user.HasRateUniversity = hru

	_, err := tx.NewUpdate().
		Model(user).
		Column("has_rate_university").
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}
