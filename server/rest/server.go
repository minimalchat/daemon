package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter" // Http router

	"github.com/minimalchat/daemon/chat"
	"github.com/minimalchat/daemon/client"
	"github.com/minimalchat/daemon/operator"
	"github.com/minimalchat/daemon/server/socket"
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
	Config ServerConfig
}

type ServerConfig struct {
	Protocol string
	Port     int
	Host     string

	SSLCertFile string
	SSLKeyFile  string
	SSLPort     int

	CORSOrigin  string
	CORSEnabled bool
}

/*
Listen starts listening on `port` and `host` */
func Initialize(ds *store.InMemory, config ServerConfig) *Server {
	s := Server{
		Router: httprouter.New(),
		Config: config,
	}

	if s.Config.CORSEnabled {
		log.Println(DEBUG, "server:", fmt.Sprintf("Setting CORS origin to %s", config.CORSOrigin))
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

	if err != nil {
		log.Fatal(err)
	}

	go sock.Listen()

	s.Router.HandlerFunc("GET", "/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		sock.ServeHTTP(w, r)
	})

	// Operators API
	operator.Routes(s.Router, ds)

	// Clients API
	client.Routes(s.Router, ds)

	// Chats API
	chat.Routes(s.Router, ds)

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
	fmt.Fprint(resp, "{\"clients\": \"/api/clients\", \"client\": \"/api/client/:id\", \"chats\":\"/api/chats\", \"chat\":\"/api/chat/:id\", \"messages\":\"/api/chat/:id/messages\", \"message\":\"/api/chat/:id/message/:mid\", \"operators\":\"/api/operators\", \"operators\":\"/api/operators\", \"operator\":\"/api/operator/:id\"}")
}
