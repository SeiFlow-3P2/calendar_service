package service

import (
	"context"
	"fmt"
	"log"

	"github.com/SeiFlow-3P2/calendar_service/internal/models"
	"github.com/SeiFlow-3P2/calendar_service/internal/producer"
	"github.com/SeiFlow-3P2/calendar_service/pkg/grpc_client"
)

type NotificationService struct {
	authClient    grpc_client.AuthServiceClient
	notifProducer *producer.NotificationProducer
}

func NewNotificationService(
	authClient grpc_client.AuthServiceClient,
	notifProducer *producer.NotificationProducer,
) (*NotificationService, error) {
	if authClient == nil {
		return nil, fmt.Errorf("authClient cannot be nil")
	}
	if notifProducer == nil {
		return nil, fmt.Errorf("notificationProducer cannot be nil")
	}
	return &NotificationService{
		authClient:    authClient,
		notifProducer: notifProducer,
	}, nil
}

func (s *NotificationService) PrepareAndSendEventNotification(ctx context.Context, event *models.Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}
	if event.UserID == "" {
		log.Printf("Event (ID: %s, Title: '%s') has no UserID, skipping notification.", event.ID, event.Title)
		return nil
	}

	telegramID, err := s.authClient.GetTelegramIDByUserID(ctx, event.UserID)
	if err != nil {
		log.Printf("Error getting TelegramID for UserID %s (EventID: %s) from auth_service: %v. Skipping notification.", event.UserID, event.ID, err)
		return fmt.Errorf("auth_service communication error for UserID %s: %w", event.UserID, err)
	}

	if telegramID == "" {
		log.Printf("TelegramID not found for UserID %s (EventID: %s). Skipping notification.", event.UserID, event.ID)
		return nil
	}

	log.Printf("Attempting to send notification for event: '%s' (ID: %s) to Telegram User (ID from Auth: %s)", event.Title, event.ID, telegramID)
	err = s.notifProducer.SendEventNotification(ctx, event.Title, telegramID)
	if err != nil {
		log.Printf("Error sending event notification via Kafka for EventID %s, TelegramID %s: %v", event.ID, telegramID, err)
		return fmt.Errorf("kafka produce error for EventID %s: %w", event.ID, err)
	}

	log.Printf("Notification successfully queued for event '%s' (ID: %s) to Telegram User (ID from Auth: %s)", event.Title, event.ID, telegramID)
	return nil
}
