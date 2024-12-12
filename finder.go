package nosql

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Finder is a part of the mongo.Collection interface.
type Finder interface {
	Find(ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOptions],
	) (*mongo.Cursor, error)
}

// FindMany finds all documents from the collection that match the filter and options.
// Any reference to mongo.Collection conforms to Finder interface.
func FindMany(
	ctx context.Context, finder Finder, filter interface{}, opts ...options.Lister[options.FindOptions],
) *ManyResult {
	cursor, err := finder.Find(ctx, filter, opts...)

	return &ManyResult{
		ctx:    ctx,
		cursor: cursor,
		err:    err,
	}
}
