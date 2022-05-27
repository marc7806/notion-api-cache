package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DbUri           string
	DbName          string
	NotionApiKey    string
	NotionDatabases []string
)

func Get(key string) string {
	return os.Getenv(key)
}

func init() {
	log.Println("Initializing environment")
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	DbUri = Get("MONGODB_URI")
	DbName = Get("MONGODB_NAME")
	NotionApiKey = Get("NOTION_API_KEY")
	NotionDatabases = strings.Split(Get("NOTION_DATABASES"), ",")
}
