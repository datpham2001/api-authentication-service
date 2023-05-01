package repository

import (
	"context"
	"errors"
	"realworld-authentication/helper"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Instance struct {
	db   *mongo.Database
	coll *mongo.Collection

	DBName         string
	ColName        string
	TemplateObject interface{}
}

func (m *Instance) ApplyDatabase(db *mongo.Database) *Instance {
	m.db = db
	m.coll = db.Collection(m.ColName)
	m.DBName = db.Name()

	return m
}

func (m *Instance) newObject() interface{} {
	t := reflect.TypeOf(m.TemplateObject)

	v := reflect.New(t)
	return v.Interface()
}

func (m *Instance) newList(limit int) interface{} {
	t := reflect.TypeOf(m.TemplateObject)
	return reflect.MakeSlice(reflect.SliceOf(t), 0, limit).Interface()
}

func (m *Instance) convertToBson(ent interface{}) (bson.M, error) {
	if ent == nil {
		return bson.M{}, nil
	}

	sel, err := bson.Marshal(ent)
	if err != nil {
		return nil, err
	}

	obj := bson.M{}
	_ = bson.Unmarshal(sel, &obj)

	return obj, nil
}

func (m *Instance) convertToObject(b bson.M) (interface{}, error) {
	obj := m.newObject()

	if b == nil {
		return obj, nil
	}

	bytes, err := bson.Marshal(b)
	if err != nil {
		return nil, err
	}

	_ = bson.Unmarshal(bytes, obj)
	return obj, nil
}

func (m *Instance) interfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

func (m *Instance) parseSingleResult(result *mongo.SingleResult, action string) *helper.APIResponse {
	// parse result
	obj := m.newObject()
	err := result.Decode(obj)
	if err != nil {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Error,
			Message:   "DB Error: " + err.Error(),
			ErrorCode: "MAP_OBJECT_FAILED",
		}
	}

	// put to slice
	list := m.newList(1)
	listValue := reflect.Append(reflect.ValueOf(list),
		reflect.Indirect(reflect.ValueOf(obj)))

	return &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: action + " " + m.ColName + " successfully.",
		Data:    listValue.Interface(),
	}
}

func (m *Instance) Create(ent interface{}) *helper.APIResponse {
	// check col
	if m.coll == nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: Collection " + m.ColName + " is not init.",
		}
	}

	// convert to bson
	obj, err := m.convertToBson(ent)
	if err != nil {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Error,
			Message:   "DB Error: " + err.Error(),
			ErrorCode: "MAP_OBJECT_FAILED",
		}
	}

	// init time
	if obj["created_time"] == nil {
		obj["created_time"] = time.Now()
	}

	// insert
	result, err := m.coll.InsertOne(context.TODO(), obj)
	if err != nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB Error: " + err.Error(),
		}
	}

	obj["_id"] = result.InsertedID
	ent, _ = m.convertToObject(obj)

	list := m.newList(1)
	listValue := reflect.Append(reflect.ValueOf(list),
		reflect.Indirect(reflect.ValueOf(ent)))

	return &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Create " + m.ColName + " successfully.",
		Data:    listValue.Interface(),
	}
}

// CreateMany insert many object into db
func (m *Instance) CreateMany(entityList interface{}) *helper.APIResponse {

	// check col
	if m.coll == nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: Create many - Collection " + m.ColName + " is not init.",
		}
	}

	list, err := m.interfaceSlice(entityList)
	if err != nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: Create many - Invalid slice.",
		}
	}

	var bsonList []interface{}
	now := time.Now()
	for _, item := range list {
		b, err := m.convertToBson(item)
		if err != nil {
			return &helper.APIResponse{
				Status:  helper.APIStatus.Error,
				Message: "DB error: Create many - Invalid bson object.",
			}
		}
		if b["created_time"] == nil {
			b["created_time"] = now
		}
		bsonList = append(bsonList, b)
	}

	result, err := m.coll.InsertMany(context.TODO(), bsonList)
	if err != nil {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Error,
			Message:   "DB Error: " + err.Error(),
			ErrorCode: "CREATE_FAILED",
		}
	}

	return &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Create " + m.ColName + "(s) successfully.",
		Data:    result.InsertedIDs,
	}
}

