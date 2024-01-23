package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	mockusecase "github.com/wahyusumartha/analytics-open-search/mock/usecase"
	"github.com/wahyusumartha/analytics-open-search/usecase"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestHandleRequest(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"empty request body":          testEmptyRequestBody,
		"invalid request body":        testInvalidRequestBody,
		"event successfully received": testPublishEvent,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testEmptyRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockReceiveEventUseCase := mockusecase.NewMockReceiveEvent(ctrl)
	receiveEventUseCase = mockReceiveEventUseCase

	response, _ := HandleRequest(context.TODO(), events.APIGatewayProxyRequest{Body: ""})
	assert.Equal(t, response.StatusCode, 400)
}

func testInvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockReceiveEventUseCase := mockusecase.NewMockReceiveEvent(ctrl)
	receiveEventUseCase = mockReceiveEventUseCase

	body := usecase.EventInput{
		Name:      "navigation_clicked",
		Timestamp: time.Now().UnixMilli(),
		Properties: map[string]interface{}{
			"screen_name": "home_screen",
			"user_id":     12345,
		},
	}

	mockReceiveEventUseCase.
		EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return("", errors.New("Invalid Request")).
		Times(1)

	jsonBody, _ := json.Marshal(&body)
	jsonBodyStr := string(jsonBody)

	resp, _ := HandleRequest(context.TODO(), events.APIGatewayProxyRequest{Body: jsonBodyStr})
	assert.Equal(t, resp.StatusCode, 400)
}

func testPublishEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockReceiveEventUseCase := mockusecase.NewMockReceiveEvent(ctrl)
	receiveEventUseCase = mockReceiveEventUseCase

	body := usecase.EventInput{
		ID:        "random-request-id",
		Name:      "navigation_clicked",
		Timestamp: time.Now().UnixMilli(),
		Properties: map[string]interface{}{
			"screen_name": "home_screen",
			"user_id":     12345,
		},
	}

	mockReceiveEventUseCase.
		EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return("event submitted", nil).
		Times(1)

	jsonBody, _ := json.Marshal(&body)
	jsonBodyStr := string(jsonBody)

	resp, _ := HandleRequest(context.TODO(), events.APIGatewayProxyRequest{Body: jsonBodyStr})
	assert.Equal(t, resp.StatusCode, 200)
}
