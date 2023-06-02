package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", handleRequest) // Sets up the request handler
	http.ListenAndServe(":8080", nil)   // Starts the HTTP server and listens on port 8080.
}

// ResponseWrites helps to onstruct and send response back to client
// Request represents incoming HTTP requests recieved from client

func handleRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body) //Decoder to decode request body into "Command" struct
	defer r.Body.Close()               //Request body is closed after erquest is processed

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
