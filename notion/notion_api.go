package notion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/marc7806/notion-cache/config"
)

type NotionResponse struct {
	Object     string                 `json:"object"`
	Results    []NotionDatabaseObject `json:"results"`
	NextCursor string                 `json:"next_cursor"`
	HasMore    bool                   `json:"has_more"`
}

type NotionDatabaseObject struct {
	Object         string    `json:"object"`
	ID             string    `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedTime    time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	CreatedBy      struct {
		Object string `json:"object"`
		ID     string `json:"id"`
	} `json:"created_by"`
	LastEditedBy struct {
		Object string `json:"object"`
		ID     string `json:"id"`
	} `json:"last_edited_by"`
	Cover  interface{} `json:"cover"`
	Icon   interface{} `json:"icon"`
	Parent struct {
		Type       string `json:"type"`
		DatabaseID string `json:"database_id"`
	} `json:"parent"`
	Archived   bool                              `json:"archived"`
	Properties map[string]map[string]interface{} `json:"properties"`
	URL        string                            `json:"url"`
}

func FetchNotionDataByDatabaseId(database_id string) *NotionResponse {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", database_id)
	payload := strings.NewReader("{\"page_size\":10}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", "2022-02-22")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.NotionApiKey))

	log.Printf("Fetch notion data from %s", url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if !strings.HasPrefix(res.Status, "2") {
		log.Fatalf("Status %s: %s", res.Status, body)
	}
	var parsedResponse *NotionResponse
	json.Unmarshal(body, &parsedResponse)
	return parsedResponse
}
