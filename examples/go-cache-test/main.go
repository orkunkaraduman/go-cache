package main

import (
	"fmt"
	"time"

	"github.com/go-cache/cache"
)

func main() {
	//fmt.Println("")
	ce := cache.NewCache()
	ce.Set("a", &cache.Value{"a"})
	ce.Set("d", &cache.Value{"d"})
	ce.Set("f", &cache.Value{"f"})
	time.Sleep(1 * time.Second)
	fmt.Println(ce.Get("d"))
}
