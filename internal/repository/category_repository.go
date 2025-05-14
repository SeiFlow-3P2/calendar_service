package repository

import (
    "calendar_service/internal/models"
    "context"
)

type CategoryRepository interface {
    Create(ctx context.Context, category *domain.Category) (string, error) 
    GetByID(ctx context.Context, id string) (*domain.Category, error)
    Update(ctx context.Context, id string, params *domain.UpdateCategoryParams) error
    Delete(ctx context.Context, id string) error
    ListByUserID(
        ctx context.Context,
        userID string, 
        limit int,
        offset int,
    ) ([]*domain.Category, int64, error) 
    GetByNameAndUserID(ctx context.Context, name string, userID string) (*domain.Category, error) 
}