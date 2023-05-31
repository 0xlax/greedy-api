package main

import (
	"sync"
	"time"
)

type KeyValueStore struct {
	data  map[string]*KeyValue
	mutex sync.RWMutex
}

type KeyValue struct {
	value      string
	expiryTime *time.Time
}
