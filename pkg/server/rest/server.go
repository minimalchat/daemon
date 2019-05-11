package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter" // Http router

	"github.com/minimalchat/daemon/pkg/api/chat"
	"github.com/minimalchat/daemon/pkg/api/client"
	"github.com/minimalchat/daemon/pkg/api/operator"
	"github.com/minimalchat/daemon/pkg/api/webhook"
	"github.com/minimalchat/daemon/pkg/server/socket"
	"github.com/minimalchat/daemon/pkg/store" // InMemory store
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
	Config
}

/*
Config holds all the necessary configuration for our REST API server */
type Config struct {
	Protocol string
	Port     string
	Host     string

	Id string

	SSLCertFile string
	SSLKeyFile  string
	SSLPort     int

	CORSOrigin  string
	CORSEnabled bool
}

/*
Initialize takes a Store and ServerConfig starts listening on port and host
provided by a ServerConfig */
func Initialize(ds *store.InMemory, c Config) *Server {
	s := Server{
		Router: httprouter.New(),
		Config: c,
	}

	if s.Config.CORSEnabled {
		log.Println(DEBUG, "server:", fmt.Sprintf("Setting CORS origin to %s", c.CORSOrigin))
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

	// Default Routes
	s.Router.GET("/", defaultRedirectRoute)
	s.Router.GET("/api", defaultRedirectRoute)
	s.Router.GET("/api/", defaultRoute)

	// Socket.io
	sock, err := socket.Create(ds)
	sock.Id = c.Id

	if err != nil {
		log.Fatal(err)
	}

	go sock.Listen()

	s.Router.HandlerFunc("GET", "/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		if s.Config.CORSEnabled {
			w.Header().Set("Access-Control-Allow-Origin", s.Config.CORSOrigin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			// resp.Header().Set("Access-Control-Allow-Headers", "X-Socket-Type")
		}

		sock.ServeHTTP(w, r)
	})

	// Operators API
	operator.Routes(s.Router, ds)

	// Clients API
	client.Routes(s.Router, ds)

	// Chats API
	chat.Routes(s.Router, ds)

	// Webhook API
	webhook.Routes(s.Router, ds)

	return &s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Config.CORSEnabled {
		w.Header().Set("Access-Control-Allow-Origin", s.Config.CORSOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// resp.Header().Set("Access-Control-Allow-Headers", "X-Socket-Type")
	}

	s.Router.ServeHTTP(w, r)
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
	// TODO: Make this less hacky?
	fmt.Fprint(resp, "{\"clients\": \"/api/clients\", \"client\": \"/api/client/:id\", \"chats\":\"/api/chats\", \"chat\":\"/api/chat/:id\", \"messages\":\"/api/chat/:id/messages\", \"message\":\"/api/chat/:id/message/:mid\", \"operators\":\"/api/operators\", \"operators\":\"/api/operators\", \"operator\":\"/api/operator/:id\", \"webhooks\":\"/api/webhooks\", \"webhook\":\"/api/webhook/:id\"}")
}
