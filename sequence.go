package nosql

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const sequenceCollectionName = "counters"

type counter struct {
	ID  string `bson:"id"`
	Seq int64  `bson:"seq"`
}

// NextSequence returns next ID value.
// Make sure that the "counters" collection has an index by "id" field
func (db *Database) NextSequence(ctx context.Context, idName string) (int64, error) {
	var c counter
	err := db.Collection(sequenceCollectionName).FindOneAndUpdate(ctx,
		bson.M{"id": idName},
		bson.M{"$inc": bson.M{"seq": 1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)).
		Decode(&c)
	return c.Seq, err
}

// NextSequences returns first ID value for specified range size
func (db *Database) NextSequences(ctx context.Context, idName string, size int64) (int64, error) {
	var c counter
	err := db.Collection(sequenceCollectionName).FindOneAndUpdate(ctx,
		bson.M{"id": idName},
		bson.M{"$inc": bson.M{"seq": size}},
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)).
		Decode(&c)
	return c.Seq - (size - 1), err
}

// InitSequence sets last used value for ID name; InitSequence needs for data migration only
func (db *Database) InitSequence(ctx context.Context, idName string, idValue int64) error {
	return db.Collection(sequenceCollectionName).FindOneAndUpdate(ctx,
		bson.M{"id": idName},
		bson.M{"$set": bson.M{"seq": idValue}},
		options.FindOneAndUpdate().SetUpsert(true)).
		Err()
}
