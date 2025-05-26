package api

import (
	"context"

	"github.com/SeiFlow-3P2/calendar_service/internal/service"

	pb "calendar_service/pkg/proto/calendar/v1"

	"github.com/SeiFlow-3P2/calendar_service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryServiceHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryServiceHandler(categoryService *service.CategoryService) *CategoryServiceHandler {
	return &CategoryServiceHandler{categoryService: categoryService}
}

func (h *CategoryServiceHandler) categoryToResponse(category *models.Category) *pb.EventCategoryResponse {
	return &pb.EventCategoryResponse{
		Id:     category.ID,
		Name:   category.Name,
		Color:  category.Color,
		UserId: category.UserID,
	}
}

func (h *CategoryServiceHandler) CreateCategory(ctx context.Context, req *pb.CreateEventCategoryRequest) (*pb.EventCategoryResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	params := service.CreateCategoryInput{
		Name:   req.Name,
		Color:  req.Color,
		UserID: req.UserId,
	}

	category, err := h.categoryService.CreateCategory(ctx, params)
	if err != nil {
		if err == service.ErrCategoryExists {
			return nil, status.Error(codes.AlreadyExists, "category already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.categoryToResponse(category), nil
}

func (h *CategoryServiceHandler) UpdateCategory(ctx context.Context, req *pb.UpdateEventCategoryRequest) (*pb.EventCategoryResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "category ID is required")
	}

	updates := service.UpdateCategoryInput{ID: req.Id}
	if req.Name != nil {
		updates.Name = &req.Name.Value
	}
	if req.Color != nil {
		updates.Color = &req.Color.Value
	}

	category, err := h.categoryService.UpdateCategory(ctx, updates)
	if err != nil {
		if err == service.ErrCategoryNotFound {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.categoryToResponse(category), nil
}

func (h *CategoryServiceHandler) DeleteCategory(ctx context.Context, req *pb.DeleteEventCategoryRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "category ID is required")
	}

	err := h.categoryService.DeleteCategory(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *CategoryServiceHandler) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	categories, err := h.categoryService.GetCategories(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetCategoriesResponse{
		Categories: make([]*pb.EventCategoryResponse, 0, len(categories)),
	}
	for _, category := range categories {
		response.Categories = append(response.Categories, h.categoryToResponse(category))
	}

	return response, nil
}
