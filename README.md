# nosql
A wrapper to make it easier to use go.mongodb.org/mongo-driver

[![PkgGoDev](https://pkg.go.dev/badge/codeation/nosql)](https://pkg.go.dev/codeation/nosql)

## FindMany(...).Decode(...) chain

You can use the FindMany Decode chain to decode an array of documents from mongodb collection.

```
	collection := client.Database("test").Collection("test")

	var data []Elem
	if err := nosql.FindMany(ctx, collection, bson.D{}).Decode(&data); err != nil {
		return err
	}

	// Some using of documents slice
	for _, e := range data {
		fmt.Println(e.ID)
	}

```

It is like calling FindOne Decode chain to decode a single document in a
[standard mongodb driver](https://godoc.org/go.mongodb.org/mongo-driver/mongo).

FindMany wraps the
[func (*Collection) Find](https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Find)
results, so the parameters are the same.

Data parameter of func Decode may be a pointer to an slice of struct.
Also data parameter may be a pointer to an slice of pointers to a struct, see below.

If no documents are found, an empty slice is returned.

## Minimal example

```
package main

import (
	"context"
	"fmt"

	"github.com/codeation/nosql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Elem struct {
	Num int    `bson:"num"`
	Str string `bson:"str"`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		return
	}

	ctx := context.Background()
	if err = client.Connect(ctx); err != nil {
		return
	}
	defer client.Disconnect(ctx)

	collection := client.Database("test").Collection("test")

	var data []*Elem
	if err := nosql.FindMany(ctx, collection, bson.D{}).Decode(&data); err != nil {
		return
	}

	for _, e := range data {
		fmt.Println(e.Num, e.Str)
	}
}
```

## Wrapper example

```
package main

import (
	"context"
	"fmt"

	"github.com/codeation/nosql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Elem struct {
	Num int    `bson:"num"`
	Str string `bson:"str"`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		return
	}

	ctx := context.Background()
	if err = client.Connect(ctx); err != nil {
		return
	}
	defer client.Disconnect(ctx)

	db := nosql.NewDatabase(client.Database("test")) // wrap mongo.Database reference
	collection := db.Collection("test")

	var data []*Elem
	if err := collection.FindMany(ctx, bson.D{}).Decode(&data); err != nil {
		return
	}

	for _, e := range data {
		fmt.Println(e.Num, e.Str)
	}
}
```

## NextSequence func

NextSequence returns next ID value.

```
    id, err := db.NextSequence(ctx, "elemid")
    if err != nil {
        return err
    }
    e := &Element {
        ID: id,
        ... // Other fields
    }
```

This can be useful when you plan to use int64 values as IDs,
or you need to know the new ID before inserting the document.

NextSequence uses the atomic operation
[$inc](https://docs.mongodb.com/manual/reference/operator/update/inc/).

Make sure that the "counters" collection has an index by "id" field:

```
db.counters.createIndex( { id: 1 } )
```
