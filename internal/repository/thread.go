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

func (r *threadRepository) UpdateByID(threadID string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		TableExpr("thread").
		Where("id = ?", threadID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DeleteByID(threadID string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		TableExpr("thread").
		Where("id = ?", threadID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) IncrementCommentsCountTx(threadID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread SET comment_count = comment_count + 1 WHERE id = ?", threadID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DecrementCommentsCountTx(threadID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread SET comment_count = comment_count - 1 WHERE id = ?", threadID).
		Exec(context.Background())

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

func (r *threadRepository) UpdateThreadActivityTx(threadID string, actorID string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewUpdate().
		Model(&updateValues).
		TableExpr("thread_activity").
		Where("thread_id = ?", threadID).
		Where("actor_id = ?", actorID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) UpdateThreadCommentActivityTx(threadCommentID string, actorID string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewUpdate().
		Model(&updateValues).
		TableExpr("thread_comment_activity").
		Where("thread_comment_id = ?", threadCommentID).
		Where("actor_id = ?", actorID).
		Where("thread_comment_reply_id IS NULL").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) UpdateThreadCommentActivityReplyTx(threadCommentReplyID string, actorID string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewUpdate().
		Model(&updateValues).
		TableExpr("thread_comment_activity").
		Where("thread_comment_reply_id = ?", threadCommentReplyID).
		Where("actor_id = ?", actorID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) SaveThreadCommentTx(threadComment *model.ThreadComment, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(threadComment).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) GetByID(id string) (model.Thread, error) {

	var (
		thread model.Thread
	)

	err := r.db.NewSelect().
		Column("th.*").
		Model(&thread).
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "username")
		}).
		Relation("User.University").
		Relation("SubThread", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "name", "label_color")
		}).
		Where("th.id = ?", id).
		Scan(context.Background())

	if err != nil {
		return thread, err
	}

	return thread, nil
}

func (r *threadRepository) GetThreadSubscriptionByUserAndThreadID(userID string, threadID string) (model.ThreadSubscription, error) {

	threadSubscription := model.ThreadSubscription{}

	err := r.db.NewSelect().
		Model(&threadSubscription).
		Where("user_id = ? and thread_id = ?", userID, threadID).
		Scan(context.Background())

	if err != nil {
		return threadSubscription, err
	}

	return threadSubscription, nil
}

func (r *threadRepository) GetList(req request.GetThreadListReq) ([]model.Thread, pkg.Pagination, error) {

	var (
		threads    []model.Thread
		nextCursor string
	)

	pagination := pkg.Pagination{}

	trendingScoreSubQuery := r.db.NewSelect().TableExpr("thread").
		ColumnExpr(`id,
							ROUND((
									(like_count * 1.5) +
									(dislike_count * 1.2) +
									(comment_count * 2)
								) * EXP(EXTRACT(EPOCH FROM (NOW() - created_at)) / -172800.0), 2) AS trending_score`)

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
		Limit(req.Limit + 1)

	if req.IncludeUserActivity {
		query.ColumnExpr("ta.action as thread_action")
		query.Join("LEFT JOIN thread_activity AS ta ON ta.thread_id = th.id AND ta.actor_id = ?", req.UserID)
	}

	if req.IsTrending {
		query.Column("ts.trending_score")
		query.Join("LEFT JOIN (?) AS ts ON (ts.id = th.id)", trendingScoreSubQuery)
	}

	if req.IsUserFollowing {
		query.Join("JOIN subthread_follower AS stf ON stf.subthread_id = th.subthread_id AND stf.user_id = ?", req.UserID)
		query.Where("stf.is_following = TRUE")
	}

	if req.UserIDParam != "" {
		query.Where("th.user_id = ?", req.UserIDParam)
	}

	if req.Search != "" {
		searchCols := []string{
			"th.title",
			"th.content",
			"th.content_summary",
		}

		query.Where("CONCAT("+strings.Join(searchCols, ", ")+") ILIKE ?", "%"+req.Search+"%")
	}

	if req.Cursor != "" {
		if req.IsTrending {
			trendingScore, _ := pkg.GetCursorData(req.Cursor)

			query.Where("ts.trending_score <= ?", trendingScore)
			query.Order("ts.trending_score desc")

		} else {

			createdAt, id := pkg.GetCursorData(req.Cursor)

			query.Where("(th.created_at, th.id) <= (?, ?)", createdAt, id)

			query.Order("th.created_at desc", "th.id desc")
		}
	} else {
		if req.IsTrending {
			query.Order("ts.trending_score desc")
		} else {
			query.Order("th.created_at desc")
		}
	}

	err := query.Scan(context.Background())
	if err != nil {
		return threads, pagination, err
	}

	if len(threads) > req.Limit {
		lastThread := threads[len(threads)-1]

		if req.IsTrending {
			nextCursor = fmt.Sprintf("%.2f_%s", lastThread.TrendingScore, lastThread.ID)
		} else {
			nextCursor = fmt.Sprintf("%s_%s", lastThread.CreatedAt.Format(time.RFC3339Nano), lastThread.ID)
		}

		threads = threads[:req.Limit] // Trim to the requested limit
	}

	pagination.CurrentCursor = req.Cursor
	pagination.NextCursor = nextCursor

	return threads, pagination, nil
}

