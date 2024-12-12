package nosql

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Aggregater is a part of the mongo.Collection interface.
type Aggregater interface {
	Aggregate(
		ctx context.Context, pipeline interface{}, opts ...options.Lister[options.AggregateOptions],
	) (*mongo.Cursor, error)
}

// AggregateMany executes an aggregate command against the collection.
// Any reference to mongo.Collection conforms to Aggregater interface.
func AggregateMany(
	ctx context.Context, aggregater Aggregater, pipeline interface{}, opts ...options.Lister[options.AggregateOptions],
) *ManyResult {
	cursor, err := aggregater.Aggregate(ctx, pipeline, opts...)

	return &ManyResult{
		ctx:    ctx,
		cursor: cursor,
		err:    err,
	}
}
