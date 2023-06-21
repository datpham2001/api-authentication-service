package file

import (
	"mime/multipart"
	"realworld-authentication/entity"
)

type fileService struct {
	storage FileStorage
}

func NewFileService(s FileStorage) *fileService {
	return &fileService{
		storage: s,
	}
}

func (u *fileService) UploadFile(fileName string, file multipart.File, fileType string) (*entity.UploadFileResponse, error) {
	resp, err := u.storage.UploadFile(fileName, file, fileType)
	if err != nil {
		return nil, err
	}

	return entity.NewUploadFileResponse(resp), nil
}

func (u *fileService) DeleteFile(fileName string) error {
	return u.storage.DeleteFile(fileName)
}
