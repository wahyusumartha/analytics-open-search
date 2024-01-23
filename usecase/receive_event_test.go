package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wahyusumartha/analytics-open-search/component/message_broker"
	mockmessagebroker "github.com/wahyusumartha/analytics-open-search/mock"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

type errorTestCase struct {
	description   string
	input         EventInput
	expectedError string
}

func TestExecuteInvalidInputShouldError(t *testing.T) {
	for _, scenario := range []errorTestCase{
		{
			description: "Invalid Input (ID is Required)",
			input: EventInput{
				Name:      "navigation_clicked",
				Timestamp: time.Now().UnixMilli(),
			},
			expectedError: "Invalid Request: ID is Required",
		},
		{
			description: "Invalid Input (Name is Required)",
			input: EventInput{
				ID:        "random-id-string",
				Timestamp: time.Now().UnixMilli(),
			},
			expectedError: "Invalid Request: Name is Required",
		},
		{
			description: "Invalid Input (Time is Required)",
			input: EventInput{
				ID:   "random-id-string",
				Name: "navigation_clicked",
			},
			expectedError: "Invalid Request: Timestamp is Required",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSQSPublisher := mockmessagebroker.NewMockPublisher[message_broker.Message, message_broker.MessageOutput](ctrl)
			sut := NewReceiveEventUseCase("", mockSQSPublisher)
			_, err := sut.Execute(context.TODO(), scenario.input)
			assert.Error(t, err, scenario.expectedError)
		})
	}
}

func TestExecutePublishSQS(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"failed publish sqs":  testExecutePublishSQSShouldError,
		"success publish sqs": testExecutePublishSQSShouldSucceed,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testExecutePublishSQSShouldSucceed(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSQSPublisher := mockmessagebroker.NewMockPublisher[message_broker.Message, message_broker.MessageOutput](ctrl)

	body := EventInput{
		ID:        "random-id-string",
		Name:      "navigation-clicked",
		Timestamp: time.Now().UnixMilli(),
		Properties: map[string]interface{}{
			"screen_name": "home_screen",
			"user_id":     12345,
		},
	}

	jsonBody, _ := json.Marshal(&body)
	jsonBodyStr := string(jsonBody)
	queueUrl := "queue-url"
	msg := message_broker.Message{
		QueueUrl: &queueUrl,
		Body:     &jsonBodyStr,
	}
	mockSQSPublisher.
		EXPECT().
		Publish(gomock.Any(), msg).
		Return(&message_broker.MessageOutput{MessageID: nil}, nil).
		Times(1)

	sut := NewReceiveEventUseCase(queueUrl, mockSQSPublisher)

	output, err := sut.Execute(context.TODO(), body)
	assert.NilError(t, err)
	assert.Equal(t, output, "event submitted")
}

func testExecutePublishSQSShouldError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSQSPublisher := mockmessagebroker.NewMockPublisher[message_broker.Message, message_broker.MessageOutput](ctrl)

	body := EventInput{
		ID:        "random-id-string",
		Name:      "navigation-clicked",
		Timestamp: time.Now().UnixMilli(),
	}

	jsonBody, _ := json.Marshal(&body)
	jsonBodyStr := string(jsonBody)
	queueUrl := "queue-url"
	msg := message_broker.Message{
		QueueUrl: &queueUrl,
		Body:     &jsonBodyStr,
	}
	mockSQSPublisher.
		EXPECT().
		Publish(gomock.Any(), msg).
		Return(nil, errors.New("failed publish sqs")).
		Times(1)

	sut := NewReceiveEventUseCase(queueUrl, mockSQSPublisher)

	_, err := sut.Execute(context.TODO(), body)
	assert.Error(t, err, "failed publish sqs")
}
