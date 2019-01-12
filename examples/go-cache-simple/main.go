package main

import (
	"fmt"

	cache "github.com/orkunkaraduman/go-cache"
)

func main() {
	fmt.Println("")
	ce := cache.NewCache()
	ce.Set("a", "1")
	ce.Set("b", "2")
	ce.Set("d", "4")
	ce.Set("e", "5")
	ce.Set("f", "6")
	ce.Del("e")
	fmt.Printf("%s => %s\n", "a", ce.Get("a"))
	fmt.Printf("%s => %s\n", "b", ce.Get("b"))
	fmt.Printf("%s => %s\n", "c", ce.Get("c"))
	fmt.Printf("%s => %s\n", "d", ce.Get("d"))
	fmt.Printf("%s => %s\n", "e", ce.Get("e"))
	fmt.Printf("%s => %s\n", "f", ce.Get("f"))
	ce.Close()
}
