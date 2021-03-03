package repository

import (
	"github.com/oklays/golang-restapi/src/modules/user/model"
)

type UserRepository interface {
	Save(*model.User) error
	Update(int64, *model.User) error
	Delete(int64) error
	FindByID(int64) (*model.User, error)
	FindAll() (model.Users, error)
}
