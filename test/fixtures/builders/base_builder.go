package builders

import (
	"math/rand"
	"time"
)

// Builder provides common builder functionality
type Builder[T any] struct {
	instance *T
}

// NewBuilder creates a new builder for type T
func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{
		instance: new(T),
	}
}

// Build returns the built instance
func (b *Builder[T]) Build() *T {
	return b.instance
}

// With sets a field using a function
func (b *Builder[T]) With(fn func(*T)) *Builder[T] {
	fn(b.instance)
	return b
}

// Random utilities
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomInt(min, max int) int {
	return min + seededRand.Intn(max-min+1)
}

// BaseBuilder provides common builder functionality that can be embedded
type BaseBuilder[T any] struct {
	Builder[T]
}

func (b *BaseBuilder[T]) With(fn func(*T)) *BaseBuilder[T] {
	b.Builder.With(fn)
	return b
}
