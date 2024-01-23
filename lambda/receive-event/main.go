package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wahyusumartha/analytics-open-search/component/message_broker"
	"github.com/wahyusumartha/analytics-open-search/usecase"
	"os"
)

type Response struct {
	Message string `json:"message"`
}

var (
	sqsClient           *message_broker.Service
	receiveEventUseCase usecase.ReceiveEvent
)

func init() {
	sqsClient = message_broker.NewSQSService(
		context.Background(),
		&message_broker.Credential{
			Key:    os.Getenv("aws_key"),
			Secret: os.Getenv("aws_secret"),
		},
		os.Getenv("role_arn"),
	)

	receiveEventUseCase = usecase.NewReceiveEventUseCase(
		os.Getenv("incoming_event_queue_url"),
		sqsClient,
	)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var eventRequest usecase.EventInput
	err := json.Unmarshal([]byte(request.Body), &eventRequest)
	if err != nil {
		return createResponse(400, err.Error()), nil
	}

	output, err := receiveEventUseCase.Execute(ctx, eventRequest)
	if err != nil {
		return createResponse(400, err.Error()), nil
	}

	return createResponse(200, output), nil
}

func createResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	responseBody := &Response{
		Message: message,
	}

	responseBytes, _ := json.Marshal(responseBody)

	response := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(responseBytes),
	}

	return response
}

func main() {
	lambda.Start(HandleRequest)
}
