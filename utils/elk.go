package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)

func returnVersion(es *elasticsearch.Client) {
	var (
		r map[string]interface{}
	)

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
}

func SendMessage(logSnippet []*Alb) {

	cfg := elasticsearch.Config{
		Addresses: []string{""},
		Username:  "",
		Password:  "",
	}

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// returnVersion(es)
	num := 1
	for _, v := range logSnippet {
		jsonString, _ := json.Marshal(v)
		body := strings.NewReader(string(jsonString))

		insertIndex(es, body, num)
		num++
	}

}

func insertIndex(es *elasticsearch.Client, body *strings.Reader, num int) {

	req := esapi.IndexRequest{
		Index: "",
		// DocumentID: strconv.Itoa(num),
		Body:    body,
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	fmt.Println(res.String())
}
