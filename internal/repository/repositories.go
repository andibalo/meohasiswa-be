package repository

import (
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/uptrace/bun"
)

type UserRepository interface {
	Save(user *model.User) error
	SaveTx(user *model.User, tx bun.Tx) error
	GetByEmail(email string) (*model.User, error)
	SaveUserVerifyEmailTx(userVerifyEmail *model.UserVerifyEmail, tx bun.Tx) error
	SetUserToEmailVerifiedTx(id string, tx bun.Tx) error
	GetUserVerifyEmailByID(id string) (*model.UserVerifyEmail, error)
	SetUserVerifyEmailToUsedTx(id string, tx bun.Tx) error
}

type SubThreadRepository interface {
	Save(subThread *model.SubThread) error
	SaveTx(subThread *model.SubThread, tx bun.Tx) error
	GetByID(subThreadID string) (*model.SubThread, error)
	GetByName(name string) (*model.SubThread, error)
	IncrementFollowersCountTx(subThreadID string, tx bun.Tx) error
	SaveSubThreadFollowerTx(subthreadFollower *model.SubThreadFollower, tx bun.Tx) error
}
