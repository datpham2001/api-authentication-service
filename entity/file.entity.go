package entity

import (
	"realworld-authentication/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadFileResponse struct {
	File struct {
		FileID *primitive.ObjectID `json:"_id"`
		Path   string              `json:"path"`
	} `json:"file"`
}

func NewUploadFileResponse(u *model.File) *UploadFileResponse {
	resp := new(UploadFileResponse)
	resp.File.FileID = u.ID
	resp.File.Path = u.Url

	return resp
}
