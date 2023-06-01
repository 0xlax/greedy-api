package main

import (
	"fmt"
	"strings"
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

// Get method to retrieve the value for a given key
func (store *KeyValueStore) Get(key string) (string, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	if keyValue, ok := store.data[key]; ok {
		if keyValue.expiryTime != nil && time.Now().After(*keyValue.expiryTime) {
			delete(store.data, key)
			return "", fmt.Errorf("key has expired")
		}

		return keyValue.value, nil
	}

	return "", fmt.Errorf("key not found")
}

// QPush method to append values to a queue
func (store *KeyValueStore) QPush(key string, values ...string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, ok := store.data[key]; !ok {
		store.data[key] = &KeyValue{
			value:      "",
			expiryTime: nil,
		}
	}

	for _, value := range values {
		store.data[key].value += value + " "
	}
}

// QPop method to retrieve the last inserted value from a queue
func (store *KeyValueStore) QPop(key string) (string, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if keyValue, ok := store.data[key]; ok {
		if keyValue.expiryTime != nil && time.Now().After(*keyValue.expiryTime) {
			delete(store.data, key)
			return "", fmt.Errorf("queue is empty")
		}

		values := strings.Fields(keyValue.value)
		if len(values) > 0 {
			lastValue := values[len(values)-1]
			values = values[:len(values)-1]
			store.data[key].value = strings.Join(values, " ")
			return lastValue, nil
		}

		delete(store.data, key)
		return "", fmt.Errorf("queue is empty")
	}

	return "", fmt.Errorf("key not found")
}
