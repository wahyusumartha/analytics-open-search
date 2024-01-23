package message_broker

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func NewSQSService(
	ctx context.Context,
	credential *Credential,
	roleArn string,
) *Service {
	cfg, err := config.LoadDefaultConfig(ctx)
	if credential != nil {
		credentialProvider := credentials.NewStaticCredentialsProvider(
			credential.Key,
			credential.Secret,
			credential.Session,
		)
		cfgProvider := config.WithCredentialsProvider(credentialProvider)
		cfg, err = config.LoadDefaultConfig(ctx, cfgProvider)
	}

	roleClient := sts.NewFromConfig(cfg)
	roleProvider := stscreds.NewAssumeRoleProvider(roleClient, roleArn)
	cfg.Credentials = roleProvider

	if err != nil {
		fmt.Printf("Could not load AWS Config: %s", err.Error())
	}

	sqsClient := sqs.NewFromConfig(cfg)
	return &Service{sqsClient}
}

type Credential struct {
	Key     string
	Secret  string
	Session string
}

type Message struct {
	QueueUrl        *string
	DeduplicationID *string
	Body            *string
}

type MessageOutput struct {
	MessageID        *string
	SequenceNumber   *string
	MD5OfMessageBody *string
}

type Service struct {
	client *sqs.Client
}

func (s *Service) Publish(ctx context.Context, message Message) (*MessageOutput, error) {
	msgInput := sqs.SendMessageInput{
		MessageBody:            message.Body,
		QueueUrl:               message.QueueUrl,
		MessageDeduplicationId: message.DeduplicationID,
	}
	msgOutput, err := s.client.SendMessage(ctx, &msgInput)
	if err != nil {
		return nil, err
	}

	output := MessageOutput{
		MessageID:        msgOutput.MessageId,
		SequenceNumber:   msgOutput.SequenceNumber,
		MD5OfMessageBody: msgOutput.MD5OfMessageBody,
	}

	return &output, nil
}
