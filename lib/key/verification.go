package key

import "context"

type Verifier interface {
	Verify(ctx context.Context, key string) bool
}
