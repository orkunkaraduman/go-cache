/*
Package cache offers concurrency-safe in-memory cache based on b-tree and hash-map indexing.
*/
package cache

import (
	"runtime"
	"sync"

	"github.com/google/btree"
)

var (
	// DefaultDegree is default b-tree degree.
	DefaultDegree = 4
)

// Cache struct is concurrency-safe in-memory cache based on b-tree and hash-map indexing.
type Cache struct {
	done   chan bool
	tr     *btree.BTree
	trMu   sync.RWMutex
	qu     map[string]item
	quMu   sync.RWMutex
	quCh   chan bool
	degree int
}

// NewCache returns a new Cache has default degree.
func NewCache() (ce *Cache) {
	return NewCacheDegree(DefaultDegree)
}

// NewCacheDegree returns a new Cache given degree.
func NewCacheDegree(degree int) (ce *Cache) {
	ce = &Cache{
		done:   make(chan bool),
		quCh:   make(chan bool, 1),
		degree: degree,
	}
	ce.Flush()
	go ce.queueWorker()
	return
}

// Flush flushes the cache.
func (ce *Cache) Flush() {
	ce.trMu.Lock()
	ce.tr = btree.New(ce.degree)
	ce.quMu.Lock()
	ce.trMu.Unlock()
	ce.qu = make(map[string]item)
	ce.quMu.Unlock()
}

// Close closes the cache. It must be called if the cache will not use.
func (ce *Cache) Close() {
	ce.done <- true
}

func (ce *Cache) queueWorker() {
	for {
		select {
		case <-ce.done:
			return
		case <-ce.quCh:
		}
		for {
			var im item
			var found bool
			ce.quMu.Lock()
			for key := range ce.qu {
				im = ce.qu[key]
				found = true
				delete(ce.qu, key)
				break
			}
			if !found {
				ce.quMu.Unlock()
				break
			}
			ce.trMu.Lock()
			ce.quMu.Unlock()
			if im.Val != nil {
				ce.tr.ReplaceOrInsert(im)
			} else {
				ce.tr.Delete(im)
			}
			ce.trMu.Unlock()
			runtime.Gosched()
		}
	}
}

// Get returns the value of given key. It returns nil, if the key wasn't exist.
func (ce *Cache) Get(key string) (val *Value) {
	ce.quMu.RLock()
	if im, ok := ce.qu[key]; ok {
		ce.quMu.RUnlock()
		val = im.Val
		return
	}
	ce.quMu.RUnlock()
	ce.trMu.RLock()
	r := ce.tr.Get(item{Key: key})
	if r == nil {
		ce.trMu.RUnlock()
		return
	}
	ce.trMu.RUnlock()
	val = r.(item).Val
	return
}

// Set sets the value of given key. It deletes the key, if the val is nil.
func (ce *Cache) Set(key string, val *Value) {
	ce.quMu.Lock()
	ce.qu[key] = item{Key: key, Val: val}
	ce.quMu.Unlock()
	select {
	case ce.quCh <- true:
	default:
	}
}

// Del deletes a key.
func (ce *Cache) Del(key string) {
	ce.Set(key, nil)
}

// GetOrSet returns the existing value for the key if present. Otherwise, it sets and returns the given value.
// The found is true if the key was exist, false if set.
func (ce *Cache) GetOrSet(key string, setval *Value) (val *Value, found bool) {
	found = true
	ce.quMu.Lock()
	if im, ok := ce.qu[key]; ok {
		ce.quMu.Unlock()
		val = im.Val
		return
	}
	ce.trMu.RLock()
	r := ce.tr.Get(item{Key: key})
	if r == nil {
		ce.qu[key] = item{Key: key, Val: setval}
		ce.quMu.Unlock()
		ce.trMu.RUnlock()
		select {
		case ce.quCh <- true:
		default:
		}
		val = setval
		found = false
		return
	}
	ce.quMu.Unlock()
	ce.trMu.RUnlock()
	val = r.(item).Val
	return
}
