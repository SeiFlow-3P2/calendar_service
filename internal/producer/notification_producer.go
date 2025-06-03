package producer

import (
	"context"
	"time"

	"github.com/SeiFlow-3P2/calendar_service/pkg/kafka"
)

type NotificationMessage struct {
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

type NotificationProducer struct {
	kafkaProducer kafka.Producer
	topic         string
}

func NewNotificationProducer(kafkaProducer kafka.Producer, topic string) *NotificationProducer {
	return &NotificationProducer{
		kafkaProducer: kafkaProducer,
		topic:         topic,
	}
}

func (p *NotificationProducer) SendEventNotification(ctx context.Context, eventTitle string, telegramID string) error {
	message := NotificationMessage{
		Text:   eventTitle,
		UserID: telegramID,
	}
	return p.kafkaProducer.Produce(ctx, message, p.topic, telegramID, 5*time.Second)
}
