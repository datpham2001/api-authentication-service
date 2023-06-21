package file

import (
	"mime/multipart"
	"realworld-authentication/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStorage interface {
	UploadFile(fileName string, file multipart.File, fileType string) (*model.File, error)
	DeleteFile(fileName string) error
	UpdateFileByID(id *primitive.ObjectID, data *model.File) (*model.File, error)
}
