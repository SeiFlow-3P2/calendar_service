package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SeiFlow-3P2/calendar_service/internal/models"
	"github.com/SeiFlow-3P2/calendar_service/internal/service"
)

type EventNotifier struct {
	eventService        *service.EventService
	notificationService *service.NotificationService
	interval            time.Duration
	ticker              *time.Ticker
	done                chan struct{}
	wg                  sync.WaitGroup
}

func NewEventNotifier(
	eventService *service.EventService,
	notificationService *service.NotificationService,
	checkInterval time.Duration,
) (*EventNotifier, error) {
	if eventService == nil {
		return nil, fmt.Errorf("eventService cannot be nil")
	}
	if notificationService == nil {
		return nil, fmt.Errorf("notificationService cannot be nil")
	}
	if checkInterval <= 0 {
		checkInterval = 1 * time.Minute
	}

	return &EventNotifier{
		eventService:        eventService,
		notificationService: notificationService,
		interval:            checkInterval,
		done:                make(chan struct{}),
	}, nil
}

func (en *EventNotifier) Start(ctx context.Context) {
	if en.ticker != nil {
		log.Println("EventNotifier already started.")
		return
	}
	en.ticker = time.NewTicker(en.interval)
	en.wg.Add(1)

	log.Printf("EventNotifier scheduler started with check interval: %s", en.interval)

	go func() {
		defer en.wg.Done()
		defer en.ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("EventNotifier: context cancelled, stopping scheduler...")
				return
			case <-en.done:
				log.Println("EventNotifier: stop signal received, stopping scheduler...")
				return
			case tickTime := <-en.ticker.C:
				log.Printf("EventNotifier: tick at %s. Checking for events to notify...", tickTime.Format(time.RFC3339))
				en.processEvents(ctx)
			}
		}
	}()
}

func (en *EventNotifier) processEvents(ctx context.Context) {
	allEvents, err := en.eventService.GetEvents(ctx)
	if err != nil {
		log.Printf("EventNotifier: Error fetching events: %v", err)
		return
	}

	if len(allEvents) == 0 {
		return
	}

	now := time.Now().UTC()
	var eventsToNotify []*models.Event

	for _, event := range allEvents {
		if !event.EndTime.IsZero() && now.After(event.EndTime) {
			eventsToNotify = append(eventsToNotify, event)
		}
	}

	if len(eventsToNotify) == 0 {
		return
	}

	log.Printf("EventNotifier: Found %d event(s) whose EndTime has passed. Attempting to send notifications.", len(eventsToNotify))

	for _, event := range eventsToNotify {
		notifCtx, cancelNotif := context.WithTimeout(ctx, 30*time.Second)

		err := en.notificationService.PrepareAndSendEventNotification(notifCtx, event)
		if err != nil {
			log.Printf("EventNotifier: Failed to send notification for event ID %s (Title: '%s'): %v", event.ID, event.Title, err)
		}
		cancelNotif()
	}
}

func (en *EventNotifier) Stop() {
	if en.ticker == nil {
		log.Println("EventNotifier not running or already stopped.")
		return
	}
	log.Println("EventNotifier: initiating stop...")
	close(en.done)
	en.wg.Wait()
	en.ticker = nil
	log.Println("EventNotifier: scheduler stopped.")
}
