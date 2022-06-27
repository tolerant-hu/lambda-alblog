package utils

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)

type Elasticer struct {
	C  *Config
	Es *elasticsearch.Client
}

func (e *Elasticer) SyncIndex(logSnippet []*Alb) error {

	for _, v := range logSnippet {
		jsonString, _ := json.Marshal(v)
		body := strings.NewReader(string(jsonString))

		req := esapi.IndexRequest{
			Index:   e.C.ElkIndex,
			Body:    body,
			Refresh: "true",
		}

		res, err := req.Do(context.Background(), e.Es)
		if err != nil {
			return err
		}

		defer res.Body.Close()
	}

	return nil

}
