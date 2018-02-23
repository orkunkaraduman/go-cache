package cache

import (
	"strings"

	"github.com/google/btree"
)

type item struct {
	Key string
	Val *Value
}

func (a item) Less(b btree.Item) bool {
	if c, ok := b.(item); ok {
		result := strings.Compare(a.Key, c.Key) < 0
		return result
	}
	return false
}
