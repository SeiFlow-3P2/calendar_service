package models

import "time"

type Event struct {
	ID          string    `json:"id" bson:"_id,omitempty"` 
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	StartTime   time.Time `json:"start_time" bson:"start_time"`
	EndTime     time.Time `json:"end_time" bson:"end_time"`
	Location    string    `json:"location,omitempty" bson:"location,omitempty"`
	CategoryID  string    `json:"category_id" bson:"category_id"` 
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	
}


type CreateEventParams struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location,omitempty"`
	CategoryID  string    `json:"category_id"`
}

type UpdateEventParams struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Location    *string    `json:"location,omitempty"`
	CategoryID  *string    `json:"category_id,omitempty"`
}


type Category struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Color     string    `json:"color" bson:"color"`      
	UserID    string    `json:"user_id" bson:"user_id"` 
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateCategoryParams struct {
	Name   string `json:"name"`
	Color  string `json:"color"`
	UserID string `json:"user_id"`
}

type UpdateCategoryParams struct {
	Name  *string `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
}