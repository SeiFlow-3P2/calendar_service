package repository

import (
	"context"
	"time"

	"github.com/SeiFlow-3P2/calendar_service/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	GetEventInfo(ctx context.Context, id string) (*models.Event, error)
	GetEvents(ctx context.Context) ([]*models.Event, error)
	UpdateEvent(ctx context.Context, id string, updates *EventUpdates) (*models.Event, error)
	DeleteEvent(ctx context.Context, id string) error
}

type EventUpdates struct {
	Title       *string    `bson:"title,omitempty"`
	Description *string    `bson:"description,omitempty"`
	StartTime   *time.Time `bson:"start_time,omitempty"`
	EndTime     *time.Time `bson:"end_time,omitempty"`
	Location    *string    `bson:"location,omitempty"`
	CategoryID  *string    `bson:"category_id,omitempty"`
	UpdatedAt   *time.Time `bson:"updated_at,omitempty"`
}

type eventRepository struct {
	db *mongo.Database
}

func NewEventRepository(db *mongo.Database) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	collection := r.db.Collection("events")
	event.ID = uuid.New().String()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *eventRepository) GetEventInfo(ctx context.Context, id string) (*models.Event, error) {
	collection := r.db.Collection("events")
	var event models.Event
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetEvents(ctx context.Context) ([]*models.Event, error) {
	collection := r.db.Collection("events")
	var events []*models.Event
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var event models.Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) UpdateEvent(ctx context.Context, id string, updates *EventUpdates) (*models.Event, error) {
	collection := r.db.Collection("events")
	updateFields := bson.M{}
	if updates.Title != nil {
		updateFields["title"] = *updates.Title
	}
	if updates.Description != nil {
		updateFields["description"] = *updates.Description
	}
	if updates.StartTime != nil {
		updateFields["start_time"] = *updates.StartTime
	}
	if updates.EndTime != nil {
		updateFields["end_time"] = *updates.EndTime
	}
	if updates.Location != nil {
		updateFields["location"] = *updates.Location
	}
	if updates.CategoryID != nil {
		updateFields["category_id"] = *updates.CategoryID
	}
	if updates.UpdatedAt != nil {
		updateFields["updated_at"] = *updates.UpdatedAt
	}

	if len(updateFields) == 0 {
		return r.GetEventInfo(ctx, id)
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if err != nil {
		return nil, err
	}
	return r.GetEventInfo(ctx, id)
}

func (r *eventRepository) DeleteEvent(ctx context.Context, id string) error {
	collection := r.db.Collection("events")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
