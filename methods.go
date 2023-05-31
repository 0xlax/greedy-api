package main

import "time"

// SET
// Writes the value to the datastore using the key and according to the specified parameters.
// Pattern: SET <key> <value> <expiry time>? <condition>?

func (store *KeyValueStore) Set(key string, value string, expiryTime time.Time, condition string) (string, error) {
	// one goroutine can access the critical section of code protected by the mutex at a time
	store.mutex.Lock()
	defer store.mutex.Unlock()
}
