package calendar_test

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	calendarpb "calendar_service/pkg/proto/calendar/v1"
)

// mockCalendarService implements calendarpb.CalendarServiceServer for testing
type mockCalendarService struct {
	calendarpb.UnimplementedCalendarServiceServer
}

func (m *mockCalendarService) CreateCalendar(ctx context.Context, req *calendarpb.CreateCalendarRequest) (*calendarpb.CalendarResponse, error) {
	if req.GetName() == "" || req.GetUserId() == "" {
		return nil, errors.New("name and user_id are required")
	}

	return &calendarpb.CalendarResponse{
		Id:        "calendar-123",
		Name:      req.GetName(),
		UserId:    req.GetUserId(),
		CreatedAt: "2023-01-01T00:00:00Z",
		UpdatedAt: "2023-01-01T00:00:00Z",
	}, nil
}

func (m *mockCalendarService) GetCalendars(ctx context.Context, req *calendarpb.GetCalendarsRequest) (*calendarpb.GetCalendarsResponse, error) {
	if req.GetUserId() == "" {
		return nil, errors.New("user_id is required")
	}

	return &calendarpb.GetCalendarsResponse{
		Calendars: []*calendarpb.CalendarResponse{
			{
				Id:        "calendar-1",
				Name:      "Personal",
				UserId:    req.GetUserId(),
				CreatedAt: "2023-01-01T00:00:00Z",
			},
			{
				Id:        "calendar-2",
				Name:      "Work",
				UserId:    req.GetUserId(),
				CreatedAt: "2023-01-01T00:00:00Z",
			},
		},
	}, nil
}

func (m *mockCalendarService) GetCalendarInfo(ctx context.Context, req *calendarpb.GetCalendarInfoRequest) (*calendarpb.CalendarResponse, error) {
	if req.GetId() == "" {
		return nil, errors.New("id is required")
	}

	return &calendarpb.CalendarResponse{
		Id:        req.GetId(),
		Name:      "Test Calendar",
		UserId:    "user-123",
		CreatedAt: "2023-01-01T00:00:00Z",
		UpdatedAt: "2023-01-01T00:00:00Z",
	}, nil
}

func (m *mockCalendarService) UpdateCalendar(ctx context.Context, req *calendarpb.UpdateCalendarRequest) (*calendarpb.CalendarResponse, error) {
	if req.GetId() == "" {
		return nil, errors.New("id is required")
	}

	resp := &calendarpb.CalendarResponse{
		Id:        req.GetId(),
		UpdatedAt: "2023-01-02T00:00:00Z",
	}
	if name := req.GetName(); name != nil {
		resp.Name = name.GetValue()
	}
	return resp, nil
}

func (m *mockCalendarService) DeleteCalendar(ctx context.Context, req *calendarpb.DeleteCalendarRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New("id is required")
	}
	return &emptypb.Empty{}, nil
}

func (m *mockCalendarService) CreateEvent(ctx context.Context, req *calendarpb.CreateEventRequest) (*calendarpb.EventResponse, error) {
	if req.GetTitle() == "" || req.GetCalendarId() == "" {
		return nil, errors.New("title and calendar_id are required")
	}

	return &calendarpb.EventResponse{
		Id:          "event-123",
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		CalendarId:  req.GetCalendarId(),
		CategoryId:  req.GetCategoryId(),
		CreatedAt:   "2023-01-01T00:00:00Z",
	}, nil
}

func (m *mockCalendarService) UpdateEvent(ctx context.Context, req *calendarpb.UpdateEventRequest) (*calendarpb.EventResponse, error) {
	if req.GetId() == "" {
		return nil, errors.New("id is required")
	}

	resp := &calendarpb.EventResponse{
		Id:        req.GetId(),
		UpdatedAt: "2023-01-02T00:00:00Z",
	}
	if title := req.GetTitle(); title != nil {
		resp.Title = title.GetValue()
	}
	if desc := req.GetDescription(); desc != nil {
		resp.Description = desc.GetValue()
	}
	return resp, nil
}

func (m *mockCalendarService) DeleteEvent(ctx context.Context, req *calendarpb.DeleteEventRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New("id is required")
	}
	return &emptypb.Empty{}, nil
}

func (m *mockCalendarService) GetEvents(ctx context.Context, req *calendarpb.GetEventsRequest) (*calendarpb.GetEventsResponse, error) {
	if req.GetCalendarId() == "" {
		return nil, errors.New("calendar_id is required")
	}

	return &calendarpb.GetEventsResponse{
		Events: []*calendarpb.EventResponse{
			{
				Id:         "event-1",
				Title:      "Meeting",
				CalendarId: req.GetCalendarId(),
				CreatedAt:  "2023-01-01T00:00:00Z",
			},
		},
	}, nil
}

func (m *mockCalendarService) CreateCategory(ctx context.Context, req *calendarpb.CreateEventCategoryRequest) (*calendarpb.EventCategoryResponse, error) {
	if req.GetName() == "" || req.GetUserId() == "" {
		return nil, errors.New("name and user_id are required")
	}

	return &calendarpb.EventCategoryResponse{
		Id:     "category-123",
		Name:   req.GetName(),
		Color:  req.GetColor(),
		UserId: req.GetUserId(),
	}, nil
}

func (m *mockCalendarService) UpdateCategory(ctx context.Context, req *calendarpb.UpdateEventCategoryRequest) (*calendarpb.EventCategoryResponse, error) {
	if req.GetId() == "" {
		return nil, errors.New("id is required")
	}

	resp := &calendarpb.EventCategoryResponse{
		Id: req.GetId(),
	}
	if name := req.GetName(); name != nil {
		resp.Name = name.GetValue()
	}
	if color := req.GetColor(); color != nil {
		resp.Color = color.GetValue()
	}
	return resp, nil
}

