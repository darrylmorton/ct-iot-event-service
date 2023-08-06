package app

import (
	"context"
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/darrylmorton/ct-iot-event-service/internal/data"
	"log"
)

type SQSReceiveMessageImpl struct{}

type ServiceConfig struct {
	EnvConfig              EnvConfig
	SqsReceiveMessageInput *sqs.ReceiveMessageInput
	SqsClient              SQSReceiveMessageAPI
	DbClient               *sql.DB
	Logger                 *log.Logger
	Models                 data.Models
}

type SQSReceiveMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

func GetQueueURL(c context.Context, api SQSReceiveMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func GetMessages(c context.Context, api SQSReceiveMessageAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}
