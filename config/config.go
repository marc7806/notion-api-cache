package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DbUri           string
	DbName          string
	NotionApiKey    string
	NotionDatabases []string

	CacheSchedulerDays    int
	CacheSchedulerHours   int
	CacheSchedulerMinutes int
	CacheOnStartup        bool
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
	CacheSchedulerHours, _ = strconv.Atoi(Get("CACHE_SCHEDULER_HOURS"))
	CacheSchedulerMinutes, _ = strconv.Atoi(Get("CACHE_SCHEDULER_MINUTES"))
	CacheSchedulerDays, _ = strconv.Atoi(Get("CACHE_SCHEDULER_DAYS"))
	CacheOnStartup, _ = strconv.ParseBool(Get("CACHE_ON_STARTUP"))
}
