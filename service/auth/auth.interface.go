package auth

import (
	"realworld-authentication/model"
)

type AuthStorage interface {
	CreateUser(data *model.User) (*model.User, error)
	UpdateUser(query, data *model.User) (*model.User, error)
	GetUserByUsernameOrEmail(username, email string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
	UpdateUserPassword(query *model.User, password string) (*model.User, error)
	DeleteToken(token string) error
}
