package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// KeyValue represents a key-value pair in the datastore.
// It stores the value and an optional expiry time for the key.
type KeyValue struct {
	Value      []string   // The value associated with the key
	ExpiryTime *time.Time // The expiry time for the key (optional)
}

// KeyValueStore represents an in-memory key-value data store.
// It stores the data and provides thread-safe access using a mutex.
type KeyValueStore struct {
	Data  map[string]*KeyValue // The underlying data store
	mutex sync.RWMutex         // Mutex for thread-safe access to the data store
}

// Mutex : Primitive used in concurrent programming to protect shared resources
// from being accessed simultaneously by multiple threads or goroutines

type Command struct {
	Command string `json:"command"` // Represents a JSON command received via the REST API.
}

type ErrorResponse struct {
	Error string `json:"error"` // Represents a JSON response containing an error message.
}

type ValueResponse struct {
	Value string `json:"value"` // Represents a JSON response containing a value.
}

var store = &KeyValueStore{
	Data: make(map[string]*KeyValue), // Initializes the key-value data store.
}

func main() {
	http.HandleFunc("/", handleRequest) // Sets up the request handler
	http.ListenAndServe(":8080", nil)   // Starts the HTTP server and listens on port 8080.
}

// Sends error response to the client.
func sendErrorResponse(w http.ResponseWriter, errorMessage string) {
	// Create ErrorResponse object as JSON with the specified error message.
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errorMessage})
}

// Sends a value response.
func sendValueResponse(w http.ResponseWriter, value string) {
	// CreateValueResponse object as JSON with the specified value.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ValueResponse{Value: value})
}

// Sends a simple OK response to the client.
func sendOKResponse(w http.ResponseWriter) {
	// Send an empty response as JSON to indicate a successful response.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

// ResponseWrites helps to onstruct and send response back to client
// Request represents incoming HTTP requests recieved from client

func handleRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body) //Decoder to decode request body into "Command" struct
	defer r.Body.Close()               //Request body is closed after request is processed

	var cmd Command
	err := decoder.Decode(&cmd)
	if err != nil {
		sendErrorResponse(w, "invalid request")
		return
	}

	parts := strings.Split(cmd.Command, " ") //Splits the command string into parts
	if len(parts) == 0 {
		sendErrorResponse(w, "invalid command")
		return
	}
	//First index is converted to uppercase and performed a switch statement to trigger appropriate function.
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
		handleBQPOP(w, parts) //Optional
	default:
		sendErrorResponse(w, "invalid command")
	}
}

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
	store.mutex.Lock() //write lock

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
		Value:      []string{value},
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
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	if kv, ok := store.Data[key]; ok {
		value := strings.Join(kv.Value, " ") // Convert the []string to a string
		sendValueResponse(w, value)
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

	// Acquire a lock on the store to ensure safe access
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	// Check if the key already exists in the store
	if kv, ok := store.Data[key]; ok {
		// Append the new values to the existing value slice
		kv.Value = append(kv.Value, values...)
	} else {
		// If the key doesn't exist, create a new KeyValue entry with the values as the slice
		store.Data[key] = &KeyValue{
			Value: values,
		}
	}

	sendOKResponse(w)
}

// OPTIONAL

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
		values := kv.Value // Retrieve the values slice directly

		if len(values) > 0 {
			// Set value to the last index
			value := values[len(values)-1]

			// Removes last index in the slice
			values = values[:len(values)-1]

			// Update the value in the store
			store.Data[key].Value = values

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

	store.mutex.RLock()
	kv, ok := store.Data[key]
	store.mutex.RUnlock()

	if ok {
		if timeout == 0 {
			// A value of 0 immediately returns a value from the queue without blocking. same as QPOP
			values := kv.Value
			if len(values) > 0 {
				value := values[len(values)-1]
				values = values[:len(values)-1]
				store.Data[key].Value = values
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
				store.mutex.RLock()
				kv, ok = store.Data[key]
				store.mutex.RUnlock()
				if ok {
					values := kv.Value
					if len(values) > 0 {
						value := values[len(values)-1]
						values = values[:len(values)-1]
						store.Data[key].Value = values
						sendValueResponse(w, value)
						return
					}
				}
			}
		}
	}

	sendErrorResponse(w, "queue is empty")
}
