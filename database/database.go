package database

import (
	"context"
	"fmt"
	"log"

	"github.com/marc7806/notion-cache/config"
	"github.com/marc7806/notion-cache/notion"
	"go.mongodb.org/mongo-driver/bson"
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

func QueryData(collectionId string, findQuery *primitive.M, sort interface{}, pageSize int64, startCursor string) (result []*notion.Page, nextCursor string, hasMore bool) {
	client := ConnectDb()
	var options options.FindOptions
	// temporary add one to page size for computing hasMore property
	pageSize = pageSize + 1
	options.Limit = &pageSize
	options.Sort = bson.M{"_id": 1}

	if startCursor != "" {
		findQuery = addStartCursor(findQuery, startCursor)
	}

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

	hasMore = len(results) > 0 && len(results) == int(pageSize)
	if hasMore {
		// remove added element for hasMore computation
		results = results[:len(results)-1]
		nextCursor = results[len(results)-1].ID
	} else {
		nextCursor = ""
		hasMore = false
	}

	if len(results) == 0 {
		results = []*notion.Page{}
	}

	return results, nextCursor, hasMore
}

func addStartCursor(findQuery *primitive.M, startCursor string) *bson.M {
	return &bson.M{
		"$and": bson.A{
			bson.M{
				"_id": bson.M{"$gt": startCursor},
			},
			findQuery,
		},
	}
}
