package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/ArmaanKatyal/tweetbit/backend/searchService/models"
	"github.com/ArmaanKatyal/tweetbit/backend/searchService/utils"
	esv7 "github.com/elastic/go-elasticsearch/v7"
	esv7api "github.com/elastic/go-elasticsearch/v7/esapi"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ElasticTweet represents the repository used for interacting with tweet records
type ElasticTweet struct {
	client *esv7.Client
	index  string
}

// NewTweet instantiates a new ElasticTweet repository
func NewElasticTweet(client *esv7.Client) *ElasticTweet {
	return &ElasticTweet{
		client: client,
		index:  "tweets",
	}
}

// Index creates or updates a tweet record in an index
func (et *ElasticTweet) IndexTweet(ctx context.Context, message models.ITweet) error {

	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("searchService.elasticsearch").Start(ctx, "ElasticTweet.IndexTweet")
	defer span.End()

	span.SetAttributes(attribute.Key("tweet_content").String(message.Content))
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(message); err != nil {
		span.SetAttributes(attribute.Key("error").Bool(true))
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}

	req := esv7api.IndexRequest{
		Index:      et.index,
		Body:       &buf,
		DocumentID: message.Id,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, et.client)
	if err != nil {
		span.SetAttributes(attribute.Key("error").Bool(true))
		log.Printf("Error getting response: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "TweetIndexRequest.Do")
	}
	defer res.Body.Close()

	if res.IsError() {
		span.SetAttributes(attribute.Key("error").Bool(true))
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), message.Id)
		return utils.NewErrorf(utils.ErrorCodeUnknown, "TweetIndexRequest.Do %s", res.StatusCode)
	}

	io.Copy(io.Discard, res.Body)

	return nil
}

// DeleteTweet removes a tweet record from an index
func (et *ElasticTweet) DeleteTweet(ctx context.Context, id string) error {
	req := esv7api.DeleteRequest{
		Index:      et.index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, et.client)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return utils.WrapErrorf(err, utils.ErrorCodeUnknown, "TweetDeleteRequest.Do")
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error deleting document ID=%s", res.Status(), id)
		return utils.NewErrorf(utils.ErrorCodeUnknown, "TweetDeleteRequest.Do %s", res.StatusCode)
	}

	io.Copy(io.Discard, res.Body)

	return nil
}

// Search returns tweets matching a query
func (et *ElasticTweet) TweetSearch(ctx context.Context, description *string) ([]models.ITweet, error) {
	if description == nil {
		return nil, nil
	}

	newCtx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("searchService.elasticsearch").Start(ctx, "ElasticTweet.TweetSearch")
	defer span.End()

	should := make([]interface{}, 0, 3)

	if description != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"content": *description,
			},
		})
	}
	// TODO: add more fields to search on

	var query map[string]interface{}
	if len(should) > 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": should[0],
		}
	}

	b, _ := json.Marshal(query)
	span.SetAttributes(attribute.Key("query").String(string(b)))

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		span.SetAttributes(attribute.Key("error").Bool(true))
		return nil, utils.WrapErrorf(err, utils.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}

	req := esv7api.SearchRequest{
		Index: []string{et.index},
		Body:  &buf,
	}
	res, err := req.Do(newCtx, et.client)
	if err != nil {
		span.SetAttributes(attribute.Key("error").Bool(true))
		log.Printf("Error getting response: %s", err)
		return nil, utils.WrapErrorf(err, utils.ErrorCodeUnknown, "TweetSearchRequest.Do")
	}
	defer res.Body.Close()

	if res.IsError() {
		span.SetAttributes(attribute.Key("error").Bool(true))
		log.Printf("[%s] Error getting response: %s", res.Status(), res.String())
		return nil, utils.NewErrorf(utils.ErrorCodeUnknown, "TweetSearchRequest.Do %s", res.StatusCode)
	}

	var hits struct {
		Hits struct {
			Hits []struct {
				Source models.ITweet `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&hits); err != nil {
		span.SetAttributes(attribute.Key("error").Bool(true))
		log.Printf("Error parsing the response body: %s", err)
		return nil, utils.WrapErrorf(err, utils.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	tweets := make([]models.ITweet, len(hits.Hits.Hits))
	for i, hit := range hits.Hits.Hits {
		tweets[i].Id = hit.Source.Id
		tweets[i].Content = hit.Source.Content
		tweets[i].UserId = hit.Source.UserId
		tweets[i].Uuid = hit.Source.Uuid
		tweets[i].CreatedAt = hit.Source.CreatedAt
		tweets[i].LikesCount = hit.Source.LikesCount
		tweets[i].RetweetsCount = hit.Source.RetweetsCount
	}

	return tweets, nil
}
