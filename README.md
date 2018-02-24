# Go in-memory cache library

[![GoDoc](https://godoc.org/github.com/go-cache/cache?status.svg)](https://godoc.org/github.com/go-cache/cache)

Package cache offers concurrency safe in-memory cache based on b-tree and hash-map indexing.
All methods of Cache struct are concurrency safe and operates cache atomically.

## Usage

### Simple code

```go
package main

import (
	"fmt"

	"github.com/go-cache/cache"
)

func main() {
	fmt.Println("")
	ce := cache.NewCache()
	ce.Set("a", &cache.Value{V: "1"})
	ce.Set("b", &cache.Value{V: "2"})
	ce.Set("d", &cache.Value{V: "4"})
	ce.Set("e", &cache.Value{V: "5"})
	ce.Set("f", &cache.Value{V: "6"})
	ce.Del("e")
	fmt.Printf("%s => %s\n", "a", ce.Get("a"))
	fmt.Printf("%s => %s\n", "b", ce.Get("b"))
	fmt.Printf("%s => %s\n", "c", ce.Get("c"))
	fmt.Printf("%s => %s\n", "d", ce.Get("d"))
	fmt.Printf("%s => %s\n", "e", ce.Get("e"))
	fmt.Printf("%s => %s\n", "f", ce.Get("f"))
	ce.Close()
}
```
