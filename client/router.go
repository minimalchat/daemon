package client

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
Routes defines the Client API routes  */
func Routes(router *httprouter.Router, ds *store.InMemory) {

	// Client
	router.GET("/api/clients", readClients(ds)) // Check
	router.GET("/api/client", readClients(ds))

	router.GET("/api/client/:id", readClients(ds)) // Check

	router.POST("/api/client", createOrUpdateClient(ds)) // Not Implement

	router.POST("/api/client/", createOrUpdateClient(ds)) // Not Implement

	router.PUT("/api/client/:id", createOrUpdateClient(ds)) // Not Implement

	router.PATCH("/api/client/:id", createOrUpdateClient(ds)) // Not Implement

	router.DELETE("/api/client/:id", deleteClient(ds)) // Not Implement
}

/*
notImplemented is a helper function for intentionally unimplemented routes */
func notImplemented(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	resp.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	resp.WriteHeader(http.StatusNotImplemented)

	fmt.Fprintf(resp, "Not Implemented")
}

// Clients

/*
GET /api/client
GET /api/clients */
func readClients(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		clients, _ := ds.Search("client.")
		result := make(map[string]interface{})

		result["clients"] = clients

		log.Println(INFO, "api/client:", "Reading clients", fmt.Sprintf("(%d records)", len(clients)))

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		resp.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(resp).Encode(result); err != nil {
			panic(err)
		}
	}
}

// Client

/*
GET /api/client/:id */
func readClient(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		cl, _ := ds.Get(fmt.Sprintf("client.%s", params.ByName("id")))

		log.Println(DEBUG, "api/client:", "Reading client", params.ByName("id"))

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		resp.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(resp).Encode(cl); err != nil {
			panic(err)
		}
	}
}

/*
POST / PUT / PATCH /api/chat/:id/message/:mid */
func createOrUpdateClient(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return notImplemented
}

/*
DELETE /api/chat/:id/message/:mid */
func deleteClient(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return notImplemented
}
