package key

import "context"

type Generator interface {
	Generate(ctx context.Context, namespace string) string
}
