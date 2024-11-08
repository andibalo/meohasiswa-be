package repository

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/uptrace/bun"
)

type threadRepository struct {
	db *bun.DB
}

func NewThreadRepository(db *bun.DB) ThreadRepository {
	return &threadRepository{
		db: db,
	}
}

func (r *threadRepository) Save(thread *model.Thread) error {

	_, err := r.db.NewInsert().Model(thread).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
