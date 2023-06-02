package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func handleSET(w http.ResponseWriter, parts []string) {
	if len(parts) < 3 {
		sendErrorResponse(w, "invalid command format")
		return
	}

	key := parts[1]   //sets key
	value := parts[2] // sets value

	//Currently - empty initialization
	var expiryTime time.Time
	var condition string

	if len(parts) >= 4 && strings.HasPrefix(parts[3], "EX") {
		// extracts the number of seconds for the expiry time, converts it to an integer
		// sets the expiryTime variable to the current time plus the specified duration.
		seconds, err := strconv.Atoi(parts[3][2:])
		if err != nil {
			sendErrorResponse(w, "invalid expiry time")
			return
		}
		expiryTime = time.Now().Add(time.Duration(seconds) * time.Second)
	}

	if len(parts) == 5 {
		condition = strings.ToUpper(parts[4])
		if condition != "NX" && condition != "XX" {
			sendErrorResponse(w, "invalid condition")
			return
		}
	}
	//Makes sure only one process can use the store at one time
	// To Support COncurrent Operations
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if condition == "NX" {
		if _, ok := store.Data[key]; ok {
			sendErrorResponse(w, "key already exists")
			return
		}
	} else if condition == "XX" {
		if _, ok := store.Data[key]; !ok {
			sendErrorResponse(w, "key does not exist")
			return
		}
	}

	store.Data[key] = &KeyValue{
		Value:      value,
		ExpiryTime: &expiryTime,
	}

	sendOKResponse(w)
}

// retrieves the value associated with a given key from the data store, ensuring concurrent access using a mutex lock.
func handleGET(w http.ResponseWriter, parts []string) {
	if len(parts) != 2 {
		sendErrorResponse(w, "invalid command format")
		return
	}

	key := parts[1]

	//Makes sure only one process can use the store at one time
	// To Support Concurrent Operations
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if kv, ok := store.Data[key]; ok {
		sendValueResponse(w, kv.Value)
		return
	}

	sendErrorResponse(w, "key not found")
}

func handleQPUSH(w http.ResponseWriter, parts []string) {
	if len(parts) < 3 {
		sendErrorResponse(w, "invalid command format")
		return
	}

	key := parts[1]
	values := parts[2:]
	//acquires a lock on the key-value store to ensure that only one process can access it at a time.
	store.mutex.Lock()
	defer store.mutex.Unlock()

	//checks if the key already exists in the store.
	// If it does, it appends the new values to the existing value by concatenating
	if kv, ok := store.Data[key]; ok {
		for _, value := range values {
			kv.Value += " " + value
		}
		// If the key doesn't exist in the store, it creates a new KeyValue entry with the joined values as the value.
	} else {
		store.Data[key] = &KeyValue{
			Value: strings.Join(values, " "),
		}
	}

	sendOKResponse(w)
}

func handleQPOP(w http.ResponseWriter, parts []string) {
	if len(parts) != 2 {
		sendErrorResponse(w, "invalid command format")
		return
	}

	key := parts[1]

	// Acquire a lock on the store to ensure safe access
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if kv, ok := store.Data[key]; ok {
		values := strings.Split(kv.Value, " ")
		if len(values) > 0 {

			// Set value to the last index
			value := values[len(values)-1]

			// Removes last index in the array
			values = values[:len(values)-1]
			kv.Value = strings.Join(values, " ")

			// Send the last value as the response
			sendValueResponse(w, value)
			return
		}
	}

	sendErrorResponse(w, "queue is empty")
}

// OPTIONAL HANDLER FUNCTION

// handleBQPOP handles the blocking queue behavior by allowing
// the caller to wait for a certain period for a value to be available in the queue
// or to immediately retrieve a value if the queue is non-empty.

func handleBQPOP(w http.ResponseWriter, parts []string) {
	if len(parts) != 3 {
		sendErrorResponse(w, "invalid command format")
		return
	}

	key := parts[1]
	timeout, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		sendErrorResponse(w, "invalid timeout")
		return
	}

	store.mutex.Lock()
	kv, ok := store.Data[key]
	store.mutex.Unlock()

	if ok {
		if timeout == 0 {
			// A value of 0 immediately returns a value from the queue without blocking. same as QPOP
			values := strings.Split(kv.Value, " ")
			if len(values) > 0 {
				value := values[len(values)-1]
				values = values[:len(values)-1]
				kv.Value = strings.Join(values, " ")
				sendValueResponse(w, value)
				return
			}
		} else if timeout > 0 {
			// convert the timeout value from seconds (represented as a float64) to a time.Duration value.
			ticker := time.NewTicker(time.Duration(timeout) * time.Second)
			select {

			// If the ticker emitted a value, it means the specified timeout duration has elapsed.
			// Send a timeout error response and return.
			case <-ticker.C:
				sendErrorResponse(w, "timeout")
				return
			// If the ticker didn't emit a value before the timeout
			case <-time.After(1 * time.Second):
				store.mutex.Lock()
				kv, ok = store.Data[key]
				store.mutex.Unlock()
				if ok {
					values := strings.Split(kv.Value, " ")
					if len(values) > 0 {
						value := values[len(values)-1]
						values = values[:len(values)-1]
						kv.Value = strings.Join(values, " ")
						sendValueResponse(w, value)
						return
					}
				}
			}
		}
	}

	sendErrorResponse(w, "queue is empty")
}
