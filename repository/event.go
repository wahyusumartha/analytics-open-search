package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"strings"
)

type Event struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	CreatedAt  int64                  `json:"created_at"`
	Properties map[string]interface{} `json:"properties"`
}

func (e Event) toJSON() string {
	eventByte, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(eventByte)
}

type EventRepository interface {
	Save(ctx context.Context, event Event) error
}

func NewOpenSearchEventRepository(client *opensearch.Client) *OpenSearchEventRepository {
	return &OpenSearchEventRepository{client: client}
}

type OpenSearchEventRepository struct {
	client *opensearch.Client
}

func (r *OpenSearchEventRepository) Save(ctx context.Context, event Event) error {

	isValidBody := len(event.toJSON()) > 0
	if !isValidBody {
		return errors.New("Invalid Body")
	}

	req := opensearchapi.IndexRequest{
		Index:      "event",
		DocumentID: event.ID,
		Body:       strings.NewReader(event.toJSON()),
	}

	_, err := req.Do(ctx, r.client)
	if err != nil {
		return err
	}

	return nil
}
