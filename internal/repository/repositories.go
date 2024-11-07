package repository

import (
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/uptrace/bun"
)

type UserRepository interface {
	Save(user *model.User) error
	SaveTx(user *model.User, tx bun.Tx) error
}
