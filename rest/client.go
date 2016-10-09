package rest

import (
  "log"
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/julienschmidt/httprouter"

  "github.com/mihok/letschat-daemon/store"
)

// Clients

// GET /api/clients
func ReadClients (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    clients, _ := db.Search("client.")
    result := make(map[string]interface{})

    result["clients"] = clients;

    log.Println(DEBUG, "client:", "Reading clients", fmt.Sprintf("(%d records)", len(clients)))

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(resp).Encode(result); err != nil {
        panic(err)
    }
  }
}

// Client

// GET /api/client/:id
func ReadClient (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    cl, _ := db.Get(fmt.Sprintf("client.%s", params.ByName("id")))

    log.Println(DEBUG, "client:", "Reading client", params.ByName("id"))

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(resp).Encode(cl); err != nil {
        panic(err)
    }
  }
}

// POST / PUT / PATCH /api/chat/:id/message/:mid
func CreateOrUpdateClient (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return NotImplemented
}

// DELETE /api/chat/:id/message/:mid
func DeleteClient (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return NotImplemented
}