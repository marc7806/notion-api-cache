package cache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/marc7806/notion-cache/config"
	"github.com/marc7806/notion-cache/database"
	"github.com/marc7806/notion-cache/notion"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshState struct {
	mu           sync.Mutex
	isRefreshing bool
}

func (s *RefreshState) setRefreshState(state bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isRefreshing = state
}

var (
	lastUpdated    time.Time
	numUpdatedDocs int
	refreshState   *RefreshState
)

func init() {
	log.Println("Initializing Notion cache")
	// initialize refresh state
	refreshState = new(RefreshState)
	HandleCacheRefresh()
}

func HandleCacheRefresh() (bool, time.Time, int) {
	log.Printf("Refresh state: %v", refreshState.isRefreshing)
	if !refreshState.isRefreshing {
		go refreshNotionCache()
	}
	return refreshState.isRefreshing, lastUpdated, numUpdatedDocs
}

func refreshNotionCache() {
	refreshState.setRefreshState(true)
	client := database.InitClient()

	// Check the connection
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	updatedDocsLength, err := cacheNotionDatabases(client, config.NotionDatabases)
	if err != nil {
		log.Fatal(err)
	}

	// Close the connection once no longer needed
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	lastUpdated = time.Now()
	numUpdatedDocs = updatedDocsLength
	refreshState.setRefreshState(false)
}

func cacheNotionDatabases(client *mongo.Client, databases []string) (updatedDocsLength int, err error) {
	var numDocuments int
	for _, notionDatabaseId := range databases {
		log.Println("Saving notion data to database")

		notionData := notion.FetchNotionDataByDatabaseId(notionDatabaseId)
		collection := client.Database(config.DbName).Collection(notionDatabaseId)

		// error handler function
		defer func() {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}()

		for _, page := range notionData.Results {
			update := bson.D{{"$set", notion.ParsePage(&page)}}
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
		numDocuments += len(notionData.Results)
	}
	return numDocuments, nil
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
