package main

import (
	"fmt"
	"time"
)

// SET
// Writes the value to the datastore using the key and according to the specified parameters.
// Pattern: SET <key> <value> <expiry time>? <condition>?

func (store *KeyValueStore) Set(key string, value string, expiryTime time.Time, condition string) (string, error) {
	// one goroutine can access the critical section of code protected by the mutex at a time
	store.mutex.Lock()
	defer store.mutex.Unlock()

	// Check the condition if the key already exists
	if condition == "NX" {
		if _, ok := store.data[key]; ok {
			return "", fmt.Errorf("key already exists")
		}
	} else if condition == "XX" {
		if _, ok := store.data[key]; !ok {
			return "", fmt.Errorf("key does not exist")
		}
	}

	store.data[key] = &KeyValue{
		value:      value,
		expiryTime: &expiryTime,
	}

	return "OK", nil
}
