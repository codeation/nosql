// Package nosql implements wrappers for go.mongodb.org/mongo-driver
package nosql

import (
	"context"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errVariableType = errors.New("wrong variable data type")

// Database is a mongo.Database wrapper.
type Database struct {
	*mongo.Database
}

// Collection is a mongo.Collection wrapper.
type Collection struct {
	*mongo.Collection
}

// ManyResult contains Find results.
type ManyResult struct {
	ctx    context.Context
	cursor *mongo.Cursor
	err    error
}

// NewDatabase creates a new Database instance.
func NewDatabase(db *mongo.Database) *Database { return &Database{Database: db} }

// Collection returns a handle for collection of database.
func (db *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{
		Collection: db.Database.Collection(name, opts...),
	}
}

// FindMany finds all documents that match the filter and options.
func (c *Collection) FindMany(
	ctx context.Context, filter interface{}, opts ...*options.FindOptions,
) *ManyResult {
	cursor, err := c.Collection.Find(ctx, filter, opts...)

	return &ManyResult{
		ctx:    ctx,
		cursor: cursor,
		err:    err,
	}
}

// Cursor returns a underlying *mongo.Cursor.
func (a *ManyResult) Cursor() *mongo.Cursor { return a.cursor }

// Err returns a underlying error.
func (a *ManyResult) Err() error { return a.err }

// Decode decodes all found documents into a variable.
// The data parameter may be a pointer to an slice of struct.
// Also data parameter may be a pointer to an slice of pointers to a struct.
// For examples:
//
//    var data1 []Struct // slice of struct
//    err := collection.FindMany(ctx, bson.D{}).Decode(&data1) // pointer to an slice of ...
//
//    var data2 []*Struct // slice of pointers to a struct
//    err := collection.FindMany(ctx, bson.D{}).Decode(&data2) // pointer to an slice of ...
//
// If no documents are found, an empty slice is returned.
func (a *ManyResult) Decode(data interface{}) error {
	if a.err != nil {
		return a.err
	}

	// detect data and element types
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return errVariableType
	}

	sliceType := reflect.TypeOf(data).Elem()
	if sliceType.Kind() != reflect.Slice {
		return errVariableType
	}

	elemType := sliceType.Elem()
	elemPtr := false

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
		elemPtr = true
	}

	// decode documents
	records := reflect.MakeSlice(sliceType, 0, 0)

	for a.cursor.Next(a.ctx) {
		record := reflect.New(elemType)

		if err := a.cursor.Decode(record.Interface()); err != nil {
			return err
		}

		if elemPtr {
			records = reflect.Append(records, record)
		} else {
			records = reflect.Append(records, record.Elem())
		}
	}

	if err := a.cursor.Err(); err != nil {
		return err
	}

	reflect.ValueOf(data).Elem().Set(records)

	return nil
}
