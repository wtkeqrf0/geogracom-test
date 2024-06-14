package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	type GreetingRequest struct {
		Name string
	}

	handler := func(ctx context.Context, param *GreetingRequest) (any, int, error) {
		return nil, http.StatusOK, nil
	}

	parser := func(c *fiber.Ctx, out any) error {
		return c.BodyParser(out)
	}

	api := API[GreetingRequest](handler).Handler(parser)

	testCases := []struct {
		name       string
		setContent bool
		statusCode int
		GreetingRequest
	}{
		{
			name:       "422 status",
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name:            "ok status",
			setContent:      true,
			GreetingRequest: GreetingRequest{Name: "Lori"},
			statusCode:      http.StatusOK,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("", api)

			buf := new(bytes.Buffer)
			_ = json.NewEncoder(buf).Encode(tt.GreetingRequest)

			req := httptest.NewRequest("GET", "http://localhost:8080", buf)
			if tt.setContent {
				req.Header.Set("Content-Type", "application/json")
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err.Error())
			}

			if resp.StatusCode != tt.statusCode {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.statusCode)
			}
		})
	}
}