func (m *mockCalendarService) DeleteCategory(ctx context.Context, req *calendarpb.DeleteEventCategoryRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New("id is required")
	}
	return &emptypb.Empty{}, nil
}

func (m *mockCalendarService) GetCategories(ctx context.Context, req *calendarpb.GetCategoriesRequest) (*calendarpb.GetCategoriesResponse, error) {
	if req.GetUserId() == "" {
		return nil, errors.New("user_id is required")
	}

	return &calendarpb.GetCategoriesResponse{
		Categories: []*calendarpb.EventCategoryResponse{
			{
				Id:     "category-1",
				Name:   "Work",
				UserId: req.GetUserId(),
			},
			{
				Id:     "category-2",
				Name:   "Personal",
				UserId: req.GetUserId(),
			},
		},
	}, nil
}

func TestCalendarService(t *testing.T) {
	ctx := context.Background()
	service := &mockCalendarService{}

	t.Run("CalendarOperations", func(t *testing.T) {
		t.Run("CreateCalendar", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *calendarpb.CreateCalendarRequest
				wantErr bool
			}{
				{
					name: "success",
					req: &calendarpb.CreateCalendarRequest{
						Name:   "Test Calendar",
						UserId: "user-123",
					},
				},
				{
					name:    "missing name",
					req:     &calendarpb.CreateCalendarRequest{UserId: "user-123"},
					wantErr: true,
				},
				{
					name:    "missing user_id",
					req:     &calendarpb.CreateCalendarRequest{Name: "Test Calendar"},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.CreateCalendar(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("CreateCalendar() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && resp.GetName() != tt.req.GetName() {
						t.Errorf("expected calendar name %q, got %q", tt.req.GetName(), resp.GetName())
					}
				})
			}
		})

		t.Run("GetCalendars", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *calendarpb.GetCalendarsRequest
				wantErr bool
			}{
				{
					name: "success",
					req:  &calendarpb.GetCalendarsRequest{UserId: "user-123"},
				},
				{
					name:    "missing user_id",
					req:     &calendarpb.GetCalendarsRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.GetCalendars(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("GetCalendars() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && len(resp.GetCalendars()) == 0 {
						t.Error("expected non-empty calendars list")
					}
				})
			}
		})

		t.Run("UpdateCalendar", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *calendarpb.UpdateCalendarRequest
				wantErr bool
			}{
				{
					name: "update name",
					req: &calendarpb.UpdateCalendarRequest{
						Id:   "calendar-123",
						Name: wrapperspb.String("New Name"),
					},
				},
				{
					name:    "missing id",
					req:     &calendarpb.UpdateCalendarRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.UpdateCalendar(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("UpdateCalendar() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && name := tt.req.GetName(); name != nil && resp.GetName() != name.GetValue() {
						t.Errorf("expected calendar name %q, got %q", name.GetValue(), resp.GetName())
					}
				})
			}
		})
	})

	t.Run("EventOperations", func(t *testing.T) {
		t.Run("CreateEvent", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *calendarpb.CreateEventRequest
				wantErr bool
			}{
				{
					name: "success",
					req: &calendarpb.CreateEventRequest{
						Title:      "Team Meeting",
						CalendarId: "calendar-123",
					},
				},
				{
					name:    "missing title",
					req:     &calendarpb.CreateEventRequest{CalendarId: "calendar-123"},
					wantErr: true,
				},
				{
					name:    "missing calendar_id",
					req:     &calendarpb.CreateEventRequest{Title: "Team Meeting"},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.CreateEvent(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && resp.GetTitle() != tt.req.GetTitle() {
						t.Errorf("expected event title %q, got %q", tt.req.GetTitle(), resp.GetTitle())
					}
				})
			}
		})

		t.Run("GetEvents", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *calendarpb.GetEventsRequest
				wantErr bool
			}{
				{
					name: "success",
					req:  &calendarpb.GetEventsRequest{CalendarId: "calendar-123"},
				},
				{
					name:    "missing calendar_id",
					req:     &calendarpb.GetEventsRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.GetEvents(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("GetEvents() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && len(resp.GetEvents()) == 0 {
						t.Error("expected non-empty events list")
					}
				})
			}
		})
	})

	t.Run("CategoryOperations", func(t *testing.T) {
		t.Run("CreateCategory", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *calendarpb.CreateEventCategoryRequest
				wantErr bool
			}{
				{
					name: "success",
					req: &calendarpb.CreateEventCategoryRequest{
						Name:   "Work",
						UserId: "user-123",
					},
				},
				{
					name:    "missing name",
					req:     &calendarpb.CreateEventCategoryRequest{UserId: "user-123"},
					wantErr: true,
				},
				{
					name:    "missing user_id",
					req:     &calendarpb.CreateEventCategoryRequest{Name: "Work"},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.CreateCategory(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("CreateCategory() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && resp.GetName() != tt.req.GetName() {
						t.Errorf("expected category name %q, got %q", tt.req.GetName(), resp.GetName())
					}
				})
			}
		})

		t.Run("GetCategories", func(t *testing.T) {
			tests := []struct {
				name    string
				req     *calendarpb.GetCategoriesRequest
				wantErr bool
			}{
				{
					name: "success",
					req:  &calendarpb.GetCategoriesRequest{UserId: "user-123"},
				},
				{
					name:    "missing user_id",
					req:     &calendarpb.GetCategoriesRequest{},
					wantErr: true,
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					resp, err := service.GetCategories(ctx, tt.req)
					if (err != nil) != tt.wantErr {
						t.Errorf("GetCategories() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && len(resp.GetCategories()) == 0 {
						t.Error("expected non-empty categories list")
					}
				})
			}
		})
	})
}