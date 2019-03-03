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

// NextSequence returns next ID value
func (db *Database) NextSequence(ctx context.Context, idName string) (int64, error) {
	var c counter
	err := db.Collection(sequenceCollectionName).FindOneAndUpdate(ctx,
		bson.M{"id": idName},
		bson.M{"$inc": bson.M{"seq": 1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)).
		Decode(&c)
	return c.Seq, err
}

// InitSequence sets last used value for ID name; InitSequence needs for data migration only
func (db *Database) InitSequence(ctx context.Context, idName string, idValue int64) error {
	_, err := db.Collection(sequenceCollectionName).InsertOne(ctx,
		&counter{
			ID:  idName,
			Seq: idValue,
		})
	return err
}
