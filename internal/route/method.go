package route

import (
	"context"
	"geogracom-test/pkg/kernel"
	"net/http"
)

//go:generate ifacemaker -f *.go -o imethod.go -i Method -s method -p route
type method struct {
	repo Repository
}

func NewMethod(repo Repository) Method {
	return &method{repo}
}

func (m *method) Create(ctx context.Context, param *Upsert) (any, int, error) {
	routeID, err := m.repo.Upsert(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	statusCode := http.StatusOK

	if routeID != param.ID {
		statusCode = http.StatusAlreadyReported
	}

	return CreateResp{ID: routeID}, statusCode, nil
}

func (m *method) Get(ctx context.Context, param *Get) (any, int, error) {
	route, err := m.repo.Get(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	if !route.Actual {
		return kernel.Response{Description: "route is not relevant"}, http.StatusGone, nil
	}
	return route, http.StatusOK, nil
}

func (m *method) Delete(_ context.Context, param *Delete) (any, int, error) {
	m.repo.DeleteBackground(param)
	return kernel.Response{Description: "route deletion accepted for processing"}, http.StatusAccepted, nil
}
