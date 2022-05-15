package database

import (
	"context"
	"fmt"
	"os"

	"github.com/marc7806/notion-cache/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitClient() *mongo.Client {
	println("This is my config")
	println(config.DbUri)
	clientOptions := options.Client().ApplyURI(config.DbUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully initialized database connection")
	MongoClient = client
	return client
}

func getDatabase(client *mongo.Client) *mongo.Database {
	return client.Database(os.Getenv("MONGO_DB_NAME"))
}

func getCollection(database *mongo.Database, collection_name string) *mongo.Collection {
	return database.Collection(collection_name)
}
