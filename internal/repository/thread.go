package repository

import (
	"context"
	"fmt"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/pkg"
	"github.com/uptrace/bun"
	"time"
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

func (r *threadRepository) IncrementLikesCountTx(threadID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread SET like_count = like_count + 1 WHERE id = ?", threadID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DecrementLikesCountTx(threadID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread SET like_count = like_count - 1 WHERE id = ?", threadID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) IncrementDislikesCountTx(threadID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread SET dislike_count = dislike_count + 1 WHERE id = ?", threadID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DecrementDislikesCountTx(threadID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread SET dislike_count = dislike_count - 1 WHERE id = ?", threadID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) SaveThreadActivity(threadActivity *model.ThreadActivity) error {

	_, err := r.db.NewInsert().Model(threadActivity).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) SaveThreadActivityTx(threadActivity *model.ThreadActivity, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(threadActivity).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) GetList(req request.GetThreadListReq) ([]model.Thread, pkg.Pagination, error) {

	//TODO: add trending filter

	var (
		threads    []model.Thread
		nextCursor string
	)

	pagination := pkg.Pagination{}

	query := r.db.NewSelect().
		Column("th.*").
		Model(&threads).
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "username")
		}).
		Relation("User.University").
		Relation("SubThread", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "name", "label_color")
		}).
		Join("JOIN subthread_follower AS stf ON stf.subthread_id = th.subthread_id").
		Where("stf.is_following = TRUE").
		Limit(req.Limit + 1)

	if req.Cursor != "" {
		createdAt, id := pkg.GetCursorData(req.Cursor)
		query.Where("(th.created_at, th.id) <= (?, ?)", createdAt, id)

		query.Order("th.created_at desc", "th.id desc")

	} else {
		query.Order("th.created_at desc")
	}

	err := query.Scan(context.Background())
	if err != nil {
		return threads, pagination, err
	}

	if len(threads) > req.Limit {
		lastThread := threads[len(threads)-1]
		nextCursor = fmt.Sprintf("%s_%s", lastThread.CreatedAt.Format(time.RFC3339Nano), lastThread.ID)
		threads = threads[:req.Limit] // Trim to the requested limit
	}

	pagination.CurrentCursor = req.Cursor
	pagination.NextCursor = nextCursor

	return threads, pagination, nil
}

func (r *threadRepository) GetLastThreadActivityByUserID(threadId string, userId string) (*model.ThreadActivity, error) {
	threadActivity := &model.ThreadActivity{}

	err := r.db.NewSelect().
		Model(threadActivity).
		Where("thread_id = ? AND actor_id = ?", threadId, userId).
		Order("created_at desc").
		Limit(1).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return threadActivity, nil
}
