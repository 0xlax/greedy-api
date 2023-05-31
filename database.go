package main

import (
	"sync"
	"time"
)

type DataItem struct {
	Key    string      `json:"key"`
	Value  interface{} `json:"value"` //empty interface that represents a type that can hold values of any type.
	Expiry time.Time   `json:"time"`
	// Other fields as needed for conditions
}

type DataStore struct {
	data  map[string]DataItem
	mutex sync.RWMutex //For concurrent safe access

}
