package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	EventRepository    EventRepository
	CategoryRepository CategoryRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		EventRepository:    NewEventRepository(db),
		CategoryRepository: NewCategoryRepository(db),
	}
}