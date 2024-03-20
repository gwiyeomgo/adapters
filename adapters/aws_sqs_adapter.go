package adapters

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type AwsSQSClientInterface interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

type AwsSQSAdapter struct {
	client AwsSQSClientInterface
}

func NewAwsSQSAdapter(client AwsSQSClientInterface) *AwsSQSAdapter {
	return &AwsSQSAdapter{client: client}
}

// 팩토리 메서드(Factory Method) 패턴
func (a AwsSQSAdapter) NewSQS() *sqs.SQS {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	cliSQS := sqs.New(sess)
	return cliSQS
}

func (a AwsSQSAdapter) SendMessage(msg, queueURL string) (*sqs.SendMessageOutput, error) {

	sqsMessage := &sqs.SendMessageInput{
		MessageGroupId:         aws.String("2a"),
		MessageDeduplicationId: aws.String("2a"),
		QueueUrl:               aws.String(queueURL),
		MessageBody:            aws.String(msg),
	}

	output, err := a.client.SendMessage(sqsMessage)
	if err != nil {
		return nil, fmt.Errorf("could not send message to queue %v: %v", queueURL, err)
	}
	return output, nil
}
