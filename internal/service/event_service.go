package service

import (
	"context"
	"errors"
	"time"

	"github.com/SeiFlow-3P2/calendar_service/internal/repository"

	"github.com/SeiFlow-3P2/calendar_service/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type EventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}

type CreateEventInput struct {
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Location    string
	CategoryID  string
}

type UpdateEventInput struct {
	ID          string
	Title       *string
	Description *string
	StartTime   *time.Time
	EndTime     *time.Time
	Location    *string
	CategoryID  *string
}

func (s *EventService) CreateEvent(ctx context.Context, input CreateEventInput) (*models.Event, error) {
	if input.Title == "" {
		return nil, errors.New("title is required")
	}
	if input.StartTime.IsZero() || input.EndTime.IsZero() {
		return nil, errors.New("start_time and end_time are required")
	}
	if input.StartTime.After(input.EndTime) {
		return nil, errors.New("start_time must be before end_time")
	}

	event := &models.Event{
		Title:       input.Title,
		Description: input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Location:    input.Location,
		CategoryID:  input.CategoryID,
	}

	return s.eventRepo.CreateEvent(ctx, event)
}

func (s *EventService) GetEvents(ctx context.Context) ([]*models.Event, error) {
	return s.eventRepo.GetEvents(ctx)
}

func (s *EventService) UpdateEvent(ctx context.Context, input UpdateEventInput) (*models.Event, error) {
	event, err := s.eventRepo.GetEventInfo(ctx, input.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	updates := &repository.EventUpdates{
		Title:       input.Title,
		Description: input.Description,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Location:    input.Location,
		CategoryID:  input.CategoryID,
		UpdatedAt:   new(time.Time),
	}
	*updates.UpdatedAt = time.Now()

	if updates.Title == nil {
		updates.Title = &event.Title
	}

	if updates.Description == nil {
		updates.Description = &event.Description
	}

	if updates.StartTime == nil {
		updates.StartTime = &event.StartTime
	}
	if updates.EndTime == nil {
		updates.EndTime = &event.EndTime
	}
	if updates.Location == nil {
		updates.Location = &event.Location
	}
	if updates.CategoryID == nil {
		updates.CategoryID = &event.CategoryID
	}

	if updates.StartTime != nil && updates.EndTime != nil && updates.StartTime.After(*updates.EndTime) {
		return nil, errors.New("start_time must be before end_time")
	}

	updatedEvent, err := s.eventRepo.UpdateEvent(ctx, input.ID, updates)
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (s *EventService) DeleteEvent(ctx context.Context, id string) error {
	_, err := s.eventRepo.GetEventInfo(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrEventNotFound
		}
		return err
	}
	return s.eventRepo.DeleteEvent(ctx, id)
}
