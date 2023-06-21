package repository

import (
	"context"
	"errors"
	"fmt"
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

func (m *Instance) parseSingleResult(result *mongo.SingleResult, action string) (interface{}, error) {
	// parse result
	obj := m.newObject()
	err := result.Decode(obj)
	if err != nil {
		return nil, err
	}

	// put to slice
	list := m.newList(1)
	listValue := reflect.Append(reflect.ValueOf(list),
		reflect.Indirect(reflect.ValueOf(obj)))

	return listValue.Interface(), nil
}

func (m *Instance) Create(ent interface{}) (interface{}, error) {
	// check col
	if m.coll == nil {
		return nil, fmt.Errorf("%v is not inited", m.ColName)
	}

	// convert to bson
	obj, err := m.convertToBson(ent)
	if err != nil {
		return nil, err
	}

	// init time
	if obj["created_time"] == nil {
		obj["created_time"] = time.Now()
	}

	// insert
	result, err := m.coll.InsertOne(context.TODO(), obj)
	if err != nil {
		return nil, err
	}

	obj["_id"] = result.InsertedID
	ent, _ = m.convertToObject(obj)

	list := m.newList(1)
	listValue := reflect.Append(reflect.ValueOf(list),
		reflect.Indirect(reflect.ValueOf(ent)))

	return listValue.Interface(), nil
}

// CreateMany insert many object into db
func (m *Instance) CreateMany(entityList interface{}) (interface{}, error) {

	// check col
	if m.coll == nil {
		return nil, fmt.Errorf("%v is not inited", m.ColName)
	}

	list, err := m.interfaceSlice(entityList)
	if err != nil {
		return nil, err
	}

	var bsonList []interface{}
	now := time.Now()
	for _, item := range list {
		b, err := m.convertToBson(item)
		if err != nil {
			return nil, err
		}
		if b["created_time"] == nil {
			b["created_time"] = now
		}
		bsonList = append(bsonList, b)
	}

	result, err := m.coll.InsertMany(context.TODO(), bsonList)
	if err != nil {
		return nil, err
	}

	return result.InsertedIDs, nil
}

// UpdateOne Update one matched object.
func (m *Instance) UpdateOne(query interface{}, updater interface{}, opts ...*options.FindOneAndUpdateOptions) (interface{}, error) {
	// check col
	if m.coll == nil {
		return nil, fmt.Errorf("%v is not inited", m.ColName)
	}

	// convert
	bUpdater, err := m.convertToBson(updater)
	if err != nil {
		return nil, err
	}
	bUpdater["last_updated_time"] = time.Now()

	// transform to bson
	converted, err := m.convertToBson(query)
	if err != nil {
		return nil, err
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
		return nil, errors.New(detail)
	}

	return m.parseSingleResult(result, "UpdateOne")
}

// Query Get all object in DB
func (m *Instance) Query(query interface{}, offset int64, limit int64, sortFields *bson.M) (interface{}, error) {
	// check col
	if m.coll == nil {
		return nil, fmt.Errorf("%v is not inited", m.ColName)
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
		return nil, err
	}

	result, err := m.coll.Find(context.TODO(), converted, opt)

	if err != nil || result.Err() != nil {
		return nil, err
	}

	list := m.newList(int(limit))
	err = result.All(context.TODO(), &list)
	result.Close(context.TODO())
	if err != nil || reflect.ValueOf(list).Len() == 0 {
		return nil, err
	}

	return list, nil
}

// Query Get all object in DB
func (m *Instance) QueryAll() (interface{}, error) {
	// check col
	if m.coll == nil {
		return nil, fmt.Errorf("%v is not inited", m.ColName)
	}
	rs, err := m.coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	list := m.newList(1000)
	rs.All(context.TODO(), &list)
	rs.Close(context.TODO())
	if reflect.ValueOf(list).Len() == 0 {
		return nil, err
	}

	return list, nil
}

// QueryOne ...
func (m *Instance) QueryOne(query interface{}) (interface{}, error) {
	// check col
	if m.coll == nil {
		return nil, fmt.Errorf("%v is not inited", m.ColName)
	}

	// transform to bson
	converted, err := m.convertToBson(query)
	if err != nil {
		return nil, err
	}

	// do find
	result := m.coll.FindOne(context.TODO(), converted)

	if result == nil || result.Err() != nil {
		return nil, errors.New("document is not existed")
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
func (m *Instance) Count(query interface{}) (interface{}, error) {
	// check col
	if m.coll == nil {
		return nil, fmt.Errorf("%v is not inited", m.ColName)
	}

	// convert query
	converted, err := m.convertToBson(query)
	if err != nil {
		return nil, err
	}

	count, err := m.coll.CountDocuments(context.TODO(), converted)
	if err != nil {
		return nil, err
	}

	return count, nil
}

func (m *Instance) DeleteOne(query interface{}) error {
	// check col
	if m.coll == nil {
		return fmt.Errorf("%v is not inited", m.ColName)
	}

	// convert query
	converted, err := m.convertToBson(query)
	if err != nil {
		return err
	}

	_, err = m.coll.DeleteOne(context.TODO(), converted)
	if err != nil {
		return err
	}

	return nil
}
