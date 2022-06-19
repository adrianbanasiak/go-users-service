package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Init(URL string) (*mongo.Client, error) {
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URL))
	if err != nil {
		return nil, err
	}

	client = c
	return client, nil
}
