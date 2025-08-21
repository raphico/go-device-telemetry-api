package telemetry

import "context"

type Repository interface {
	Create(ctx context.Context, t *Telemetry) error
}
