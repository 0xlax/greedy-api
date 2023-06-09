package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleSET(t *testing.T) {
	// Create a new HTTP request for SET command
	body := strings.NewReader(`{"command": "SET key value"}`)
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handleRequest function with the request and recorder
	handleRequest(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// TODO: Add more assertions to test the behavior of the handleSET function
	// For example, you can check if the key-value pair is correctly stored in the data store.
}

func TestHandleGET(t *testing.T) {
	// Create a new HTTP request for GET command
	body := strings.NewReader(`{"command": "GET key"}`)
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handleRequest function with the request and recorder
	handleRequest(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// TODO: Add more assertions to test the behavior of the handleGET function
	// For example, you can check if the correct value is returned for the specified key.
}
