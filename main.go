package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
	store := &KeyValueStore{
		data:  make(map[string]*KeyValue),
		mutex: sync.RWMutex{},
	}

	http.HandleFunc("/", handleRequest(store))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
