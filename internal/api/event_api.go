package api

import (
	"context"
	"time"

	pb "github.com/SeiFlow-3P2/calendar_service/pkg/proto/v1"

	"github.com/SeiFlow-3P2/calendar_service/internal/service"

	"github.com/SeiFlow-3P2/calendar_service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type EventServiceHandler struct {
	eventService *service.EventService
}

func NewEventServiceHandler(eventService *service.EventService) *EventServiceHandler {
	return &EventServiceHandler{eventService: eventService}
}

func (h *EventServiceHandler) eventToResponse(event *models.Event) *pb.EventResponse {
	return &pb.EventResponse{
		Id:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime.Format(time.RFC3339),
		EndTime:     event.EndTime.Format(time.RFC3339),
		Location:    wrapperspb.String(event.Location),
		CategoryId:  event.CategoryID,
		CreatedAt:   event.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   event.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *EventServiceHandler) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.EventResponse, error) {
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.StartTime == "" {
		return nil, status.Error(codes.InvalidArgument, "start_time is required")
	}
	if req.EndTime == "" {
		return nil, status.Error(codes.InvalidArgument, "end_time is required")
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid start_time format")
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid end_time format")
	}

	params := service.CreateEventInput{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		Location:    req.GetLocation().GetValue(),
		CategoryID:  req.CategoryId,
	}

	event, err := h.eventService.CreateEvent(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.eventToResponse(event), nil
}

func (h *EventServiceHandler) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "event ID is required")
	}

	updates := service.UpdateEventInput{ID: req.Id}
	if req.Title != nil {
		updates.Title = &req.Title.Value
	}
	if req.Description != nil {
		updates.Description = &req.Description.Value
	}
	if req.StartTime != nil {
		startTime, err := time.Parse(time.RFC3339, req.StartTime.Value)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid start_time format")
		}
		updates.StartTime = &startTime
	}
	if req.EndTime != nil {
		endTime, err := time.Parse(time.RFC3339, req.EndTime.Value)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid end_time format")
		}
		updates.EndTime = &endTime
	}
	if req.Location != nil {
		updates.Location = &req.Location.Value
	}
	if req.CategoryId != nil {
		updates.CategoryID = &req.CategoryId.Value
	}

	event, err := h.eventService.UpdateEvent(ctx, updates)
	if err != nil {
		if err == service.ErrEventNotFound {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.eventToResponse(event), nil
}

func (h *EventServiceHandler) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "event ID is required")
	}

	err := h.eventService.DeleteEvent(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *EventServiceHandler) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := h.eventService.GetEvents(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetEventsResponse{
		Events: make([]*pb.EventResponse, 0, len(events)),
	}
	for _, event := range events {
		response.Events = append(response.Events, h.eventToResponse(event))
	}

	return response, nil
}
