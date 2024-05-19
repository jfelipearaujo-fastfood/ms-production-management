package cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jfelipearaujo-org/ms-production-management/internal/adapter/cloud/mocks"
	service_mocks "github.com/jfelipearaujo-org/ms-production-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetQueueName(t *testing.T) {
	t.Run("Should return queue name", func(t *testing.T) {
		// Arrange
		fakeProcessor := service_mocks.NewMockCreateOrderProductionService[create.CreateOrderProductionInput](t)
		updateOrderTopic := mocks.NewMockTopicService(t)

		service := NewQueueService("test-queue", aws.Config{}, fakeProcessor, updateOrderTopic)

		// Act
		queueName := service.GetQueueName()

		// Assert
		assert.Equal(t, "test-queue", queueName)
	})
}

func TestUpdateQueueUrl(t *testing.T) {
	t.Run("Should return nil when queue is found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		fakeProcessor := service_mocks.NewMockCreateOrderProductionService[create.CreateOrderProductionInput](t)
		updateOrderTopic := mocks.NewMockTopicService(t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor, updateOrderTopic)

		// Act
		err := service.UpdateQueueUrl(ctx)

		// Assert
		assert.NoError(t, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("Should return error when GetQueueUrl operation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		raiseErr := &testtools.StubError{Err: errors.New("ClientError")}

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Error:         raiseErr,
		})

		fakeProcessor := service_mocks.NewMockCreateOrderProductionService[create.CreateOrderProductionInput](t)
		updateOrderTopic := mocks.NewMockTopicService(t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor, updateOrderTopic)

		// Act
		err := service.UpdateQueueUrl(ctx)

		// Assert
		testtools.VerifyError(err, raiseErr, t)
		testtools.ExitTest(stubber, t)
	})
}

func TestStartConsuming(t *testing.T) {
	t.Run("Should start consuming messages", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"c3fdab1b-3c06-4db2-9edc-4760a2429462\",\"items\":[{\"id\": \"cfdab175-1f86-4fb0-9bcb-15f2c58df30c\",\"name\": \"Hamburger\",\"quantity\": 1}]}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
					{
						MessageId:     aws.String("456"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567890"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567891"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		fakeProcessor := service_mocks.NewMockCreateOrderProductionService[create.CreateOrderProductionInput](t)
		updateOrderTopic := mocks.NewMockTopicService(t)

		fakeProcessor.On("Handle", ctx, mock.Anything).
			Return(nil, nil).
			Times(2)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor, updateOrderTopic)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})

	t.Run("Should log error when cannot receive message", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		raiseErr := &testtools.StubError{Err: errors.New("ClientError")}

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5,
			},
			Error: raiseErr,
		})

		fakeProcessor := service_mocks.NewMockCreateOrderProductionService[create.CreateOrderProductionInput](t)
		updateOrderTopic := mocks.NewMockTopicService(t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor, updateOrderTopic)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})

	t.Run("Should log error when cannot unmarshal message", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"c3fdab1b-3c06-4db2-9edc-4760a2429462\",\"items\":[{\"id\": \"cfdab175-1f86-4fb0-9bcb-15f2c58df30c\",\"name\": \"Hamburger\",\"quantity\": 1}]}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("123"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
					{
						MessageId:     aws.String("456"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567890"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567891"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		fakeProcessor := service_mocks.NewMockCreateOrderProductionService[create.CreateOrderProductionInput](t)
		updateOrderTopic := mocks.NewMockTopicService(t)

		fakeProcessor.On("Handle", ctx, mock.Anything).
			Return(nil, nil).
			Times(2)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor, updateOrderTopic)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})

	t.Run("Should log error when cannot process message", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"c3fdab1b-3c06-4db2-9edc-4760a2429462\",\"items\":[{\"id\": \"cfdab175-1f86-4fb0-9bcb-15f2c58df30c\",\"name\": \"Hamburger\",\"quantity\": 1}]}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("123"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
					{
						MessageId:     aws.String("456"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567890"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567891"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		fakeProcessor := service_mocks.NewMockCreateOrderProductionService[create.CreateOrderProductionInput](t)
		updateOrderTopic := mocks.NewMockTopicService(t)

		fakeProcessor.On("Handle", ctx, mock.Anything).
			Return(nil, assert.AnError).
			Times(2)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor, updateOrderTopic)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})

	t.Run("Should log error when cannot delete message", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		raiseErr := &testtools.StubError{Err: errors.New("ClientError")}

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"c3fdab1b-3c06-4db2-9edc-4760a2429462\",\"items\":[{\"id\": \"cfdab175-1f86-4fb0-9bcb-15f2c58df30c\",\"name\": \"Hamburger\",\"quantity\": 1}]}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("123"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Error: raiseErr,
		})

		fakeProcessor := service_mocks.NewMockCreateOrderProductionService[create.CreateOrderProductionInput](t)
		updateOrderTopic := mocks.NewMockTopicService(t)

		fakeProcessor.On("Handle", ctx, mock.Anything).
			Return(nil, nil).
			Once()

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor, updateOrderTopic)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})
}
