package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/marc7806/notion-cache/cache"
	"github.com/marc7806/notion-cache/config"
	"github.com/marc7806/notion-cache/database"
)

func main() {
	config.LoadEnvironment()
	client := database.InitClient()

	// Check the connection
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	cache.CacheNotionDatabases(client, config.NotionDatabases)

	// Close the connection once no longer needed
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}

	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World")
}
