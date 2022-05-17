package database

import (
	"context"
	"fmt"

	"github.com/marc7806/notion-cache/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.DbUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully initialized database connection")
	MongoClient = client
	return client
}
