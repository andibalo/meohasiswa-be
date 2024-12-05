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

type subThreadRepository struct {
	db *bun.DB
}

func NewSubThreadRepository(db *bun.DB) SubThreadRepository {
	return &subThreadRepository{
		db: db,
	}
}

func (r *subThreadRepository) GetList(req request.GetSubThreadListReq) ([]model.SubThread, pkg.Pagination, error) {

	var (
		subThreads = []model.SubThread{}
		nextCursor string
	)

	pagination := pkg.Pagination{}

	query := r.db.NewSelect().
		Column("st.*").
		Model(&subThreads).
		Limit(req.Limit + 1)

	if !req.IncludeUniversitySubThread {
		query.Where("st.is_university_subthread = FALSE")
	}

	if req.IsFollowing {
		query.Join("JOIN subthread_follower AS stf ON stf.subthread_id = st.id").
			Where("stf.user_id = ?", req.UserID).
			Where("stf.is_following = TRUE")
	}

	if req.ShouldExcludeFollowing {
		query.Join("LEFT JOIN subthread_follower AS stf ON stf.subthread_id = st.id AND stf.user_id = ?", req.UserID).
			Where("stf.user_id IS NULL OR stf.is_following = FALSE")
	}

	if req.Search != "" {
		searchCols := []string{
			"st.name",
			"st.description",
		}

		query.Where("CONCAT("+strings.Join(searchCols, ", ")+") ILIKE ?", "%"+req.Search+"%")
	}

	if req.Cursor != "" {
		createdAt, id := pkg.GetCursorData(req.Cursor)
		query.Where("(st.created_at, st.id) <= (?, ?)", createdAt, id)

		query.Order("st.created_at desc", "st.id desc")

	} else {
		query.Order("st.created_at desc")
	}

	err := query.Scan(context.Background())
	if err != nil {
		return subThreads, pagination, err
	}

	if len(subThreads) > req.Limit {
		lastSubThread := subThreads[len(subThreads)-1]
		nextCursor = fmt.Sprintf("%s_%s", lastSubThread.CreatedAt.Format(time.RFC3339Nano), lastSubThread.ID)
		subThreads = subThreads[:req.Limit] // Trim to the requested limit
	}

	pagination.CurrentCursor = req.Cursor
	pagination.NextCursor = nextCursor

	return subThreads, pagination, nil
}

func (r *subThreadRepository) IncrementFollowersCountTx(subThreadID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE subthread SET followers_count = followers_count + 1 WHERE id = ?", subThreadID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *subThreadRepository) DecrementFollowersCountTx(subThreadID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE subthread SET followers_count = followers_count - 1 WHERE id = ?", subThreadID).
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

func (r *subThreadRepository) SaveSubThreadFollowerTx(subThreadFollower *model.SubThreadFollower, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(subThreadFollower).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *subThreadRepository) DeleteSubThreadFollowerTx(subThreadFollower *model.SubThreadFollower, tx bun.Tx) error {

	_, err := tx.NewDelete().
		Model(subThreadFollower).
		Where("user_id = ? and subthread_id = ?", subThreadFollower.UserID, subThreadFollower.SubThreadID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *subThreadRepository) GetSubThreadFollowerByUserIDAndSubThreadID(userID string, subThreadID string) (*model.SubThreadFollower, error) {
	subThreadFollower := &model.SubThreadFollower{}

	err := r.db.NewSelect().
		Model(subThreadFollower).
		Where("user_id = ? and subthread_id = ?", userID, subThreadID).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return subThreadFollower, nil
}

func (r *subThreadRepository) UpdateSubThreadFollowerIsFollowingTx(id string, isFollowing bool, tx bun.Tx) error {
	subThreadFollower := &model.SubThreadFollower{}
	subThreadFollower.IsFollowing = isFollowing

	_, err := tx.NewUpdate().
		Model(subThreadFollower).
		Column("is_following").
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *subThreadRepository) UpdateByID(subThreadID string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		TableExpr("subthread").
		Where("id = ?", subThreadID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *subThreadRepository) DeleteByID(subThreadID string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		TableExpr("subthread").
		Where("id = ?", subThreadID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
