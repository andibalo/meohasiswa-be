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
