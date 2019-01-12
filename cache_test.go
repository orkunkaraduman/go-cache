package cache

import (
	"fmt"
	"testing"
)

func BenchmarkCacheSet(b *testing.B) {
	ce := NewCache()
	b.ResetTimer()
	for i := 0; i < 4; i++ {
		b.Run(fmt.Sprintf("pass %d", i+1), func(b *testing.B) {
			keys := make([]string, b.N)
			for i := range keys {
				keys[i] = fmt.Sprintf("key %d", i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ce.Set(keys[i], i)
			}
		})
	}
}
