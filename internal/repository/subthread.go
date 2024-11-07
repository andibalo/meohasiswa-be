package repository

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/uptrace/bun"
)

type subThreadRepository struct {
	db *bun.DB
}

func NewSubThreadRepository(db *bun.DB) SubThreadRepository {
	return &subThreadRepository{
		db: db,
	}
}

func (r *subThreadRepository) IncrementFollowersCountTx(subThreadID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE subthread SET followers_count = followers_count + 1 WHERE id = ?", subThreadID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *subThreadRepository) Save(user *model.SubThread) error {

	_, err := r.db.NewInsert().Model(user).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *subThreadRepository) GetByName(name string) (*model.SubThread, error) {
	subThread := &model.SubThread{}

	err := r.db.NewSelect().Model(subThread).Where("name = ?", name).Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return subThread, nil
}

func (r *subThreadRepository) GetByID(subThreadID string) (*model.SubThread, error) {
	subThread := &model.SubThread{}

	err := r.db.NewSelect().Model(subThread).Where("id = ?", subThreadID).Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return subThread, nil
}

func (r *subThreadRepository) SaveTx(user *model.SubThread, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(user).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *subThreadRepository) SaveSubThreadFollowerTx(user *model.SubThreadFollower, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(user).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
