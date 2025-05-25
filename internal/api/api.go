package api

import (
	"context"

	pb "calendar_service/pkg/proto/calendar/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedCalendarServiceServer
	eventHandler    *EventServiceHandler
	categoryHandler *CategoryServiceHandler
}

func NewHandler(
	eventHandler *EventServiceHandler,
	categoryHandler *CategoryServiceHandler,
) *Handler {
	return &Handler{
		eventHandler:    eventHandler,
		categoryHandler: categoryHandler,
	}
}

func (h *Handler) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.EventResponse, error) {
	return h.eventHandler.CreateEvent(ctx, req)
}

func (h *Handler) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	return h.eventHandler.UpdateEvent(ctx, req)
}

func (h *Handler) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	return h.eventHandler.DeleteEvent(ctx, req)
}

func (h *Handler) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return h.eventHandler.GetEvents(ctx, req)
}

func (h *Handler) CreateCategory(ctx context.Context, req *pb.CreateEventCategoryRequest) (*pb.EventCategoryResponse, error) {
	return h.categoryHandler.CreateCategory(ctx, req)
}

func (h *Handler) UpdateCategory(ctx context.Context, req *pb.UpdateEventCategoryRequest) (*pb.EventCategoryResponse, error) {
	return h.categoryHandler.UpdateCategory(ctx, req)
}

func (h *Handler) DeleteCategory(ctx context.Context, req *pb.DeleteEventCategoryRequest) (*emptypb.Empty, error) {
	return h.categoryHandler.DeleteCategory(ctx, req)
}

func (h *Handler) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesResponse, error) {
	return h.categoryHandler.GetCategories(ctx, req)
}