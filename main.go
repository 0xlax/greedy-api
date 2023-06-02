package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type KeyValue struct {
	Value      string
	ExpiryTime *time.Time
}

type KeyValueStore struct {
	Data  map[string]*KeyValue
	mutex sync.Mutex
}

type Command struct {
	Command string `json:"command"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValueResponse struct {
	Value string `json:"value"`
}

var store = &KeyValueStore{
	Data: make(map[string]*KeyValue),
}

func main() {
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var cmd Command
	err := decoder.Decode(&cmd)
	if err != nil {
		sendErrorResponse(w, "invalid request")
		return
	}

	parts := strings.Split(cmd.Command, " ")
	if len(parts) == 0 {
		sendErrorResponse(w, "invalid command")
		return
	}

	switch strings.ToUpper(parts[0]) {
	case "SET":
		handleSET(w, parts)
	case "GET":
		handleGET(w, parts)
	case "QPUSH":
		handleQPUSH(w, parts)
	case "QPOP":
		handleQPOP(w, parts)
	case "BQPOP":
		handleBQPOP(w, parts)
	default:
		sendErrorResponse(w, "unknown command")
	}
}

// SET
// Writes the value to the datastore using the key and according to the specified parameters.
// Pattern: SET <key> <value> <expiry time>? <condition>?

func SET(input string, store *KeyValueStore) (string, error) {
	parts := strings.Split(input, " ")

	// Validate the input format
	if len(parts) < 3 || len(parts) > 5 {
		return "", fmt.Errorf("invalid input format")
	}

	key := parts[1]
	value := parts[2]
	var expiryTime time.Time
	var condition string

	// Parse expiry time if provided
	if len(parts) >= 4 {
		if strings.HasPrefix(parts[3], "EX") {
			seconds, err := strconv.Atoi(parts[3][2:])
			if err != nil {
				return "", fmt.Errorf("invalid expiry time")
			}
			expiryTime = time.Now().Add(time.Duration(seconds) * time.Second)
		} else {
			return "", fmt.Errorf("invalid input format")
		}
	}

	// Parse condition if provided
	if len(parts) == 5 {
		condition = parts[4]
		if condition != "NX" && condition != "XX" {
			return "", fmt.Errorf("invalid condition")
		}
	}

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
			if keyValue.expiryTime.IsZero() {
				delete(store.data, key)
				return "", fmt.Errorf("key has expired")
			}
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

// BQPop method to perform a blocking queue read operation
func (store *KeyValueStore) BQPop(key string, timeout float64) (string, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	// Check if the queue exists
	if _, ok := store.data[key]; !ok {
		return "", fmt.Errorf("queue does not exist")
	}

	queue := store.data[key]

	// Check if the queue is empty
	if len(strings.Fields(queue.value)) > 0 {
		return store.QPop(key)
	}

	// If timeout is 0, return immediately without blocking
	if timeout == 0 {
		return "", nil
	}

	// Create a channel to wait for the value to be pushed
	valueChan := make(chan string, 1)

	// Start a goroutine to read from the channel after the specified timeout
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		valueChan <- ""
	}()

	select {
	case value := <-valueChan:
		if value == "" {
			return "", fmt.Errorf("timeout exceeded")
		}
		return value, nil
	}
}
