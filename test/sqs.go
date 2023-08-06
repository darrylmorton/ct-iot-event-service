package test

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/darrylmorton/ct-iot-event-service/internal/app"
)

type SQSReceiveMessageImpl struct{}

func (dt SQSReceiveMessageImpl) GetQueueUrl(ctx context.Context,
	params *sqs.GetQueueUrlInput,
	optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {

	output := &sqs.GetQueueUrlOutput{
		QueueUrl: aws.String(app.QueueName),
	}

	return output, nil
}

func (dt SQSReceiveMessageImpl) ReceiveMessage(ctx context.Context,
	params *sqs.ReceiveMessageInput,
	optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {

	messageOne := Events[0]
	messageOneMarshalled, _ := json.Marshal(messageOne)

	messageTwo := Events[1]
	messageTwoMarshalled, _ := json.Marshal(messageTwo)

	messages := []types.Message{
		{
			MessageId:     aws.String("message-one-id"),
			ReceiptHandle: aws.String("message-one-receipt-handle"),
			Body:          aws.String(string(messageOneMarshalled)),
		},
		{
			MessageId:     aws.String("message-two-id"),
			ReceiptHandle: aws.String("message-two-receipt-handle"),
			Body:          aws.String(string(messageTwoMarshalled)),
		},
	}

	output := &sqs.ReceiveMessageOutput{
		Messages: messages,
	}

	return output, nil
}
