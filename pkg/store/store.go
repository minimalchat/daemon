// Package store provides an in-memory key-value store.
package store

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

/*
Log Levels */
const (
	DEBUG   string = "DEBUG"
	INFO    string = "INFO"
	WARNING string = "WARN"
	ERROR   string = "ERROR"
	FATAL   string = "FATAL"
)

/*
Keyer is an object that can be kept in an InMemory store */
type Keyer interface {
	Key() string
}

/*
InMemory store is thread-safe and handles any object implemented as a Keyer */
type InMemory struct {
	sync.RWMutex
	data map[string]Keyer
}

/*
Put stores a value */
func (db *InMemory) Put(v Keyer) error {
	db.Lock()
	defer db.Unlock()

	if db.data == nil {
		db.data = make(map[string]Keyer)
	}

	log.Println(DEBUG, "store:", fmt.Sprintf("Storing %s ...", v.Key()))

	db.data[v.Key()] = v
	return nil
}

/*
Get retrieves a value */
func (db *InMemory) Get(k string) (Keyer, error) {
	db.RLock()
	defer db.RUnlock()

	log.Println(DEBUG, "store:", fmt.Sprintf("Getting %s ...", k))

	return db.data[k], nil
}

/*
Remove a value */
func (db *InMemory) Remove(k string) error {
	db.Lock()
	defer db.Unlock()

	if db.data == nil {
		db.data = make(map[string]Keyer)
	}

	log.Println(DEBUG, "store:", fmt.Sprintf("Deleting %s ...", k))

	delete(db.data, k)
	return nil
}

/*
Search on keys */
func (db *InMemory) Search(q string) ([]Keyer, error) {
	var result []Keyer

	log.Println(DEBUG, "store:", fmt.Sprintf("Searching for %s ...", q))

	for key := range db.data {
		if strings.Contains(key, q) {
			value, _ := db.Get(key)
			result = append(result, value)
		}
	}

	return result, nil
}
