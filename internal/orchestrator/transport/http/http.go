package http

import "context"

type Calculator interface {
	Evaluate(ctx context.Context, expression string) (float64, error)
}
