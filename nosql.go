// Package nosql implements (*Collection) FindAll method
package nosql

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is a mongo.Database wrapper
type Database struct {
	*mongo.Database
}

// Collection is a mongo.Collection wrapper
type Collection struct {
	*mongo.Collection
}

// AllResult contains Find results
type AllResult struct {
	cursor *mongo.Cursor
	err    error
}

// NewDatabase creates a new Database instance
func NewDatabase(db *mongo.Database) *Database { return &Database{Database: db} }

// Collection returns a handle for collection of database
func (db *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{
		Collection: db.Database.Collection(name, opts...),
	}
}

// FindAll finds all documents that match the filter and options
func (c *Collection) FindAll(
	ctx context.Context, filter interface{}, opts ...*options.FindOptions,
) *AllResult {
	cursor, err := c.Collection.Find(ctx, filter, opts...)
	return &AllResult{
		cursor: cursor,
		err:    err,
	}
}

// Cursor returns a underlying *mongo.Cursor
func (a *AllResult) Cursor() *mongo.Cursor { return a.cursor }

// Err returns a underlying error
func (a *AllResult) Err() error { return a.err }

// All decodes all found documents into a variable.
// The data parameter must be a pointer to an slice of pointers to a struct.
// For example:
//
//    var data []*Struct // slice of pointers to a struct
//    err := collection.FindAll(ctx, bson.D{}).All(&data) // pointer to an slice of ...
//
func (a *AllResult) All(data interface{}) error {
	if a.err != nil {
		return a.err
	}
	elemType := reflect.TypeOf(data).Elem().Elem().Elem()
	records := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(elemType)), 0, 200)
	for a.cursor.Next(context.Background()) {
		record := reflect.New(elemType)
		if err := a.cursor.Decode(record.Interface()); err != nil {
			return err
		}
		records = reflect.Append(records, record)
	}
	if err := a.cursor.Err(); err != nil {
		return err
	}
	reflect.ValueOf(data).Elem().Set(records)
	return nil
}
