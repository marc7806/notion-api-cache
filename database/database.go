package database

import (
	"context"
	"fmt"
	"log"

	"github.com/marc7806/notion-cache/config"
	sortparser "github.com/marc7806/notion-cache/database/sort-parser"
	"github.com/marc7806/notion-cache/notion"
	"github.com/marc7806/notion-cache/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoDbHandle *mongo.Database

func ConnectDb() *mongo.Database {

	if mongoClient != nil && mongoDbHandle != nil {
		return mongoDbHandle
	}

	clientOptions := options.Client().ApplyURI(config.DbUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Ping failed...")
	}

	fmt.Printf("Successfully initialized database connection")
	mongoClient = client
	mongoDbHandle = client.Database(config.DbName)
	return mongoDbHandle
}

func DisconnectDb() {
	err := mongoClient.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func QueryData(collectionId string, findQuery *primitive.M, sort []types.QuerySort, pageSize int64, startCursor string) (result []*notion.Page, nextCursor string, hasMore bool) {
	db := ConnectDb()
	var options options.FindOptions
	// temporary add one to page size for computing hasMore property
	pageSize = pageSize + 1
	options.Limit = &pageSize
	options.Sort = sortparser.ParseSortOptions(sort)

	if startCursor != "" {
		var startCursorPage notion.Page
		err := db.Collection(collectionId).FindOne(context.Background(), bson.D{{"_id", startCursor}}).Decode(&startCursorPage)
		if err != nil {
			log.Fatal("Start cursor not found", err)
		} else {
			findQuery = addStartCursor(findQuery, startCursorPage, sort)
		}
	}

	var results []*notion.Page
	cur, err := db.Collection(collectionId).Find(context.Background(), findQuery, &options)
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

func addStartCursor(findQuery *primitive.M, startCursor notion.Page, sort []types.QuerySort) *bson.M {
	cursorQuery := bson.M{}
	if sort != nil {
		for _, sortEntry := range sort {
			cursorQuery[notion.BuildPropertyValueAccessorString(sortEntry.Property)] = bson.M{"$gt": startCursor.Properties[sortEntry.Property].Value}
		}
	} else {
		cursorQuery["_id"] = bson.M{"$gt": startCursor}
	}

	return &bson.M{
		"$and": bson.A{
			cursorQuery,
			findQuery,
		},
	}
}
