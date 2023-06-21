package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserProfileUpdateDto struct {
	User struct {
		Email    string              `json:"email,omitempty"`
		Username string              `json:"username,omitempty"`
		Bio      *string             `json:"bio,omitempty"`
		Avatar   *primitive.ObjectID `json:"avatar,omitempty"`
	} `json:"user" validate:"required"`
}
