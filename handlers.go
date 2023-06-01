package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func handleRequest(store *KeyValueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var request map[string]string
		err := decoder.Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		command, ok := request["command"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		parts := strings.Fields(command)
		if len(parts) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response := make(map[string]interface{})

		switch parts[0] {
		case "SET":
			if len(parts) < 3 {
				response["error"] = "invalid command"
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			key := parts[1]
			value := parts[2]
			expiryTime := time.Time{}
			condition := ""

			for i := 3; i < len(parts); i++ {
				if parts[i] == "EX" && i+1 < len(parts) {
					expirySeconds, err := strconv.Atoi(parts[i+1])
					if err != nil {
						response["error"] = "invalid expiry time"
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					expiryTime = time.Now().Add(time.Duration(expirySeconds) * time.Second)
					i++
				} else if parts[i] == "NX" || parts[i] == "XX" {
					condition = parts[i]
				} else {
					response["error"] = "invalid command"
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			result, err := store.Set(key, value, expiryTime, condition)
			if err != nil {
				response["error"] = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			response["result"] = result
			w.WriteHeader(http.StatusOK)
		case "GET":
			if len(parts) != 2 {
				response["error"] = "invalid command"
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			key := parts[1]
			value, err := store.Get(key)
			if err != nil {
				response["error"] = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			response["value"] = value
			w.WriteHeader(http.StatusOK)
		case "QPUSH":
			if len(parts) < 3 {
				response["error"] = "invalid command"
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			key := parts[1]
			values := parts[2:]
			store.QPush(key, values...)
			w.WriteHeader(http.StatusOK)
		case "QPOP":
			if len(parts) != 2 {
				response["error"] = "invalid command"
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			key := parts[1]
			value, err := store.QPop(key)
			if err != nil {
				response["error"] = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			response["value"] = value
			w.WriteHeader(http.StatusOK)
		case "BQPOP":
			if len(parts) != 3 {
				response["error"] = "invalid command"
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			key := parts[1]
			timeout, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				response["error"] = "invalid timeout"
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			value, err := store.BQPop(key, timeout)
			if err != nil {
				response["error"] = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				break
			}

			response["value"] = value
			w.WriteHeader(http.StatusOK)
		default:
			response["error"] = "unknown command"
			w.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(response)
	}
}
