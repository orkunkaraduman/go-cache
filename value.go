package cache

// Value encloses needing value of key-value pair.
type Value struct {
	// Value of key-value pair.
	// It's not concurrency safe. Don't change it after set operation without synchronizing them.
	// You may use Value struct in "sync/atomic" package as type of V.
	V interface{}
}
