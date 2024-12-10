package repository

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/request"
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

func (r *userRepository) SaveUserDevice(userDevice *model.UserDevice) error {

	_, err := r.db.NewInsert().Model(userDevice).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserVerifyCodeByID(id string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		TableExpr("user_verify_code").
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserVerifyCodeByIDTx(id string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewUpdate().
		Model(&updateValues).
		TableExpr("user_verify_code").
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserPasswordByUserID(id string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		Table("user").
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserVerifyCodeByEmail(email string, verifyCodeType string) (*model.UserVerifyCode, error) {
	userVerifyCode := &model.UserVerifyCode{}

	err := r.db.NewSelect().
		Model(userVerifyCode).
		Where("email = ?", email).
		Where("type = ?", verifyCodeType).
		Order("created_at desc").
		Limit(1).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return userVerifyCode, nil
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

func (r *userRepository) GetByID(id string) (*model.User, error) {
	user := &model.User{}

	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserDevices(req request.GetUserDevicesReq) ([]model.UserDevice, error) {

	var userDevices = []model.UserDevice{}

	query := r.db.NewSelect().
		Model(&userDevices)

	if req.UserID != "" {
		query.Where("user_id = ?", req.UserID)
	}

	if req.NotificationToken != "" {
		query.Where("notification_token = ?", req.NotificationToken)
	}

	err := query.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return userDevices, nil
}

func (r *userRepository) SaveUserVerifyCode(userVerifyCode *model.UserVerifyCode) error {

	_, err := r.db.NewInsert().Model(userVerifyCode).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) SaveUserVerifyCodeTx(userVerifyCode *model.UserVerifyCode, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(userVerifyCode).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserVerifyCodeByID(id string, verifyCodeType string) (*model.UserVerifyCode, error) {
	userVerifyCode := &model.UserVerifyCode{}

	err := r.db.NewSelect().
		Model(userVerifyCode).
		Where("id = ?", id).
		Where("type = ?", verifyCodeType).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return userVerifyCode, nil
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

func (r *userRepository) SetUserVerifyCodeToUsedTx(id string, tx bun.Tx) error {
	userVerifyCode := &model.UserVerifyCode{}
	userVerifyCode.IsUsed = true

	_, err := tx.NewUpdate().
		Model(userVerifyCode).
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

func (r *userRepository) UpdateUser(id string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		Table("user").
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUserTx(id string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewUpdate().
		Model(&updateValues).
		Table("user").
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) IncrementUserReputationPointsTx(id string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewRaw(`update
									"user"
								set 
									reputation_points = reputation_points + ?, 
									updated_at = now(), 
									updated_by = ?
								where
								id = ?`, updateValues["reputation_points"], updateValues["updated_by"], id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DecrementUserReputationPointsTx(id string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewRaw(`update
									"user"
								set 
									reputation_points = reputation_points - ?, 
									updated_at = now(), 
									updated_by = ?
								where
								id = ?`, updateValues["reputation_points"], updateValues["updated_by"], id).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
