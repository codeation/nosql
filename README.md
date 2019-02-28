# nosql
Layer for convenient use of go.mongodb.org/mongo-driver

[![GoDoc](https://godoc.org/github.com/codeation/nosql?status.svg)](https://godoc.org/github.com/codeation/nosql)

# FindAll(...).Decode(...) chain

You can use the FindAll Decode chain to decode an array of documents from mongodb.

```
	var data []Elem
	if err := collection.FindAll(ctx, bson.D{}).Decode(&data); err != nil {
		return err
	}

	// Some using of documents slice
	for _, e := range data {
		fmt.Println(e.ID)
	}

```

It is like calling FindOne Decode chain to decode a single document in a
[standard mongodb driver](https://godoc.org/go.mongodb.org/mongo-driver/mongo).

FindAll wraps the
[func (*Collection) Find](https://godoc.org/go.mongodb.org/mongo-driver/mongo#Collection.Find)
results, so the parameters are the same.

Data parameter of func Decode may be a pointer to an slice of struct.
Also data parameter may be a pointer to an slice of pointers to a struct, see below.

If no documents are found, an empty slice is returned.

# Minimal example

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

	db := nosql.NewDatabase(client.Database("test")) // wrap mongo.Database handle
	collection := db.Collection("test")

	var data []*Elem
	if err := collection.FindAll(ctx, bson.D{}).Decode(&data); err != nil {
		return
	}

	for _, e := range data {
		fmt.Println(e.Num, e.Str)
	}
}
```

See the [documentation](https://godoc.org/github.com/codeation/nosql) for details.
