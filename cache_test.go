package cache

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkCacheGet(b *testing.B) {
	ce := NewCache()
	keys := make([]string, b.N)
	for i := range keys {
		keys[i] = fmt.Sprintf("key %d", i)
	}
	for i := 0; i < b.N; i++ {
		ce.Set(keys[i], i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ce.Get(keys[i]) == nil {
			b.FailNow()
		}
	}
}

func BenchmarkCacheSet(b *testing.B) {
	ce := NewCache()
	keys := make([]string, b.N)
	for i := range keys {
		keys[i] = fmt.Sprintf("key %d", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ce.Set(keys[i], i)
	}
}

func TestCacheTime(t *testing.T) {
	ce := NewCache()
	n := 2097152
	keys := make([]string, n)
	for i := range keys {
		keys[i] = fmt.Sprintf("key %d", i)
	}
	var buf [4096]byte
	for i := range buf {
		buf[i] = 0xff
	}
	tm := time.Now()
	for i := 0; i < n; i++ {
		b := make([]byte, len(buf))
		copy(b, buf[:])
		ce.Set(keys[i], b)
	}
	t.Log("set", time.Now().Sub(tm))
	tm = time.Now()
	for i := 0; i < n; i++ {
		b := ce.Get(keys[i])
		copy(buf[:], b.([]byte))
		if i == 1000000 {
			fmt.Println(buf)
		}
	}
	t.Log("get", time.Now().Sub(tm))
}
