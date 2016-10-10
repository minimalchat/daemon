package rest

import (
  "log"
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/julienschmidt/httprouter"

  "github.com/mihok/letschat-daemon/store"
  "github.com/mihok/letschat-daemon/chat"
)


// Chats

// GET /api/chats
func ReadChats (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    chats, _ := db.Search("chat.")
    result := make(map[string]interface{})

    result["chats"] = chats;

    log.Println(DEBUG, "chat:", "Reading chats", fmt.Sprintf("(%d records)", len(chats)))

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

    log.Println(DEBUG, "chat:", "Reading chat", params.ByName("id"))

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
    messages, _ := db.Search(fmt.Sprintf("message.%s-", params.ByName("id")))
    result := make(map[string]interface{})

    result["messages"] = messages;

    log.Println(DEBUG, "message:", "Reading messages", fmt.Sprintf("(%d records)", len(messages)))

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(resp).Encode(result); err != nil {
        panic(err)
    }
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
    var msg *chat.Message

    id := params.ByName("id")
    decoder := json.NewDecoder(req.Body)

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")

    if err := decoder.Decode(&msg); err != nil {
      log.Println(DEBUG, "message:", "Bad Request", err)
      resp.WriteHeader(http.StatusBadRequest)
      return
    }

    if (id == "") {
      log.Println(DEBUG, "message:", "Bad Request ID", id)
      resp.WriteHeader(http.StatusBadRequest)
      return
    }

    result, _ := db.Get(fmt.Sprintf("chat.%s", id))

    if (result == nil) {
      log.Println(DEBUG, "message:", "Unknown Chat ID", id, result)
      resp.WriteHeader(http.StatusNotFound)
      return
    }

    if ch, ok := result.(chat.Chat); ok {
      log.Println(DEBUG, "operator:", msg.Content)

      // Fix if missing in Message object
      if (msg.Chat == "") {
        msg.Chat = id
      }

      db.Put(msg)

      ch.Client.Socket.Emit("operator:message", msg.Content, nil)
    }

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