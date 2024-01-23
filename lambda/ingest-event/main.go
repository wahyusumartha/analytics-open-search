package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/opensearch-project/opensearch-go"
	"github.com/wahyusumartha/analytics-open-search/component/infra"
	"github.com/wahyusumartha/analytics-open-search/repository"
	"github.com/wahyusumartha/analytics-open-search/usecase"
	"log"
	"os"
)

var (
	openSearchClient   *opensearch.Client
	repo               *repository.OpenSearchEventRepository
	ingestEventUseCase usecase.IngestEvent
)

func init() {
	var err error
	openSearchClient, err = infra.NewOpenSearchClient(
		os.Getenv("opensearch_username"),
		os.Getenv("opensearch_password"),
		os.Getenv("opensearch_endpoint"),
	)

	if err != nil {
		log.Print("Could not connect to open search client")
	}

	repo = repository.NewOpenSearchEventRepository(openSearchClient)

	ingestEventUseCase = usecase.NewIngestEventUseCase(repo)
}

func HandleEvent(ctx context.Context, event events.SQSEvent) (events.SQSEventResponse, error) {
	var failures []events.SQSBatchItemFailure
	for _, message := range event.Records {
		err := consumeMessage(ctx, message)
		if err != nil {
			failedItem := events.SQSBatchItemFailure{
				message.MessageId,
			}
			failures = append(failures, failedItem)
		}
	}

	resp := events.SQSEventResponse{BatchItemFailures: failures}
	return resp, nil
}

func consumeMessage(ctx context.Context, message events.SQSMessage) error {
	var event usecase.EventInput

	err := json.Unmarshal([]byte(message.Body), &event)
	if err != nil {
		return err
	}

	return ingestEventUseCase.Execute(ctx, event)
}

func main() {
	lambda.Start(HandleEvent)
}