func (r *threadRepository) GetThreadCommentReplyByID(id string) (model.ThreadCommentReply, error) {

	var (
		tcr model.ThreadCommentReply
	)

	err := r.db.NewSelect().
		Model(&tcr).
		Where("thcr.id = ?", id).
		Scan(context.Background())

	if err != nil {
		return tcr, err
	}

	return tcr, nil
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

func (r *threadRepository) IncrementCommentReplyCountTx(commentID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread_comment SET reply_count = reply_count + 1 WHERE id = ?", commentID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) SaveCommentReplyTx(threadCommentReply *model.ThreadCommentReply, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(threadCommentReply).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) IncrementCommentLikesCountTx(threadCommentID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread_comment SET like_count = like_count + 1 WHERE id = ?", threadCommentID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DecrementCommentLikesCountTx(threadCommentID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread_comment SET like_count = like_count - 1 WHERE id = ?", threadCommentID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) IncrementCommentDislikesCountTx(threadCommentID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread_comment SET dislike_count = dislike_count + 1 WHERE id = ?", threadCommentID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DecrementCommentDislikesCountTx(threadCommentID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread_comment SET dislike_count = dislike_count - 1 WHERE id = ?", threadCommentID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) IncrementCommentReplyLikesCountTx(threadCommentReplyID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread_comment_reply SET like_count = like_count + 1 WHERE id = ?", threadCommentReplyID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DecrementCommentReplyLikesCountTx(threadCommentID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread_comment_reply SET like_count = like_count - 1 WHERE id = ?", threadCommentID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) IncrementCommentReplyDislikesCountTx(threadCommentID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread_comment_reply SET dislike_count = dislike_count + 1 WHERE id = ?", threadCommentID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DecrementCommentReplyDislikesCountTx(threadCommentID string, tx bun.Tx) error {

	_, err := tx.NewRaw("UPDATE thread_comment_reply SET dislike_count = dislike_count - 1 WHERE id = ?", threadCommentID).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) SaveThreadCommentActivityTx(tca *model.ThreadCommentActivity, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(tca).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) GetLastThreadCommentActivityByUserID(threadId string, commentId string, userId string) (*model.ThreadCommentActivity, error) {
	threadCommentActivity := &model.ThreadCommentActivity{}

	err := r.db.NewSelect().
		Model(threadCommentActivity).
		Where("thread_id = ? AND thread_comment_id = ? AND actor_id = ? AND thread_comment_reply_id IS NULL", threadId, commentId, userId).
		Order("created_at desc").
		Limit(1).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return threadCommentActivity, nil
}

func (r *threadRepository) GetLastThreadCommentActivityReplyByUserID(threadId string, commentReplyId string, userId string) (*model.ThreadCommentActivity, error) {
	threadCommentActivity := &model.ThreadCommentActivity{}

	err := r.db.NewSelect().
		Model(threadCommentActivity).
		Where("thread_id = ? AND thread_comment_reply_id = ? AND actor_id = ?", threadId, commentReplyId, userId).
		Order("created_at desc").
		Limit(1).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return threadCommentActivity, nil
}

func (r *threadRepository) GetThreadCommentsByThreadID(threadId string, userId string) ([]model.ThreadComment, error) {
	var (
		threadComments []model.ThreadComment
	)

	err := r.db.NewSelect().
		Column("thc.*").
		ColumnExpr("tca.action as comment_action").
		Model(&threadComments).
		Join("LEFT JOIN thread_comment_activity AS tca ON tca.thread_comment_id = thc.id AND tca.actor_id = ? AND tca.thread_comment_reply_id IS NULL", userId).
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "username")
		}).
		Relation("User.University", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "name", "abbreviated_name", "image_url")
		}).
		Relation("Replies", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ColumnExpr("thcr.*, tca.action as comment_reply_action").
				Join("LEFT JOIN thread_comment_activity AS tca ON tca.thread_comment_reply_id = thcr.id AND tca.actor_id = ?", userId).
				Order("thcr.created_at desc")
		}).
		Relation("Replies.User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "username")
		}).
		Relation("Replies.User.University", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "name", "abbreviated_name", "image_url")
		}).
		Where("thc.thread_id = ?", threadId).
		Order("thc.created_at desc").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return threadComments, nil
}

func (r *threadRepository) GetThreadCommentByID(id string) (model.ThreadComment, error) {

	var (
		tc model.ThreadComment
	)

	err := r.db.NewSelect().
		Model(&tc).
		Where("thc.id = ?", id).
		Scan(context.Background())

	if err != nil {
		return tc, err
	}

	return tc, nil
}

func (r *threadRepository) DeleteThreadCommentByID(threadCommentID string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		TableExpr("thread_comment").
		Where("id = ?", threadCommentID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DeleteThreadCommentByIDTx(threadCommentID string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewUpdate().
		Model(&updateValues).
		TableExpr("thread_comment").
		Where("id = ?", threadCommentID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) UpdateThreadCommentByID(threadCommentID string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		TableExpr("thread_comment").
		Where("id = ?", threadCommentID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DeleteThreadCommentReplyByID(threadCommentReplyID string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		TableExpr("thread_comment_reply").
		Where("id = ?", threadCommentReplyID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) DeleteThreadCommentReplyByIDTx(threadCommentReplyID string, updateValues map[string]interface{}, tx bun.Tx) error {

	_, err := tx.NewUpdate().
		Model(&updateValues).
		TableExpr("thread_comment_reply").
		Where("id = ?", threadCommentReplyID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) UpdateThreadCommentReplyByID(threadCommentReplyID string, updateValues map[string]interface{}) error {

	_, err := r.db.NewUpdate().
		Model(&updateValues).
		TableExpr("thread_comment_reply").
		Where("id = ?", threadCommentReplyID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) SaveThreadSubscription(threadSubscription *model.ThreadSubscription) error {

	_, err := r.db.NewInsert().Model(threadSubscription).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *threadRepository) UpdateThreadSubscriptionIsSubscribed(id string, isSubscribed bool) error {
	threadSubscription := &model.ThreadSubscription{}
	threadSubscription.IsSubscribed = isSubscribed

	_, err := r.db.NewUpdate().
		Model(threadSubscription).
		Column("is_subscribed").
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}
