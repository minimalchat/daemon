package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
Routes defines the Chat and Chat Message API routes */
func Routes(router *httprouter.Router, ds *store.InMemory) {

	// Chat
	router.GET("/api/chats", readChats(ds)) // Check
	router.GET("/api/chat", readChats(ds))

	router.GET("/api/chat/:id", readChat(ds)) // Check

	router.POST("/api/chat", createOrUpdateChat(ds)) // Not Implement

	router.POST("/api/chat/", createOrUpdateChat(ds)) // Not Implement

	router.PUT("/api/chat/:id", createOrUpdateChat(ds)) // Not Implement

	router.PATCH("/api/chat/:id", createOrUpdateChat(ds)) // Not Implement

	router.DELETE("/api/chat/:id", deleteChat(ds)) // Not Implement

	// Chat Messages
	router.GET("/api/chat/:id/messages", readMessages(ds)) // Check
	router.GET("/api/chat/:id/message", readMessages(ds))

	router.GET("/api/chat/:id/message/:mid", readMessage(ds)) // Not Implement

	router.POST("/api/chat/:id/message", createMessage(ds)) // Check

	router.POST("/api/chat/:id/message/", createMessage(ds)) // Check

	router.PUT("/api/chat/:id/message/:mid", updateMessage(ds)) // Not Implement

	router.PATCH("/api/chat/:id/message/:mid", updateMessage(ds)) // Not Implement

	router.DELETE("/api/chat/:id/message/:mid", deleteMessage(ds)) // Not Implement

}

/*
notImplemented is a helper function for intentionally unimplemented routes */
func notImplemented(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	resp.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	resp.WriteHeader(http.StatusNotImplemented)

	fmt.Fprintf(resp, "Not Implemented")
}

// Chats

/*
GET /api/chat
GET /api/chats */
func readChats(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		chats, _ := ds.Search("chat.")
		result := make(map[string]interface{})

		result["chats"] = chats

		log.Println(INFO, "api/chat:", "Reading chats", fmt.Sprintf("(%d records)", len(chats)))

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		resp.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(resp).Encode(result); err != nil {
			panic(err)
		}
	}
}

// Chat

/*
GET /api/chat/:id */
func readChat(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ch, _ := ds.Get(fmt.Sprintf("chat.%s", params.ByName("id")))

		log.Println(DEBUG, "api/chat:", "Reading chat", params.ByName("id"))

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if ch != nil {
			resp.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(resp).Encode(ch); err != nil {
				panic(err)
			}
		} else {
			resp.WriteHeader(http.StatusNotFound)

			fmt.Fprintf(resp, "Not Found")
		}
	}
}

/*
POST / PUT / PATCH /api/chat/:id? */
func createOrUpdateChat(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return notImplemented
}

/*
DELETE /api/chat/:id? */
func deleteChat(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return notImplemented
}

// Chat Messages

/*
GET /api/chat/:id/message
GET /api/chat/:id/messages */
func readMessages(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		messages, _ := ds.Search(fmt.Sprintf("message.%s-", params.ByName("id")))
		result := make(map[string]interface{})

		result["messages"] = messages

		log.Println(INFO, "api/message:", "Reading messages", fmt.Sprintf("(%d records)", len(messages)))

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		resp.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(resp).Encode(result); err != nil {
			panic(err)
		}
	}
}

// Chat Message

/*
GET /api/chat/:id/message/:mid */
func readMessage(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return notImplemented
}

/*
POST / PUT /api/chat/:id/message */
func createMessage(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		var msg *Message

		id := params.ByName("id")
		decoder := json.NewDecoder(req.Body)

		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := decoder.Decode(&msg); err != nil {
			log.Println(ERROR, "api/message:", "Bad Request", err)
			resp.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(resp, "Bad Request")
			return
		}

		if id == "" {
			log.Println(ERROR, "api/message:", "Bad Request ID", id)
			resp.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(resp, "Bad Request")
			return
		}

		result, _ := ds.Get(fmt.Sprintf("chat.%s", id))

		if result == nil {
			log.Println(DEBUG, "api/message:", "Unknown Chat ID", id, result)
			resp.WriteHeader(http.StatusNotFound)

			fmt.Fprintf(resp, "Not Found")
			return
		}

		if ch, ok := result.(Chat); ok {
			log.Println(DEBUG, "api/operator:", msg.Content, ch.Uid)

			// Fix if missing in Message object
			if msg.Chat == "" {
				msg.Chat = id
			}

			ds.Put(msg)

			// TODO: Pass API post message to Operator via socket somehow
			// ch.Client.Socket.Emit("operator:message", msg.Content, nil)
		} else {
			log.Println(ERROR, "api/message:", "Could not cast store data to struct", ok, result.(Chat))
			resp.WriteHeader(http.StatusInternalServerError)

			fmt.Fprintf(resp, "Bad Request")
			return
		}

		resp.WriteHeader(http.StatusOK)
	}
}

/*
PATCH /api/chat/:id/message/:mid */
func updateMessage(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return notImplemented
}

/*
DELETE /api/chat/:id/message/:mid */
func deleteMessage(ds *store.InMemory) func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	return notImplemented
}
