package key

type Verifier interface {
	Verify(key string) bool
}
