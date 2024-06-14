package route

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/mock/gomock"
	"math/rand/v2"
	"net/http"
	"testing"
)

func TestMethodCreate(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name        string
		modify      func(r *MockRepository)
		shouldError bool
		shouldEqual bool
		*Upsert
	}{
		{
			name: "Insert instantly",
			modify: func(r *MockRepository) {
				r.EXPECT().Upsert(ctx, gomock.Any()).Return(int64(123456), nil)
			},
			shouldEqual: true,
			Upsert: &Upsert{
				ID: 123456,
				Insert: Insert{
					Name:  "Москва-Тюмень",
					Load:  93.234,
					Cargo: "уголь",
				},
			},
		},
		{
			name: "Error case",
			modify: func(r *MockRepository) {
				r.EXPECT().Upsert(ctx, gomock.Any()).Return(int64(0), sql.ErrConnDone)
			},
			shouldError: true,
			Upsert: &Upsert{
				ID: 123456,
				Insert: Insert{
					Name:  "Москва-Тюмень",
					Load:  93.234,
					Cargo: "уголь",
				},
			},
		},
		{
			name: "Insert with other id",
			modify: func(r *MockRepository) {
				r.EXPECT().Upsert(ctx, gomock.Any()).Return(rand.Int64(), nil)
			},
			Upsert: &Upsert{
				ID: 123456,
				Insert: Insert{
					Name:  "Москва-Тюмень",
					Load:  93.234,
					Cargo: "уголь",
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mRepo := NewMockRepository(ctrl)
			tt.modify(mRepo)

			resp, statusCode, err := NewMethod(mRepo).Create(ctx, tt.Upsert)
			if err != nil {
				if tt.shouldError {
					return
				}
				t.Fatal(err.Error())
			}

			r, ok := resp.(CreateResp)
			if !ok {
				t.Fatal("failed to get response")
			}

			if tt.shouldEqual {
				if r.ID == tt.ID && statusCode == http.StatusOK {
					return
				} else {
					t.Error("ids are not equal")
				}
			} else {
				if r.ID == tt.ID {
					t.Error("ids are equal")
				} else if statusCode == http.StatusAlreadyReported {
					return
				}
			}
		})
	}
}

func TestMethodGet(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name        string
		modify      func(r *MockRepository)
		statusCode  int
		shouldError bool
		*Get
	}{
		{
			name: "Get with actual true",
			modify: func(r *MockRepository) {
				r.EXPECT().Get(ctx, gomock.Any()).Return(&Route{
					Actual: true,
				}, nil)
			},
			statusCode: http.StatusOK,
			Get: &Get{
				ID: 123456,
			},
		},
		{
			name: "Error case",
			modify: func(r *MockRepository) {
				r.EXPECT().Get(ctx, gomock.Any()).Return(nil, sql.ErrConnDone)
			},
			shouldError: true,
			Get: &Get{
				ID: 123456,
			},
		},
		{
			name: "Get with actual false",
			modify: func(r *MockRepository) {
				r.EXPECT().Get(ctx, gomock.Any()).Return(&Route{
					Actual: false,
				}, nil)
			},
			statusCode: http.StatusGone,
			Get: &Get{
				ID: 123456,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mRepo := NewMockRepository(ctrl)
			tt.modify(mRepo)

			_, statusCode, err := NewMethod(mRepo).Get(ctx, tt.Get)
			if err != nil {
				if tt.shouldError {
					return
				}
				t.Fatal(err.Error())
			}

			if tt.statusCode != statusCode {
				t.Errorf("got %d, want %d", statusCode, tt.statusCode)
			}
		})
	}
}

func TestMethodDelete(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name   string
		modify func(r *MockRepository)
		*Delete
	}{
		{
			name: "Delete in background mode",
			modify: func(r *MockRepository) {
				r.EXPECT().DeleteBackground(gomock.Any())
			},
			Delete: &Delete{
				Ids: []int64{123456},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mRepo := NewMockRepository(ctrl)
			tt.modify(mRepo)

			_, statusCode, err := NewMethod(mRepo).Delete(ctx, tt.Delete)
			if err != nil {
				t.Fatal(err.Error())
			}

			if http.StatusAccepted != statusCode {
				t.Errorf("got %d, want %d", statusCode, http.StatusAccepted)
			}
		})
	}
}

func TestMethodMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := fiber.New()

	NewMethod(NewMockRepository(ctrl)).MapHandlers(app)

	if len(app.GetRoutes()) != 4 {
		t.Error("routes are not registered")
	}
}
