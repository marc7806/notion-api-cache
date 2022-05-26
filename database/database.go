package database

import (
	"context"
	"fmt"
	"log"

	"github.com/marc7806/notion-cache/config"
	"github.com/marc7806/notion-cache/notion"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectDb() *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.DbUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully initialized database connection")
	mongoClient = client
	return client
}

func DisconnectDb() {
	err := mongoClient.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func QueryData(collectionId string, findQuery *primitive.M, sort interface{}, pageSize int64, start int64) []*notion.Page {
	client := ConnectDb()
	var options options.FindOptions
	options.Limit = &pageSize
	options.Skip = &start
	// options.Sort = sort

	var results []*notion.Page
	cur, err := client.Database(config.DbName).Collection(collectionId).Find(context.Background(), findQuery, &options)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var elem notion.Page
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	return results
}
