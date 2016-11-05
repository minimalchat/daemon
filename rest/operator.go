package rest

import (
  "log"
  "fmt"
  // "reflect"
  "net/http"
  "encoding/json"

  "github.com/julienschmidt/httprouter" // Router

  "github.com/minimalchat/mnml-daemon/operator"
  "github.com/minimalchat/mnml-daemon/store"
)


// Log levels
const (
  DEBUG string = "DEBUG"
  INFO string = "INFO"
  WARNING string = "WARN"
  ERROR string = "ERROR"
  FATAL string = "FATAL"
)


// Operators

// GET /api/operators
func ReadOperators (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    operators, _ := db.Search("operator.")
    result := make(map[string]interface{})

    result["operators"] = operators

    log.Println(DEBUG, "operator:", "Reading operators", fmt.Sprintf("(%d records)", len(operators)))

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(resp).Encode(result); err != nil {
        panic(err)
    }
  }
}

// Operator

// GET /api/operator/:id
func ReadOperator (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    id := params.ByName("id")
    op, _ := db.Get(fmt.Sprintf("operator.%s", id))

    log.Println(DEBUG, "operator:", "Reading operator", params.ByName("id"))

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    resp.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(resp).Encode(op); err != nil {
        panic(err)
    }
  }
}

// POST / PUT / PATCH /api/operator/:id?
func CreateOrUpdateOperator (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    var op *operator.Operator

    id := params.ByName("id")
    decoder := json.NewDecoder(req.Body)

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")

    if err := decoder.Decode(&op); err != nil {
        resp.WriteHeader(http.StatusBadRequest)

        fmt.Fprintf(resp, "Bad Request")
        return
    }

    if (id == "") { // Create

      // Save new record
      db.Put(op)
      log.Println(DEBUG, "operator:", "Creating new operator", op)

    } else { // Update

      //  Get old record
      result, _ := db.Get(fmt.Sprintf("operator.%s", id))

      if (result == nil) {
        resp.WriteHeader(http.StatusNotFound)

        fmt.Fprintf(resp, "Not Found")
        return
      }

      if old, ok := result.(operator.Operator); ok {

        // Update fields of old record
        if (op.FirstName != "") {
          old.FirstName = op.FirstName
        }

        if (op.LastName != "") {
          old.LastName = op.LastName
        }

        // Save old record
        db.Put(old)
        log.Println(DEBUG, "operator:", "Updating operator", old)
      }
    }

    resp.WriteHeader(http.StatusOK)
  }
}

// DELETE /api/operator/:id
func DeleteOperator (db *store.InMemory) func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
  return func (resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
    id := params.ByName("id")
    err := db.Remove(fmt.Sprintf("operator.%s", id))

    log.Println(DEBUG, "operator:", "Removing operator", id)

    resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err != nil {
        resp.WriteHeader(http.StatusBadRequest)

        fmt.Fprintf(resp, "Bad Request")
        return
    }

    resp.WriteHeader(http.StatusOK)
  }
}