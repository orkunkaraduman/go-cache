package main

import (
	"github.com/go-cache/cache"
)

func main() {
	ce := cache.NewCache()
	ce.Close()
}
