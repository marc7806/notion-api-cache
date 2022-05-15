package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/marc7806/notion-cache/config"
	"github.com/marc7806/notion-cache/notion"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CacheNotionDatabases(client *mongo.Client, databases []string) {
	for _, notionDatabaseId := range databases {
		log.Println("Saving notion data to database")

		notionData := notion.FetchNotionDataByDatabaseId(notionDatabaseId)
		collection := client.Database(config.DbName).Collection(notionDatabaseId)

		for _, page := range notionData.Results {
			update := bson.D{{"$set", page}}
			opts := options.Update().SetUpsert(true)
			result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": page.ID}, update, opts)
			if err != nil {
				panic(err)
			}

			if result.UpsertedCount > 0 {
				log.Printf("Inserted new document %s", page.ID)
			} else {
				log.Printf("Updated existing document %s", page.ID)
			}
		}
	}
}

func ClearCache(client *mongo.Client, databases []string) {
	log.Println("Start clearing database")
	for _, notionDatabaseId := range databases {
		collection := client.Database(config.DbName).Collection(notionDatabaseId)

		// Delete all the documents in the collection
		deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted %v documents in collection %s\n", deleteResult.DeletedCount, notionDatabaseId)
	}
}
