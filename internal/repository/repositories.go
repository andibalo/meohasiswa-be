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
	SaveUserDevice(userDevice *model.UserDevice) error
	GetUserProfileByEmail(email string) (*model.User, error)
	GetByID(id string) (*model.User, error)
	GetUserDevices(req request.GetUserDevicesReq) ([]model.UserDevice, error)
	GetByEmail(email string) (*model.User, error)
	SaveUserVerifyEmailTx(userVerifyEmail *model.UserVerifyEmail, tx bun.Tx) error
	SetUserToEmailVerifiedTx(id string, tx bun.Tx) error
	GetUserVerifyEmailByID(id string) (*model.UserVerifyEmail, error)
	SetUserVerifyEmailToUsedTx(id string, tx bun.Tx) error
	SetUserHasRateUniversityTx(id string, hru bool, tx bun.Tx) error
	GetUserVerifyEmail(email string) (*model.UserVerifyEmail, error)
	UpdateUserVerifyEmailByIDTx(id string, updateValues map[string]interface{}, tx bun.Tx) error
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
	DeleteByID(subThreadID string, updateValues map[string]interface{}) error
	UpdateByID(subThreadID string, updateValues map[string]interface{}) error
}

type ThreadRepository interface {
	Save(thread *model.Thread) error
	UpdateByID(threadID string, updateValues map[string]interface{}) error
	DeleteByID(threadID string, updateValues map[string]interface{}) error
	GetList(req request.GetThreadListReq) ([]model.Thread, pkg.Pagination, error)
	GetByID(id string) (model.Thread, error)
	GetThreadCommentByID(id string) (model.ThreadComment, error)
	DeleteThreadCommentByID(threadCommentID string, updateValues map[string]interface{}) error
	UpdateThreadCommentByID(threadCommentID string, updateValues map[string]interface{}) error
	GetThreadCommentReplyByID(id string) (model.ThreadCommentReply, error)
	DeleteThreadCommentReplyByID(threadCommentReplyID string, updateValues map[string]interface{}) error
	UpdateThreadCommentReplyByID(threadCommentReplyID string, updateValues map[string]interface{}) error
	GetLastThreadActivityByUserID(threadId string, userId string) (*model.ThreadActivity, error)
	GetLastThreadCommentActivityByUserID(threadId string, commentId string, userId string) (*model.ThreadCommentActivity, error)
	GetLastThreadCommentActivityReplyByUserID(threadId string, commentReplyId string, userId string) (*model.ThreadCommentActivity, error)
	GetThreadCommentsByThreadID(threadId string, userId string) ([]model.ThreadComment, error)
	SaveThreadActivity(threadActivity *model.ThreadActivity) error
	SaveThreadActivityTx(threadActivity *model.ThreadActivity, tx bun.Tx) error
	UpdateThreadActivityTx(threadActivityID string, actorID string, updateValues map[string]interface{}, tx bun.Tx) error
	UpdateThreadCommentActivityTx(threadCommentID string, actorID string, updateValues map[string]interface{}, tx bun.Tx) error
	UpdateThreadCommentActivityReplyTx(threadCommentReplyID string, actorID string, updateValues map[string]interface{}, tx bun.Tx) error
	SaveThreadCommentTx(threadComment *model.ThreadComment, tx bun.Tx) error
	SaveCommentReplyTx(threadCommentReply *model.ThreadCommentReply, tx bun.Tx) error
	IncrementLikesCountTx(threadID string, tx bun.Tx) error
	DecrementLikesCountTx(threadID string, tx bun.Tx) error
	IncrementDislikesCountTx(threadID string, tx bun.Tx) error
	DecrementDislikesCountTx(threadID string, tx bun.Tx) error
	IncrementCommentsCountTx(threadID string, tx bun.Tx) error
	IncrementCommentReplyCountTx(commentID string, tx bun.Tx) error
	IncrementCommentLikesCountTx(threadCommentID string, tx bun.Tx) error
	DecrementCommentLikesCountTx(threadCommentID string, tx bun.Tx) error
	IncrementCommentDislikesCountTx(threadCommentID string, tx bun.Tx) error
	DecrementCommentDislikesCountTx(threadCommentID string, tx bun.Tx) error
	IncrementCommentReplyLikesCountTx(threadCommentReplyID string, tx bun.Tx) error
	DecrementCommentReplyLikesCountTx(threadCommentReplyID string, tx bun.Tx) error
	IncrementCommentReplyDislikesCountTx(threadCommentReplyID string, tx bun.Tx) error
	DecrementCommentReplyDislikesCountTx(threadCommentReplyID string, tx bun.Tx) error
	SaveThreadCommentActivityTx(tca *model.ThreadCommentActivity, tx bun.Tx) error
}

type UniversityRepository interface {
	GetUniversityRatingByID(id string) (model.UniversityRating, error)
	GetList(req request.GetUniversityRatingListReq) ([]model.UniversityRating, pkg.Pagination, error)
	GetUniversityRatingByUserIDAndUniversityID(userID string, universityID string) (*model.UniversityRating, error)
	Save(university *model.University) error
	SaveUniversityRatingTx(universityRating *model.UniversityRating, tx bun.Tx) error
	UpdateUniversityRatingByIDTx(universityRatingID string, updateValues map[string]interface{}, tx bun.Tx) error
	BulkSaveUniversityRatingPointsTx(urp []model.UniversityRatingPoints, tx bun.Tx) error
	UpdateUniversityRatingPointByIDTx(universityRatingPointID string, updateValues map[string]interface{}, tx bun.Tx) error
	DeleteUniversityRatingPointsTx(universityRatingID string, tx bun.Tx) error
}

type FileRepository interface {
	Upload(ctx context.Context, uploadFileData model.UploadFileDTO) (model.UploadFileOutputDTO, error)
}
