package cloud

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/create"
)

type QueueService interface {
	GetQueueName() string
	UpdateQueueUrl(ctx context.Context) error
	ConsumeMessages(ctx context.Context)
}

type AwsSqsService struct {
	QueueName string
	QueueUrl  string
	Client    *sqs.Client

	MessageProcessor        service.CreateOrderProductionService[create.CreateOrderProductionInput]
	UpdateOrderTopicService TopicService

	ChanMessage chan types.Message

	Mutex     sync.Mutex
	WaitGroup sync.WaitGroup
}

func NewQueueService(
	queueName string,
	config aws.Config,
	messageProcessor service.CreateOrderProductionService[create.CreateOrderProductionInput],
	updateOrderTopicService TopicService,
) QueueService {
	client := sqs.NewFromConfig(config)

	return &AwsSqsService{
		QueueName: queueName,
		Client:    client,

		MessageProcessor:        messageProcessor,
		UpdateOrderTopicService: updateOrderTopicService,

		ChanMessage: make(chan types.Message, 10),

		Mutex:     sync.Mutex{},
		WaitGroup: sync.WaitGroup{},
	}
}

func (s *AwsSqsService) GetQueueName() string {
	return s.QueueName
}

func (s *AwsSqsService) UpdateQueueUrl(ctx context.Context) error {
	output, err := s.Client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &s.QueueName,
	})
	if err != nil {
		return err
	}

	s.QueueUrl = *output.QueueUrl

	return nil
}

func (s *AwsSqsService) ConsumeMessages(ctx context.Context) {
	output, err := s.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &s.QueueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     5,
	})
	if err != nil {
		slog.ErrorContext(ctx, "error receiving message from queue", "queue_url", s.QueueUrl, "error", err)
		return
	}

	s.WaitGroup.Add(len(output.Messages))

	for _, message := range output.Messages {
		go s.processMessage(ctx, message)
	}

	s.WaitGroup.Wait()
}

func (s *AwsSqsService) processMessage(ctx context.Context, message types.Message) {
	defer s.WaitGroup.Done()
	s.Mutex.Lock()

	slog.InfoContext(ctx, "message received", "message_id", *message.MessageId)

	var notification TopicNotification

	err := json.Unmarshal([]byte(*message.Body), &notification)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshalling message", "message_id", *message.MessageId, "error", err)
	}

	if notification.Type != "Notification" {
		slog.ErrorContext(ctx, "invalid notification type", "message_id", *message.MessageId, "type", notification.Type)
		return
	}

	var request create.CreateOrderProductionInput

	err = json.Unmarshal([]byte(notification.Message), &request)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshalling message", "message_id", *message.MessageId, "error", err)
	}

	if err == nil {
		slog.InfoContext(ctx, "message unmarshalled", "request", request)
		order, err := s.MessageProcessor.Handle(ctx, request)
		if err != nil {
			slog.ErrorContext(ctx, "error processing message", "message_id", *message.MessageId, "error", err)
		}

		if order != nil {
			messageId, err := s.UpdateOrderTopicService.PublishMessage(ctx, NewUpdateOrderContractFromPayment(order))
			if err != nil {
				slog.ErrorContext(ctx, "error publishing message to update order topic", "error", err)
			}

			if messageId != nil {
				slog.InfoContext(ctx, "message published to update order topic", "message_id", *messageId)
			}
		}
	}

	if err := s.deleteMessage(ctx, message); err != nil {
		slog.ErrorContext(ctx, "error deleting message", "message_id", *message.MessageId, "error", err)
	}

	s.Mutex.Unlock()
}

func (s *AwsSqsService) deleteMessage(ctx context.Context, message types.Message) error {
	_, err := s.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &s.QueueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return err
	}

	return nil
}
