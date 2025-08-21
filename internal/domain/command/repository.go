package command

import "context"

type Repository interface {
	Create(ctx context.Context, c *Command) error
}
