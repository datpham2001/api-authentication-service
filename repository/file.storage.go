package repository

import (
	"context"
	"mime/multipart"
	"realworld-authentication/config/env"
	"realworld-authentication/model"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type fileStorage struct {
	Instace *Instance
	S3      *s3.Client
}

func NewFileStorage(db *mongo.Database) *fileStorage {
	ins := &Instance{
		ColName:        "file",
		TemplateObject: &model.File{},
	}
	ins.ApplyDatabase(db)

	creds := credentials.NewStaticCredentialsProvider(env.AppConfig.AWSAccessKeyID, env.AppConfig.AWSSecretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds))
	if err != nil {
		return nil
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = env.AppConfig.AWSRegion
	})

	r := &fileStorage{
		Instace: ins,
		S3:      client,
	}

	return r
}

func (f *fileStorage) UploadFile(fileName string, file multipart.File, fileType string) (*model.File, error) {
	objInput := &s3.PutObjectInput{
		Bucket:      aws.String(env.AppConfig.AWSBucketName),
		Key:         aws.String(fileName),
		Body:        file,
		ACL:         "public-read",
		ContentType: &fileType,
	}

	uploader := manager.NewUploader(f.S3)
	uploadResult, err := uploader.Upload(context.TODO(), objInput)
	if err != nil {
		return nil, err
	}

	// insert data to db
	resp, err := f.createFile(&model.File{
		Key: fileName,
		Url: uploadResult.Location,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (f *fileStorage) DeleteFile(fileName string) error {
	delObj := &s3.DeleteObjectInput{
		Bucket: aws.String(env.AppConfig.AWSBucketName),
		Key:    aws.String(fileName),
	}
	_, err := f.S3.DeleteObject(context.TODO(), delObj)
	if err != nil {
		return err
	}

	// delete data in db
	return f.deleteFile(fileName)
}

func (f *fileStorage) createFile(data *model.File) (*model.File, error) {
	resp, err := f.Instace.Create(data)
	if err != nil {
		return nil, err
	}

	return resp.([]*model.File)[0], nil
}

func (f *fileStorage) UpdateFileByID(id *primitive.ObjectID, data *model.File) (*model.File, error) {
	resp, err := f.Instace.UpdateOne(model.File{
		ID: id,
	}, data)
	if err != nil {
		return nil, err
	}

	return resp.([]*model.File)[0], nil
}

func (f *fileStorage) deleteFile(fileName string) error {
	return f.Instace.DeleteOne(model.File{
		Key: fileName,
	})
}
