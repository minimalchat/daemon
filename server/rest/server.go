package rest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter" // Http router

	"github.com/minimalchat/daemon/chat"
	"github.com/minimalchat/daemon/client"
	"github.com/minimalchat/daemon/operator"
	"github.com/minimalchat/daemon/store" // InMemory store
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
func Listen(ds *store.InMemory) *Server {
	s := Server{
		Router: httprouter.New(),
	}

	// 404
	s.Router.NotFound = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		resp.WriteHeader(http.StatusNotFound)

		fmt.Fprintf(resp, "Not Found")
	})

	// 405
	s.Router.HandleMethodNotAllowed = true
	s.Router.MethodNotAllowed = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		resp.WriteHeader(http.StatusMethodNotAllowed)

		fmt.Fprintf(resp, "Method Not Allowed")
	})

	s.Router.GET("/", defaultRedirectRoute)
	s.Router.GET("/api", defaultRedirectRoute)
	s.Router.GET("/api/", defaultRoute)

	// Operators
	operator.Routes(s.Router, ds)

	// Clients
	client.Routes(s.Router, ds)

	// Chats
	chat.Routes(s.Router, ds)

	return &s
}

// GET /
func defaultRedirectRoute(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// resp.Header().Set("Content-Type", "text/html; charset=UTF-8")
	// resp.WriteHeader(http.StatusMovedPermanently)
	http.Redirect(resp, req, "/api/", http.StatusMovedPermanently)
}

// GET /api/
func defaultRoute(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprint(resp, "{\"clients\": \"/api/clients\", \"client\": \"/api/client/:id\", \"chats\":\"/api/chats\", \"chat\":\"/api/chat/:id\", \"messages\":\"/api/chat/:id/messages\", \"message\":\"/api/chat/:id/message/:mid\", \"operators\":\"/api/operators\", \"operators\":\"/api/operators\", \"operator\":\"/api/operator/:id\"}")
}
