package repository

import (
	"realworld-authentication/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type authStorage struct {
	Instance *Instance
}

func NewAuthStorage(db *mongo.Database) *authStorage {
	ins := &Instance{
		ColName:        "auth",
		TemplateObject: &model.User{},
	}
	ins.ApplyDatabase(db)

	r := &authStorage{
		Instance: ins,
	}

	return r
}

func (r *authStorage) CreateUser(data *model.User) (*model.User, error) {
	dataRes, err := r.Instance.Create(data)
	if err != nil {
		return nil, err
	}

	return dataRes.([]*model.User)[0], nil
}

func (r *authStorage) UpdateUser(query, data *model.User) (*model.User, error) {
	dataRes, err := r.Instance.UpdateOne(query, data)
	if err != nil {
		return nil, err
	}

	return dataRes.([]*model.User)[0], nil
}

func (r *authStorage) GetUserByUsernameOrEmail(username, email string) (*model.User, error) {
	dataRes, err := r.Instance.QueryOne(model.User{
		ComplexQuery: []*bson.M{
			{
				"$or": []*bson.M{{
					"username": username,
				}, {
					"email": email,
				}},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return dataRes.([]*model.User)[0], nil
}

func (r *authStorage) GetUserByEmail(email string) (*model.User, error) {
	dataRes, err := r.Instance.QueryOne(model.User{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	return dataRes.([]*model.User)[0], nil
}

func (r *authStorage) GetUserByID(id string) (*model.User, error) {
	dataRes, err := r.Instance.QueryOne(model.User{
		UserID: id,
	})
	if err != nil {
		return nil, err
	}

	return dataRes.([]*model.User)[0], nil
}

func (r *authStorage) UpdateUserPassword(query *model.User, password string) (*model.User, error) {
	dataRes, err := r.Instance.UpdateOne(query, &model.User{
		HashedPassword: password,
	})
	if err != nil {
		return nil, err
	}

	return dataRes.([]*model.User)[0], nil
}

func (r *authStorage) DeleteToken(token string) error {
	return r.Instance.DeleteOne(model.User{
		RefreshToken: token,
	})
}
