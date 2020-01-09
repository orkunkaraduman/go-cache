package main

import (
	"fmt"
	"sync"
	"time"

	cache "github.com/orkunkaraduman/go-cache"
)

func main() {
	ce := cache.NewCache()
	n := 2097152
	keys := make([]string, n)
	for i := range keys {
		keys[i] = fmt.Sprintf("key %d", i)
	}
	var buf [4096]byte
	for i := range buf {
		buf[i] = byte(i)
	}
	tm := time.Now()
	fmt.Println("set start", time.Now())
	for i := 0; i < n; i++ {
		buf[i%len(buf)] = byte(i)
		b := make([]byte, len(buf))
		copy(b, buf[:])
		ce.Set(keys[i], b)
	}
	fmt.Println("set done", time.Now(), time.Now().Sub(tm))
	time.Sleep(1 * time.Second)
	tm = time.Now()
	fmt.Println("get start", time.Now())
	var wg sync.WaitGroup
	for k := 0; k < 4; k++ {
		wg.Add(1)
		go func(k int) {
			for j := 0; j < n/4; j++ {
				i := j * k
				b := ce.Get(keys[i])
				copy(buf[:], b.([]byte))
				if buf[i%len(buf)] != byte(i) {
					panic("a")
				}
			}
			wg.Done()
		}(k)
	}
	wg.Wait()
	fmt.Println("get done", time.Now(), time.Now().Sub(tm))
}
