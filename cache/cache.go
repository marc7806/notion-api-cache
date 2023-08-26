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

type CacheStatus string

const (
	Idle       CacheStatus = "idle"
	Refreshing             = "refreshing"
	Clearing               = "clearing"
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
	LastUpdated    time.Time
	NumUpdatedDocs int
	Status         CacheStatus
	refreshState   *RefreshState
)

func InitCache() {
	log.Println("Initializing Notion cache")
	// initialize refresh state
	refreshState = new(RefreshState)

	if config.CacheOnStartup {
		HandleCacheRefresh()
	}
}

func HandleCacheRefresh() bool {
	if !refreshState.isRefreshing {
		LastUpdated = time.Now()
		Status = Refreshing
		go refreshNotionCache()
	}
	return refreshState.isRefreshing
}

func HandleCacheClear() bool {
	if !refreshState.isRefreshing {
		LastUpdated = time.Now()
		Status = Clearing
		go clearCache()
	}
	return refreshState.isRefreshing
}

func refreshNotionCache() {
	refreshState.setRefreshState(true)
	defer refreshState.setRefreshState(false)

	updatedDocsLength, err := cacheNotionDatabases(database.DataStore.Db, config.NotionDatabases)
	if err != nil {
		log.Fatal(err)
	}

	NumUpdatedDocs = updatedDocsLength
	Status = Idle
}

func cacheNotionDatabases(db *mongo.Database, databases []string) (updatedDocsLength int, err error) {
	var numDocuments int

	// error handler function
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	for _, notionDatabaseId := range databases {
		log.Println("Syncing notion data to database")
		collection := db.Collection(notionDatabaseId)
		syncTime := time.Now()

		startCursor := ""
		hasMore := true
		for hasMore {
			notionData := notion.FetchNotionDataByDatabaseId(notionDatabaseId, startCursor)

			for _, page := range notionData.Results {
				update := bson.D{{"$set", notion.ParsePage(&page, &syncTime)}}
				// if query matches then update document, otherwise insert as new
				opts := options.Update().SetUpsert(true)
				result, err := collection.UpdateOne(context.Background(), bson.M{"_id": page.ID}, update, opts)
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
			startCursor = notionData.NextCursor
			hasMore = notionData.HasMore
			log.Printf("Finished page. Go to next cursor: %s", startCursor)
		}

		// delete all entities with older sync timestamp - Those represent entities that got deleted in notion
		result, err := collection.DeleteMany(context.Background(), bson.M{"lastsynctime": bson.M{"$lt": syncTime}})
		if err != nil {
			panic(err)
		}
		log.Printf("Deleted %d documents without sync", result.DeletedCount)

	}
	return numDocuments, nil
}

func clearCache() {
	refreshState.setRefreshState(true)
	defer refreshState.setRefreshState(false)

	log.Println("Start clearing database")
	var updateDocumentsCount int
	for _, notionDatabaseId := range config.NotionDatabases {
		collection := database.DataStore.Db.Collection(notionDatabaseId)

		// Delete all the documents in the collection
		deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted %v documents in collection %s\n", deleteResult.DeletedCount, notionDatabaseId)
		updateDocumentsCount += int(deleteResult.DeletedCount)
	}

	NumUpdatedDocs = updateDocumentsCount
	Status = Idle
}
