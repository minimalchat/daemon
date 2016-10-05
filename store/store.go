// Package store provides an in-memory key-value store.
package store

import (
  // "log"
  // "errors"
  // "math/rand"
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

// var errUnexpected = errors.New("an unexpected error")

// Put stores a value
func (db *InMemory) Put(v StoreKeyer) error {
  db.Lock()
  defer db.Unlock()

  if db.data == nil {
    db.data = make(map[string]StoreKeyer)
  }

  // if rand.Intn(10) < 5 {
  //   return errUnexpected
  // }

  db.data[v.StoreKey()] = v
  return nil
}

// Get retrieves a value
func (db *InMemory) Get(k string) (StoreKeyer, error) {
  db.RLock()
  defer db.RUnlock()

  // if rand.Intn(10) < 5 {
  //   return nil, errUnexpected
  // }

  return db.data[k], nil
}

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
