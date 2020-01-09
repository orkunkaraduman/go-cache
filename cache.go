/*
Package cache offers concurrency safe in-memory cache based on b-tree and hash-map indexing.
All methods of Cache struct are concurrency safe and operates cache atomically.
*/
package cache

import (
	"sync"

	"github.com/google/btree"
)

var (
	// DefaultDegree is default b-tree degree.
	DefaultDegree = 4
)

// Cache struct is concurrency safe in-memory cache based on b-tree and hash-map indexing.
// All methods of Cache struct are concurrency safe and operates cache atomically.
type Cache struct {
	done   chan struct{}
	tr     *btree.BTree
	trMu   sync.RWMutex
	qu     chan *item
	degree int
}

// NewCache returns a new Cache has default degree.
func NewCache() (ce *Cache) {
	return NewCacheDegree(DefaultDegree)
}

// NewCacheDegree returns a new Cache given degree.
func NewCacheDegree(degree int) (ce *Cache) {
	ce = &Cache{
		done:   make(chan struct{}),
		qu:     make(chan *item, 1024),
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
	ce.trMu.Unlock()
	//ce.qu = make(chan *item, 1024)
}

// Close closes the cache. It must be called if the cache will not use.
func (ce *Cache) Close() {
	ce.done <- struct{}{}
	close(ce.qu)
}

func (ce *Cache) queueWorker() {
	for im := range ce.qu {
		ce.trMu.Lock()
		if im.Val != nil {
			ce.tr.ReplaceOrInsert(*im)
		} else {
			ce.tr.Delete(*im)
		}
		ce.trMu.Unlock()
	}
}

// Get returns the value of given key. It returns nil, if the key wasn't exist.
func (ce *Cache) Get(key string) (val interface{}) {
	done := false
	for !done {
		select {
		case im := <-ce.qu:
			if im.Key == key {
				val = im.Val
				return
			}
			ce.qu <- im
		default:
			done = true
		}
	}
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
func (ce *Cache) Set(key string, val interface{}) {
	ce.qu <- &item{Key: key, Val: val}
}

// Del deletes the key.
func (ce *Cache) Del(key string) {
	ce.Set(key, nil)
}

/*// GetOrSet returns the existing value for the key if present. Otherwise, it sets and returns the given value.
// If the key was exist, the found is true.
func (ce *Cache) GetOrSet(key string, newVal interface{}) (oldVal interface{}, found bool) {
	found = true
	ce.quMu.Lock()
	if im, ok := ce.qu[key]; ok {
		ce.quMu.Unlock()
		oldVal = im.Val
		return
	}
	ce.trMu.RLock()
	r := ce.tr.Get(item{Key: key})
	if r == nil {
		ce.qu[key] = item{Key: key, Val: newVal}
		ce.quMu.Unlock()
		ce.trMu.RUnlock()
		select {
		case ce.quCh <- struct{}{}:
		default:
		}
		oldVal = newVal
		found = false
		return
	}
	ce.quMu.Unlock()
	ce.trMu.RUnlock()
	oldVal = r.(item).Val
	return
}

// GetAndSet returns the replaced value for the key if present. Otherwise, returns nil.
// Value replaces by f.
func (ce *Cache) GetAndSet(key string, f func(interface{}) interface{}) (newVal interface{}) {
	ce.quMu.Lock()
	if im, ok := ce.qu[key]; ok {
		newVal = f(im.Val)
		ce.qu[key] = item{Key: key, Val: newVal}
		ce.quMu.Unlock()
		select {
		case ce.quCh <- struct{}{}:
		default:
		}
		return
	}
	ce.trMu.RLock()
	r := ce.tr.Get(item{Key: key})
	if r == nil {
		ce.quMu.Unlock()
		ce.trMu.RUnlock()
		return
	}
	newVal = f(r.(item).Val)
	ce.qu[key] = item{Key: key, Val: newVal}
	ce.quMu.Unlock()
	ce.trMu.RUnlock()
	select {
	case ce.quCh <- struct{}{}:
	default:
	}
	return
}

// Inc increases and the value of given key if the value is int or int64, and after returns new value.
// Otherwise returns old value.
func (ce *Cache) Inc(key string, x int64) (val interface{}) {
	return ce.GetAndSet(key, func(val2 interface{}) interface{} {
		switch val2.(type) {
		case int:
			return val2.(int) + int(x)
		case int64:
			return val2.(int64) + int64(x)
		}
		return val2
	})
}

// Dec decreases and the value of given key if the value is int or int64, and after returns new value.
// Otherwise returns old value.
func (ce *Cache) Dec(key string, x int64) (val interface{}) {
	return ce.GetAndSet(key, func(val2 interface{}) interface{} {
		switch val2.(type) {
		case int:
			return val2.(int) - int(x)
		case int64:
			return val2.(int64) - int64(x)
		}
		return val2
	})
}*/
