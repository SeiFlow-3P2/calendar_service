package repository

import (
    "calendar_service/internal/models"
    "context"
    "time"
)

type EventRepository interface {
    Create(ctx context.Context, event *domain.Event) (string, error) 
    GetByID(ctx context.Context, id string) (*domain.Event, error)
    Update(ctx context.Context, id string, params *domain.UpdateEventParams) error
    Delete(ctx context.Context, id string) error
    List(
        ctx context.Context,
        startTimeFrom *time.Time, 
        startTimeTo *time.Time,   
        categoryID *string,      
        limit int,
        offset int,
    ) ([]*domain.Event, int64, error) 

    GetEventsForNotification(ctx context.Context, notifyTimeThreshold time.Time) ([]*domain.Event, error)

}