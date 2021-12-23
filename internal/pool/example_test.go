package pool_test

import (
	"github.com/tier2pool/tier2pool/internal/pool"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool.Get()
	}
}

func BenchmarkPut(b *testing.B) {
	buffer := pool.Get()
	for i := 0; i < b.N; i++ {
		pool.Put(buffer)
	}
}
