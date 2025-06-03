package repository

import (
	"context"
	"time"

	"github.com/SeiFlow-3P2/calendar_service/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error)
	GetCategoryInfo(ctx context.Context, id string) (*models.Category, error)
	GetCategories(ctx context.Context, userID string) ([]*models.Category, error)
	UpdateCategory(ctx context.Context, id string, updates *CategoryUpdates) (*models.Category, error)
	DeleteCategory(ctx context.Context, id string) error
	EnsureIndexes(ctx context.Context) error // Новый метод
}

type CategoryUpdates struct {
	Name      *string    `bson:"name,omitempty"`
	Color     *string    `bson:"color,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

type categoryRepository struct {
	db *mongo.Database
}

func NewCategoryRepository(db *mongo.Database) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error) {
	collection := r.db.Collection("categories")
	category.ID = uuid.New().String()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) GetCategoryInfo(ctx context.Context, id string) (*models.Category, error) {
	collection := r.db.Collection("categories")
	var category models.Category
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetCategories(ctx context.Context, userID string) ([]*models.Category, error) {
	collection := r.db.Collection("categories")
	var categories []*models.Category
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, id string, updates *CategoryUpdates) (*models.Category, error) {
	collection := r.db.Collection("categories")
	updateFields := bson.M{}
	if updates.Name != nil {
		updateFields["name"] = *updates.Name
	}
	if updates.Color != nil {
		updateFields["color"] = *updates.Color
	}
	if updates.UpdatedAt != nil {
		updateFields["updated_at"] = *updates.UpdatedAt
	}

	if len(updateFields) == 0 {
		return r.GetCategoryInfo(ctx, id)
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if err != nil {
		return nil, err
	}
	return r.GetCategoryInfo(ctx, id)
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id string) error {
	collection := r.db.Collection("categories")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *categoryRepository) EnsureIndexes(ctx context.Context) error {
	collection := r.db.Collection("categories")

	// Создаём индекс по полю user_id для оптимизации запросов
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	// Создаём уникальный индекс по комбинации name и user_id
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}, {Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}