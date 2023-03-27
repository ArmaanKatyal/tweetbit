package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/models"
	esv7 "github.com/elastic/go-elasticsearch/v7"
	esv7api "github.com/elastic/go-elasticsearch/v7/esapi"
)

// Tweet represents the repository used for interacting with tweet records
type Tweet struct {
	client *esv7.Client
	index  string
}

// NewTweet instantiates a new Tweet repository
func NewTweet(client *esv7.Client) *Tweet {
	return &Tweet{
		client: client,
		index:  "tweets",
	}
}

// Index creates or updates a tweet record in an index
func (t *Tweet) Index(message models.ITweet) error {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(message); err != nil {
		return err
	}

	req := esv7api.IndexRequest{
		Index:      t.index,
		Body:       &buf,
		DocumentID: message.Id,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), t.client)
	if err != nil {
		log.Printf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), message.Id)
	}

	io.Copy(io.Discard, res.Body)

	return nil
}
