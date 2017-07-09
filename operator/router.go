package operator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter" // Router

	"github.com/minimalchat/daemon/store"
)

// Log levels
const (
	DEBUG   string = "DEBUG"
	INFO    string = "INFO"
	WARNING string = "WARN"
	ERROR   string = "ERROR"
	FATAL   string = "FATAL"
)

/*
Routes defines the Operator API CRUD uris */
func Routes(router *httprouter.Router, ds *store.InMemory) {

	// Operator
	router.GET("/api/operators", readOperators(ds)) // Check
	router.GET("/api/operator", readOperators(ds))

	router.GET("/api/operator/:id", readOperator(ds)) // Check

	router.POST("/api/operator", createOrUpdateOperator(ds)) // Check

	router.POST("/api/operator/", createOrUpdateOperator(ds)) // Check

	router.PUT("/api/operator/:id", createOrUpdateOperator(ds)) // Check

	router.PATCH("/api/operator/:id", createOrUpdateOperator(ds)) // Check

	router.DELETE("/api/operator/:id", deleteOperator(ds)) // Check
}

// Operators

/*
GET /api/operator
GET /api/operators */
func readOperators(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		operators, _ := ds.Search("operator.")
		result := make(map[string]interface{})

		result["operators"] = operators

		log.Println(INFO, "api/operator:", "Reading operators", fmt.Sprintf("(%d records)", len(operators)))

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		resp.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(resp).Encode(result); err != nil {
			panic(err)
		}
	}
}

// Operator

/*
GET /api/operator/:id */
func readOperator(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		op, _ := ds.Get(fmt.Sprintf("operator.%s", id))

		log.Println(DEBUG, "api/operator:", "Reading operator", params.ByName("id"))

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		resp.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(resp).Encode(op); err != nil {
			panic(err)
		}
	}
}

/*
POST / PUT / PATCH /api/operator/:id? */
func createOrUpdateOperator(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var op *Operator

		id := params.ByName("id")
		decoder := json.NewDecoder(req.Body)

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := decoder.Decode(&op); err != nil {
			resp.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(resp, "Bad Request")
			return
		}

		if id == "" { // Create

			// Save new record
			ds.Put(op)
			log.Println(DEBUG, "api/operator:", "Creating new operator", op)

		} else { // Update

			//  Get old record
			result, _ := ds.Get(fmt.Sprintf("operator.%s", id))

			if result == nil {
				resp.WriteHeader(http.StatusNotFound)

				fmt.Fprintf(resp, "Not Found")
				return
			}

			if old, ok := result.(Operator); ok {

				// Update fields of old record
				if op.GetFirstName() != "" {
					old.FirstName = op.FirstName
				}

				if op.GetLastName() != "" {
					old.LastName = op.LastName
				}

				// Save old record
				ds.Put(old)
				log.Println(DEBUG, "api/operator:", "Updating operator", old)
			}
		}

		resp.WriteHeader(http.StatusOK)
	}
}

/*
DELETE /api/operator/:id */
func deleteOperator(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		err := ds.Remove(fmt.Sprintf("operator.%s", id))

		log.Println(DEBUG, "api/operator:", "Removing operator", id)

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(resp, "Bad Request")
			return
		}

		resp.WriteHeader(http.StatusOK)
	}
}
