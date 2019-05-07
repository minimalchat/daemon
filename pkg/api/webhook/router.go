package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/julienschmidt/httprouter" // Router

	"github.com/minimalchat/daemon/pkg/store"
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
Routes defines the Webhook API CRUD uris */
func Routes(router *httprouter.Router, ds *store.InMemory) {

	// Operator
	// Read
	// router.GET("/api/operators", readOperators(ds)) // Check
	// router.GET("/api/operator", readOperators(ds))
	router.GET("/api/webhooks", readWebhooks(ds))
	router.GET("/api/webhook", readWebhooks(ds))

	router.GET("/api/webhook/:id", readWebhooks(ds))

	// router.GET("/api/operator/:id", readOperator(ds)) // Check

	// Create / Update
	// router.POST("/api/operator", createOrUpdateOperator(ds))      // Check
	// router.POST("/api/operator/:id", createOrUpdateOperator(ds))  // Check
	// router.PUT("/api/operator/", createOrUpdateOperator(ds))      // Check
	// router.PATCH("/api/operator/:id", createOrUpdateOperator(ds)) // Check
	router.POST("/api/webhook", createOrUpdateWebhook(ds))
	router.POST("/api/webhook/:id", createOrUpdateWebhook(ds))
	router.PUT("/api/webhook/", createOrUpdateWebhook(ds))
	router.PATCH("/api/webhook/:id", createOrUpdateWebhook(ds))

	// Delete
	// router.DELETE("/api/operator/:id", deleteOperator(ds)) // Check
	router.DELETE("/api/webhook/:id", deleteWebhook(ds))
}

/*
GET /api/webhook
GET /api/webhooks
GET /api/webhook/:id */
func readWebhooks(ds *store.InMemory) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		id := params.ByName("id")

		if id == "" {
			ws, _ := ds.Search("webhook.")
			result := make(map[string]interface{})

			result["webhooks"] = ws

			log.Println(INFO, "api/webhook:", "Reading webhooks", fmt.Sprintf("(%d records)", len(ws)))

			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			resp.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(resp).Encode(result); err != nil {
				panic(err)
			}
		} else {
			w, _ := ds.Get(fmt.Sprintf("webhook.%s", id))

			log.Println(DEBUG, "api/webhook:", "Reading webhook", id)

			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if w != nil {
				resp.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(resp).Encode(w); err != nil {
					panic(err)
				}
			} else {
				resp.WriteHeader(http.StatusNotFound)

				fmt.Fprintf(resp, "Not Found")
			}
		}
	}
}

/*
POST /api/webhook/
POST /api/webhook/:id
PUT /api/webhook/
PATCH /api/webhook/:id */
func createOrUpdateWebhook(ds *store.InMemory) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var w *Webhook

		id := params.ByName("id")
		decoder := json.NewDecoder(req.Body)

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := decoder.Decode(&w); err != nil {
			log.Println(ERROR, "api/webhook:", err)

			resp.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(resp, "Bad Request")
			return
		}

		if id == "" { // Create

			log.Println(DEBUG, "api/webhook:", "Creating new webhook", w)

			// Use our creation function to generate an ID and secret
			wh := CreateWebhook(w.Uri, w.EventTypes)

			if w.GetSecret() != "" {
				wh.Secret = w.Secret
			}

			// Save new record
			ds.Put(wh)

		} else { // Update

			//  Get old record
			result, _ := ds.Get(fmt.Sprintf("webhook.%s", id))

			if result == nil {
				resp.WriteHeader(http.StatusNotFound)

				fmt.Fprintf(resp, "Not Found")
				return
			}

			if old, ok := result.(*Webhook); ok {
				// Update fields of old record

				if w.GetUri() != "" {
					old.Uri = w.Uri
				}

				now := time.Now()
				seconds := now.Unix()
				nanos := int32(now.Sub(time.Unix(seconds, 0)))

				old.UpdatedTime = &timestamp.Timestamp{
					Seconds: seconds,
					Nanos:   nanos,
				}

				log.Println(DEBUG, "api/webhook:", "Updating webhook", old)

				// Save old record
				ds.Put(old)
			}
		}

		resp.WriteHeader(http.StatusOK)
	}
}

/*
DELETE /api/webhook/:id */
func deleteWebhook(ds *store.InMemory) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		id := params.ByName("id")
		err := ds.Remove(fmt.Sprintf("webhook.%s", id))

		log.Println(DEBUG, "api/operator:", "Removing operator", id)

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			log.Println(ERROR, "api/webhook:", err)

			resp.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(resp, "Bad Request")
			return
		}

		resp.WriteHeader(http.StatusOK)
	}
}
