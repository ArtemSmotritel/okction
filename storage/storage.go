package storage

import "github.com/artemsmotritel/oktion/types"

type Storage interface {
	GetUserByID(id int64) (*types.User, error)
	GetUsers() ([]types.User, error)
	SaveUser(user *types.User) error
	UpdateUser(id int64, request *types.UserUpdateRequest) error
	DeleteUser(id int64) error
	SeedData() error
}
