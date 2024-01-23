package message_broker

import "context"

type Publisher[Message any, MessageOutput any] interface {
	Publish(ctx context.Context, message Message) (*MessageOutput, error)
}
