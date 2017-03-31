// Package store provides an in-memory key-value store.
package store

import (
  "fmt"
  "log"
  "sync"
  "strings"
)

// Log levels
const (
  DEBUG string = "DEBUG"
  INFO string = "INFO"
  WARNING string = "WARN"
  ERROR string = "ERROR"
  FATAL string = "FATAL"
)

// StoreKeyer is an object that can be kept in an InMemory store
type StoreKeyer interface {
  StoreKey() string
}

// An InMemory store is thread-safe and handles any StoreKeyer
type InMemory struct {
  sync.RWMutex
  data map[string]StoreKeyer
}

// Put stores a value
func (db *InMemory) Put(v StoreKeyer) error {
  db.Lock()
  defer db.Unlock()

  if db.data == nil {
    db.data = make(map[string]StoreKeyer)
  }

  log.Println(DEBUG, "store:", fmt.Sprintf("Storing %s ...", v.StoreKey()))

  db.data[v.StoreKey()] = v
  return nil
}

// Get retrieves a value
func (db *InMemory) Get(k string) (StoreKeyer, error) {
  db.RLock()
  defer db.RUnlock()

  log.Println(DEBUG, "store:", fmt.Sprintf("Getting %s ...", k))

  return db.data[k], nil
}

// Remove a value
func (db *InMemory) Remove(k string) error {
  db.Lock()
  defer db.Unlock()

  if db.data == nil {
    db.data = make(map[string]StoreKeyer)
  }

  log.Println(DEBUG, "store:", fmt.Sprintf("Deleting %s ...", k))

  delete(db.data, k)
  return nil
}

// Search on keys
func (db *InMemory) Search(q string) ([]StoreKeyer, error) {
  var result []StoreKeyer

  log.Println(DEBUG, "store:", fmt.Sprintf("Searching for %s ...", q))

  for key, _:= range db.data {
    if (strings.Contains(key, q)) {
      value, _ := db.Get(key)
      result = append(result, value)
    }
  }

  return result, nil
}