// UpdateOne Update one matched object.
func (m *Instance) UpdateOne(query interface{}, updater interface{}, opts ...*options.FindOneAndUpdateOptions) *helper.APIResponse {
	// check col
	if m.coll == nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: Collection " + m.ColName + " is not init.",
		}
	}

	// convert
	bUpdater, err := m.convertToBson(updater)
	if err != nil {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Error,
			Message:   "DB Error: " + err.Error(),
			ErrorCode: "MAP_OBJECT_FAILED",
		}
	}
	bUpdater["last_updated_time"] = time.Now()

	// transform to bson
	converted, err := m.convertToBson(query)
	if err != nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: UpdateOne - Cannot convert object - " + err.Error(),
		}
	}

	// do update
	if opts == nil {
		after := options.After
		opts = []*options.FindOneAndUpdateOptions{
			{
				ReturnDocument: &after,
			},
		}
	}
	result := m.coll.FindOneAndUpdate(context.TODO(), converted, bson.M{"$set": bUpdater}, opts...)
	if result.Err() != nil {
		detail := ""
		if result != nil {
			detail = result.Err().Error()
		}
		return &helper.APIResponse{
			Status:    helper.APIStatus.Notfound,
			Message:   "Not found any matched " + m.ColName + ". Error detail: " + detail,
			ErrorCode: "NOT_FOUND",
		}
	}

	return m.parseSingleResult(result, "UpdateOne")
}

// Query Get all object in DB
func (m *Instance) Query(query interface{}, offset int64, limit int64, sortFields *bson.M) *helper.APIResponse {
	// check col
	if m.coll == nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: Collection " + m.ColName + " is not init.",
		}
	}
	opt := &options.FindOptions{}
	k := int64(1000)
	if limit <= 0 {
		opt.Limit = &k
	} else {
		opt.Limit = &limit
	}
	if offset > 0 {
		opt.Skip = &offset
	}
	if sortFields != nil {
		opt.Sort = sortFields
	}

	// transform to bson
	converted, err := m.convertToBson(query)
	if err != nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: QueryOne - Cannot convert object - " + err.Error(),
		}
	}

	result, err := m.coll.Find(context.TODO(), converted, opt)

	if err != nil || result.Err() != nil {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Notfound,
			Message:   "Not found any matched " + m.ColName + ".",
			ErrorCode: "NOT_FOUND",
		}
	}

	list := m.newList(int(limit))
	err = result.All(context.TODO(), &list)
	result.Close(context.TODO())
	if err != nil || reflect.ValueOf(list).Len() == 0 {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Notfound,
			Message:   "Not found any matched " + m.ColName + ".",
			ErrorCode: "NOT_FOUND",
		}
	}

	return &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Query " + m.ColName + " successfully.",
		Data:    list,
	}
}

// Query Get all object in DB
func (m *Instance) QueryAll() *helper.APIResponse {
	// check col
	if m.coll == nil {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Error,
			Message:   "DB error: Collection " + m.ColName + " is not init.",
			ErrorCode: "NOT_INIT_YET",
		}
	}
	rs, err := m.coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Notfound,
			Message:   "Not found any " + m.ColName + ".",
			ErrorCode: "NOT_FOUND",
		}
	}

	list := m.newList(1000)
	rs.All(context.TODO(), &list)
	rs.Close(context.TODO())
	if reflect.ValueOf(list).Len() == 0 {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Notfound,
			Message:   "Not found any matched " + m.ColName + ".",
			ErrorCode: "NOT_FOUND",
		}
	}
	return &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Query " + m.ColName + " successfully.",
		Data:    list,
	}
}

// QueryOne ...
func (m *Instance) QueryOne(query interface{}) *helper.APIResponse {
	// check col
	if m.coll == nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: Collection " + m.ColName + " is not init.",
		}
	}

	// transform to bson
	converted, err := m.convertToBson(query)
	if err != nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: QueryOne - Cannot convert object - " + err.Error(),
		}
	}

	// do find
	result := m.coll.FindOne(context.TODO(), converted)

	if result == nil || result.Err() != nil {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Notfound,
			Message:   "Not found any matched " + m.ColName + ".",
			ErrorCode: "NOT_FOUND",
		}
	}

	return m.parseSingleResult(result, "Query")
}

func (m *Instance) CreateIndex(keys bson.D, options *options.IndexOptions) error {
	_, err := m.coll.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    keys,
		Options: options,
	})

	return err
}

// Count Count object which matched with query.
func (m *Instance) Count(query interface{}) *helper.APIResponse {
	// check col
	if m.coll == nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: Collection " + m.ColName + " is not init.",
		}
	}

	// convert query
	converted, err := m.convertToBson(query)
	if err != nil {
		return &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: "DB error: Count - Cannot convert object - " + err.Error(),
		}
	}

	count, err := m.coll.CountDocuments(context.TODO(), converted)
	if err != nil {
		return &helper.APIResponse{
			Status:    helper.APIStatus.Error,
			Message:   "Count error: " + err.Error(),
			ErrorCode: "COUNT_FAILED",
		}
	}

	return &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Count query executed successfully.",
		Data:    count,
	}

}
