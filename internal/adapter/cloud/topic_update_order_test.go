package cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
	"github.com/stretchr/testify/assert"
)

func TestUpdateOrderGetTopicName(t *testing.T) {
	t.Run("Should return topic name", func(t *testing.T) {
		// Arrange
		service := NewUpdateOrderTopicService("test-topic", aws.Config{})

		// Act
		topicName := service.GetTopicName()

		// Assert
		assert.Equal(t, "test-topic", topicName)
	})
}

func TestUpdateOrderUpdateTopicArn(t *testing.T) {
	t.Run("Should return nil when topic is found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "ListTopics",
			Input:         &sns.ListTopicsInput{},
			Output: &sns.ListTopicsOutput{
				Topics: []types.Topic{
					{
						TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:test-topic"),
					},
				},
			},
		})

		service := NewUpdateOrderTopicService("test-topic", *stubber.SdkConfig)

		// Act
		err := service.UpdateTopicArn(ctx)

		// Assert
		assert.NoError(t, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("Should return error when topic is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "ListTopics",
			Input:         &sns.ListTopicsInput{},
			Output: &sns.ListTopicsOutput{
				Topics: []types.Topic{
					{
						TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:another-topic"),
					},
				},
			},
		})

		service := NewUpdateOrderTopicService("test-topic", *stubber.SdkConfig)

		// Act
		err := service.UpdateTopicArn(ctx)

		// Assert
		assert.ErrorIs(t, err, custom_error.ErrTopicNotFound)
		testtools.ExitTest(stubber, t)
	})

	t.Run("Should return error when ListTopics operation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		raiseErr := &testtools.StubError{Err: errors.New("ClientError")}

		stubber.Add(testtools.Stub{
			OperationName: "ListTopics",
			Error:         raiseErr,
		})

		service := NewUpdateOrderTopicService("test-topic", *stubber.SdkConfig)

		// Act
		err := service.UpdateTopicArn(ctx)

		// Assert
		testtools.VerifyError(err, raiseErr, t)
		testtools.ExitTest(stubber, t)
	})
}

func TestUpdateOrderPublishMessage(t *testing.T) {
	t.Run("Should return nil when message is published", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "ListTopics",
			Input:         &sns.ListTopicsInput{},
			Output: &sns.ListTopicsOutput{
				Topics: []types.Topic{
					{
						TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:test-topic"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "Publish",
			Input: &sns.PublishInput{
				TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:test-topic"),
				Message:  aws.String(`{"message":"test"}`),
			},
			Output: &sns.PublishOutput{
				MessageId: aws.String("1234"),
			},
		})

		service := NewUpdateOrderTopicService("test-topic", *stubber.SdkConfig)

		err := service.UpdateTopicArn(ctx)
		assert.NoError(t, err)

		message := map[string]string{"message": "test"}

		// Act
		resp, err := service.PublishMessage(ctx, message)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "1234", *resp)
		testtools.ExitTest(stubber, t)
	})

	t.Run("Should return error when message is not published", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "ListTopics",
			Input:         &sns.ListTopicsInput{},
			Output: &sns.ListTopicsOutput{
				Topics: []types.Topic{
					{
						TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:test-topic"),
					},
				},
			},
		})

		raiseErr := &testtools.StubError{Err: errors.New("ClientError")}

		stubber.Add(testtools.Stub{
			OperationName: "Publish",
			Input: &sns.PublishInput{
				TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:test-topic"),
				Message:  aws.String(`{"message":"test"}`),
			},
			Error: raiseErr,
		})

		service := NewUpdateOrderTopicService("test-topic", *stubber.SdkConfig)

		err := service.UpdateTopicArn(ctx)
		assert.NoError(t, err)

		message := map[string]string{"message": "test"}

		// Act
		resp, err := service.PublishMessage(ctx, message)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		testtools.VerifyError(err, raiseErr, t)
		testtools.ExitTest(stubber, t)
	})
}
