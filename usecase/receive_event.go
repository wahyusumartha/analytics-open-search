package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wahyusumartha/analytics-open-search/component/message_broker"
)

type EventInput struct {
	ID         string                 `json:"event_id"`
	Name       string                 `json:"event_name"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties"`
}

func (input EventInput) Validate() error {
	if len(input.ID) == 0 {
		return errors.New("Invalid Request: ID is Required")
	}

	if len(input.Name) == 0 {
		return errors.New("Invalid Request: Name is Required")
	}

	if input.Timestamp == 0 {
		return errors.New("Invalid Request: Timestamp is Required")
	}
	return nil
}

type ReceiveEvent interface {
	Execute(ctx context.Context, input EventInput) (string, error)
}

func NewReceiveEventUseCase(
	queueUrl string,
	publisher message_broker.Publisher[
		message_broker.Message,
		message_broker.MessageOutput,
	],
) ReceiveEventUseCase {
	return ReceiveEventUseCase{
		QueueUrl:     queueUrl,
		SQSPublisher: publisher,
	}
}

type ReceiveEventUseCase struct {
	QueueUrl     string
	SQSPublisher message_broker.Publisher[message_broker.Message, message_broker.MessageOutput]
}

func (uc ReceiveEventUseCase) Execute(ctx context.Context, input EventInput) (string, error) {
	err := input.Validate()
	if err != nil {
		return "", err
	}

	jsonInput, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	jsonInputStr := string(jsonInput)
	msg := message_broker.Message{
		QueueUrl: &uc.QueueUrl,
		Body:     &jsonInputStr,
	}

	_, err = uc.SQSPublisher.Publish(ctx, msg)
	if err != nil {
		return "", err
	}

	return "event submitted", nil
}
