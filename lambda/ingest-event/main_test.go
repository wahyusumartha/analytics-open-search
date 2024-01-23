package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	mock_usecase "github.com/wahyusumartha/analytics-open-search/mock/usecase"
	"github.com/wahyusumartha/analytics-open-search/usecase"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestHandleEvent(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"should put back failed message when input invalid":     testInvalidInput,
		"should put back failed message when failed to execute": testIngestEventUseCaseExecuteFailed,
		"should process all messages successfully":              testAllConsumeMessageSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testInvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)

	useCaseMock := mock_usecase.NewMockIngestEvent(ctrl)
	ingestEventUseCase = useCaseMock

	message := events.SQSMessage{
		MessageId: "12345",
	}
	event := events.SQSEvent{
		Records: []events.SQSMessage{message},
	}

	resp, err := HandleEvent(context.TODO(), event)

	expectedBatchItemFailures := []events.SQSBatchItemFailure{
		{
			ItemIdentifier: "12345",
		},
	}

	assert.Equal(
		t,
		resp.BatchItemFailures[0].ItemIdentifier,
		expectedBatchItemFailures[0].ItemIdentifier,
	)
	assert.NilError(t, err)
}

func testIngestEventUseCaseExecuteFailed(t *testing.T) {
	input := usecase.EventInput{
		ID:         "12345",
		Name:       "event_clicked",
		Timestamp:  time.Now().UnixMilli(),
		Properties: nil,
	}
	jsonInput, _ := json.Marshal(input)

	ctrl := gomock.NewController(t)

	useCaseMock := mock_usecase.NewMockIngestEvent(ctrl)

	firstExecution := useCaseMock.
		EXPECT().
		Execute(
			gomock.Any(),
			input,
		).Return(errors.New("unknown error")).
		Times(1)

	secondExecution := useCaseMock.
		EXPECT().
		Execute(
			gomock.Any(),
			input,
		).Return(nil).
		Times(1)

	gomock.InOrder(firstExecution, secondExecution)

	ingestEventUseCase = useCaseMock

	firstMsg := events.SQSMessage{
		MessageId: "12345",
		Body:      string(jsonInput),
	}

	secondMsg := events.SQSMessage{
		MessageId: "12346",
		Body:      string(jsonInput),
	}
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			firstMsg,
			secondMsg,
		},
	}

	resp, err := HandleEvent(context.TODO(), event)

	expectedBatchItemFailures := []events.SQSBatchItemFailure{
		{
			ItemIdentifier: "12345",
		},
	}

	assert.Equal(
		t,
		resp.BatchItemFailures[0].ItemIdentifier,
		expectedBatchItemFailures[0].ItemIdentifier,
	)
	assert.NilError(t, err)
}

func testAllConsumeMessageSuccess(t *testing.T) {
	input := usecase.EventInput{
		ID:         "12345",
		Name:       "event_clicked",
		Timestamp:  time.Now().UnixMilli(),
		Properties: nil,
	}
	jsonInput, _ := json.Marshal(input)

	ctrl := gomock.NewController(t)

	useCaseMock := mock_usecase.NewMockIngestEvent(ctrl)

	firstExecution := useCaseMock.
		EXPECT().
		Execute(
			gomock.Any(),
			input,
		).Return(nil).
		Times(1)

	secondExecution := useCaseMock.
		EXPECT().
		Execute(
			gomock.Any(),
			input,
		).Return(nil).
		Times(1)

	gomock.InOrder(firstExecution, secondExecution)

	ingestEventUseCase = useCaseMock

	firstMsg := events.SQSMessage{
		MessageId: "12345",
		Body:      string(jsonInput),
	}

	secondMsg := events.SQSMessage{
		MessageId: "12346",
		Body:      string(jsonInput),
	}
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			firstMsg,
			secondMsg,
		},
	}

	resp, err := HandleEvent(context.TODO(), event)
	assert.Equal(
		t,
		len(resp.BatchItemFailures),
		0,
	)
	assert.NilError(t, err)
}
