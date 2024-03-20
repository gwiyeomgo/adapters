package adapters

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockSQSClient struct {
	mock.Mock
}

func (m MockSQSClient) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	d := sqs.SendMessageOutput{}
	d.SetMessageId("123214")
	return &d, nil
}
func TestAwsSQSAdapter_SendMessage(t *testing.T) {
	mockSQS := new(MockSQSClient)
	adapter := NewAwsSQSAdapter(mockSQS)

	s, _ := adapter.SendMessage("test", "test-url")

	assert.Equal(t, "123214", *s.MessageId)
}
