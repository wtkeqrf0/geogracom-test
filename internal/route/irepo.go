// Code generated by ifacemaker; DO NOT EDIT.

package route

import (
	"context"
)

// Repository ...
type Repository interface {
	Upsert(ctx context.Context, param *Upsert) (routeID int64, err error)
	Get(ctx context.Context, param *Get) (*Route, error)
	Delete(ctx context.Context, param *Delete) error
	DeleteBackground(param *Delete)
}
