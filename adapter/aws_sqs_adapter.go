package adapter

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewSQS() *sqs.SQS {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	cliSQS := sqs.New(sess)
	return cliSQS
}

func SendMessage(sqsClient *sqs.SQS, msg, queueURL string) (*sqs.SendMessageOutput, error) {
	sqsMessage := &sqs.SendMessageInput{
		MessageGroupId:         aws.String("2a"),
		MessageDeduplicationId: aws.String("2a"),
		QueueUrl:               aws.String(queueURL),
		MessageBody:            aws.String(msg),
	}

	output, err := sqsClient.SendMessage(sqsMessage)
	if err != nil {
		return nil, fmt.Errorf("could not send message to queue %v: %v", queueURL, err)
	}
	return output, nil
}
