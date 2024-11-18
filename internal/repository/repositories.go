package repository

import (
	"context"
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/request"
	"github.com/andibalo/meowhasiswa-be/pkg"
	"github.com/uptrace/bun"
)

type UserRepository interface {
	Save(user *model.User) error
	SaveTx(user *model.User, tx bun.Tx) error
	GetUserProfileByEmail(email string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	SaveUserVerifyEmailTx(userVerifyEmail *model.UserVerifyEmail, tx bun.Tx) error
	SetUserToEmailVerifiedTx(id string, tx bun.Tx) error
	GetUserVerifyEmailByID(id string) (*model.UserVerifyEmail, error)
	SetUserVerifyEmailToUsedTx(id string, tx bun.Tx) error
	SetUserHasRateUniversityTx(id string, hru bool, tx bun.Tx) error
}

type SubThreadRepository interface {
	Save(subThread *model.SubThread) error
	SaveTx(subThread *model.SubThread, tx bun.Tx) error
	GetList(req request.GetSubThreadListReq) ([]model.SubThread, pkg.Pagination, error)
	GetByID(subThreadID string) (*model.SubThread, error)
	GetByName(name string) (*model.SubThread, error)
	GetSubThreadFollowerByUserIDAndSubThreadID(userID string, subThreadID string) (*model.SubThreadFollower, error)
	IncrementFollowersCountTx(subThreadID string, tx bun.Tx) error
	DecrementFollowersCountTx(subThreadID string, tx bun.Tx) error
	SaveSubThreadFollowerTx(subthreadFollower *model.SubThreadFollower, tx bun.Tx) error
	DeleteSubThreadFollowerTx(subThreadFollower *model.SubThreadFollower, tx bun.Tx) error
	UpdateSubThreadFollowerIsFollowingTx(id string, isFollowing bool, tx bun.Tx) error
}

type ThreadRepository interface {
	Save(thread *model.Thread) error
	GetList(req request.GetThreadListReq) ([]model.Thread, pkg.Pagination, error)
	GetByID(id string) (model.Thread, error)
	GetLastThreadActivityByUserID(threadId string, userId string) (*model.ThreadActivity, error)
	SaveThreadActivity(threadActivity *model.ThreadActivity) error
	SaveThreadActivityTx(threadActivity *model.ThreadActivity, tx bun.Tx) error
	SaveThreadCommentTx(threadComment *model.ThreadComment, tx bun.Tx) error
	SaveCommentReplyTx(threadCommentReply *model.ThreadCommentReply, tx bun.Tx) error
	IncrementLikesCountTx(threadID string, tx bun.Tx) error
	DecrementLikesCountTx(threadID string, tx bun.Tx) error
	IncrementDislikesCountTx(threadID string, tx bun.Tx) error
	DecrementDislikesCountTx(threadID string, tx bun.Tx) error
	IncrementCommentsCountTx(threadID string, tx bun.Tx) error
	IncrementCommentReplyCountTx(commentID string, tx bun.Tx) error
}

type UniversityRepository interface {
	GetList(req request.GetUniversityRatingListReq) ([]model.UniversityRating, pkg.Pagination, error)
	GetUniversityRatingByUserIDAndUniversityID(userID string, universityID string) (*model.UniversityRating, error)
	Save(university *model.University) error
	SaveUniversityRatingTx(universityRating *model.UniversityRating, tx bun.Tx) error
	BulkSaveUniversityRatingPointsTx(urp []model.UniversityRatingPoints, tx bun.Tx) error
}

type FileRepository interface {
	Upload(ctx context.Context, uploadFileData model.UploadFileDTO) (model.UploadFileOutputDTO, error)
}
