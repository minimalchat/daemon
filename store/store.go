// Package store provides an in-memory key-value store.
package store

import (
  "sync"
  "strings"
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

  db.data[v.StoreKey()] = v
  return nil
}

// Get retrieves a value
func (db *InMemory) Get(k string) (StoreKeyer, error) {
  db.RLock()
  defer db.RUnlock()

  return db.data[k], nil
}

// Remove a value
func (db *InMemory) Remove(k string) error {
  db.Lock()
  defer db.Unlock()

  if db.data == nil {
    db.data = make(map[string]StoreKeyer)
  }

  delete(db.data, k)
  return nil
}

// Search on keys
func (db *InMemory) Search(q string) ([]StoreKeyer, error) {
  var result []StoreKeyer

  for key, _:= range db.data {
    if (strings.Contains(key, q)) {
      value, _ := db.Get(key)
      result = append(result, value)
    }
  }

  return result, nil
}
