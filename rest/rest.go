package rest

import (
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/julienschmidt/httprouter"

  "github.com/mihok/letschat-daemon/store"
)

func Clients (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    clients, _ := db.Search("client.")
    result := make(map[string]interface{})

    result["clients"] = clients;

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(resp).Encode(result); err != nil {
        panic(err)
    }
  }
}

func Operators (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    operators, _ := db.Search("operator.")
    result := make(map[string]interface{})

    result["operators"] = operators;

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(resp).Encode(result); err != nil {
        panic(err)
    }
  }
}


// Default Route
func CatchAll (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
  resp.WriteHeader(http.StatusOK)
  fmt.Fprint(resp, "{\"clients\": \"/api/clients\", \"operators\":\"/api/operators\"}")
}