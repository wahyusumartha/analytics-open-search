package usecase

import (
	"context"
	"errors"
	mock_repository "github.com/wahyusumartha/analytics-open-search/mock/repository"
	"github.com/wahyusumartha/analytics-open-search/repository"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestIngestEventUseCaseExecute(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"should return error when input invalid":       testInvalidInput,
		"should return error when failed to save data": testErrorSaveEvent,
		"should successfully save data":                testSuccessSaveEvent,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testInvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockEventRepository(ctrl)

	input := EventInput{Name: "event_clicked"}

	sut := NewIngestEventUseCase(mockRepo)
	err := sut.Execute(context.TODO(), input)

	assert.Error(t, err, "Invalid Request: ID is Required")
}

func testErrorSaveEvent(t *testing.T) {
	input := EventInput{
		ID:         "random-id",
		Name:       "event-clicked",
		Timestamp:  time.Now().UnixMilli(),
		Properties: nil,
	}

	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockEventRepository(ctrl)

	event := repository.Event{
		ID:         input.ID,
		Name:       input.Name,
		CreatedAt:  input.Timestamp,
		Properties: input.Properties,
	}
	mockRepo.
		EXPECT().Save(gomock.Any(), event).
		Return(errors.New("failed save data")).
		Times(1)

	sut := NewIngestEventUseCase(mockRepo)
	err := sut.Execute(context.TODO(), input)

	assert.Error(t, err, "failed save data")
}

func testSuccessSaveEvent(t *testing.T) {
	input := EventInput{
		ID:         "random-id",
		Name:       "event-clicked",
		Timestamp:  time.Now().UnixMilli(),
		Properties: nil,
	}

	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockEventRepository(ctrl)

	event := repository.Event{
		ID:         input.ID,
		Name:       input.Name,
		CreatedAt:  input.Timestamp,
		Properties: input.Properties,
	}
	mockRepo.
		EXPECT().Save(gomock.Any(), event).
		Return(nil).
		Times(1)

	sut := NewIngestEventUseCase(mockRepo)
	err := sut.Execute(context.TODO(), input)

	assert.NilError(t, err)
}
