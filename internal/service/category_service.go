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
	ErrCategoryExists   = errors.New("category already exists")
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

type CreateCategoryInput struct {
	Name   string
	Color  string
	UserID string
}

type UpdateCategoryInput struct {
	ID    string
	Name  *string
	Color *string
}

func (s *CategoryService) CreateCategory(ctx context.Context, input CreateCategoryInput) (*models.Category, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	categories, err := s.categoryRepo.GetCategories(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	for _, cat := range categories {
		if cat.Name == input.Name {
			return nil, ErrCategoryExists
		}
	}

	category := &models.Category{
		Name:   input.Name,
		Color:  input.Color,
		UserID: input.UserID,
	}

	return s.categoryRepo.CreateCategory(ctx, category)
}

func (s *CategoryService) GetCategories(ctx context.Context, userID string) ([]*models.Category, error) {
	return s.categoryRepo.GetCategories(ctx, userID)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, input UpdateCategoryInput) (*models.Category, error) {
	category, err := s.categoryRepo.GetCategoryInfo(ctx, input.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	updates := &repository.CategoryUpdates{}
	now := time.Now()

	if input.Name != nil {
		categories, err := s.categoryRepo.GetCategories(ctx, category.UserID)
		if err != nil {
			return nil, err
		}
		for _, cat := range categories {
			if cat.ID != input.ID && cat.Name == *input.Name {
				return nil, ErrCategoryExists
			}
		}
		updates.Name = input.Name
	}
	updates.Color = input.Color
	updates.UpdatedAt = &now

	return s.categoryRepo.UpdateCategory(ctx, input.ID, updates)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	_, err := s.categoryRepo.GetCategoryInfo(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrCategoryNotFound
		}
		return err
	}
	return s.categoryRepo.DeleteCategory(ctx, id)
}
