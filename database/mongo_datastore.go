package database

import (
	"context"
	"fmt"
	"log"

	"github.com/marc7806/notion-cache/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDataStore struct {
	Db     *mongo.Database
	Client *mongo.Client
}

var DataStore *MongoDataStore

func InitMongoDataStore() (*MongoDataStore, error) {
	log.Println("Initialize Mongo Datastore")

	if DataStore != nil {
		return DataStore, nil
	}

	clientOptions := options.Client().ApplyURI(config.DbUri)
	client, db, err := connectDb(clientOptions, config.DbName)

	if err != nil {
		return nil, err
	}

	DataStore = &MongoDataStore{Db: db, Client: client}
	return DataStore, nil
}

func DisconnectDb(datastore *MongoDataStore) {
	log.Println("Disconnect client...")
	err := datastore.Client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func connectDb(opts *options.ClientOptions, dbName string) (*mongo.Client, *mongo.Database, error) {
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot ping mongodb: %v", err)
	}

	log.Println("Successfully initialized database connection")
	var db = client.Database(dbName)
	return client, db, nil
}
