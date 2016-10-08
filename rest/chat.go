package rest

import (
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/julienschmidt/httprouter"

  "github.com/mihok/letschat-daemon/store"
)


// Chats

// GET /api/chats
func ReadChats (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    chats, _ := db.Search("chat.")
    result := make(map[string]interface{})

    result["chats"] = chats;

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(resp).Encode(result); err != nil {
        panic(err)
    }
  }
}


// Chat

// GET /api/chat/:id
func ReadChat (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    ch, _ := db.Get(fmt.Sprintf("chat.%s", params.ByName("id")))

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if (ch != nil) {
      resp.WriteHeader(http.StatusOK)
      if err := json.NewEncoder(resp).Encode(ch); err != nil {
          panic(err)
      }
    } else {
      resp.WriteHeader(http.StatusNotFound)
    }
  }
}

// POST / PUT / PATCH /api/chat/:id?
func CreateOrUpdateChat (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return NotImplemented
}

// DELETE /api/chat/:id?
func DeleteChat (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return NotImplemented
}


// Chat Messages

// GET /api/chat/:id/messages
func ReadMessages (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
  }
}

// Chat Message

// GET /api/chat/:id/message/:mid
func ReadMessage (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return NotImplemented
}

// POST / PUT /api/chat/:id/message
func CreateMessage (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
  }
}

// PATCH /api/chat/:id/message/:mid
func UpdateMessage (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return NotImplemented
}

// DELETE /api/chat/:id/message/:mid
func DeleteMessage (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return NotImplemented
}