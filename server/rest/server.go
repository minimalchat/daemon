package rest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter" // Http router

	"github.com/minimalchat/mnml-daemon/chat"
	"github.com/minimalchat/mnml-daemon/client"
	"github.com/minimalchat/mnml-daemon/operator"
	"github.com/minimalchat/mnml-daemon/store" // InMemory store
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
Server is the REST API server for Minimal Chat */
type Server struct {
	Router *httprouter.Router
}

/*
Listen starts listening on `port` and `host` */
func Listen(host string, port int, ds *store.InMemory) *Server {
	srv := Server{
		Router: httprouter.New(),
	}

	// 404
	srv.Router.NotFound = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		resp.WriteHeader(http.StatusNotFound)

		fmt.Fprintf(resp, "Not Found")
	})

	// 405
	srv.Router.HandleMethodNotAllowed = true
	srv.Router.MethodNotAllowed = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		resp.WriteHeader(http.StatusMethodNotAllowed)

		fmt.Fprintf(resp, "Method Not Allowed")
	})

	srv.Router.GET("/api", defaultRoute)
	srv.Router.GET("/api/", defaultRoute)

	// Operators
	operator.Routes(srv.Router, ds)

	// Clients
	client.Routes(srv.Router, ds)

	// Chats
	chat.Routes(srv.Router, ds)

	return &srv
}

// GET /api/
func defaultRoute(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, "{\"clients\": \"/api/clients\", \"client\": \"/api/client/:id\", \"chats\":\"/api/chats\", \"chat\":\"/api/chat/:id\", \"messages\":\"/api/chat/:id/messages\", \"message\":\"/api/chat/:id/message/:mid\", \"operators\":\"/api/operators\", \"operators\":\"/api/operators\", \"operator\":\"/api/operator/:id\"}")
}
