package socket

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	// TODO: Move away from this library
	"github.com/minimalchat/go-socket.io" // Socket
	// "github.com/googollee/go-socket.io" // Socket

	// "github.com/minimalchat/daemon/chat"
	// "github.com/minimalchat/daemon/client"
	// "github.com/minimalchat/daemon/operator"
	// "github.com/minimalchat/daemon/person"
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
Server is the socket.io abstraction for Minimal Chat */
type Server struct {
	ID    string
	store *store.InMemory
	sock  *socketio.Server

	sockets map[*Conn]bool

	registerClient       chan *Conn
	registerOperator     chan *Conn
	unregister           chan *Conn
	broadcastToOperators chan *Message
	broadcastToClient    chan *Message
}

/*
Create takes a data store and returns a new socket server */
func Create(ds *store.InMemory) (*Server, error) {
	log.Println(DEBUG, "socket:", "Starting WebSocket server ...")

	ping, _ := time.ParseDuration("5s")

	srv := &Server{
		store: ds,

		registerClient:   make(chan *Conn),
		registerOperator: make(chan *Conn),

		unregister: make(chan *Conn),

		broadcastToOperators: make(chan *Message),
		broadcastToClient:    make(chan *Message),

		sockets: make(map[*Conn]bool),
	}

	sock, err := socketio.NewServer(nil)

	if err != nil {
		return nil, err
	}

	srv.sock = sock
	srv.sock.SetPingInterval(ping)

	srv.sock.On("connection", func(s socketio.Socket) {
		go srv.onConnect(s)
	})

	return srv, nil
}

/*
Listen creates a new Server instance and begins listening for ws://
connections. */
func (s Server) Listen() {

	for {
		select {
		case data := <-s.broadcastToClient:
			for c := range s.sockets {
				if c.raw.Id() == data.target {
					select {
					case c.send <- data:
					default:
						close(c.send)
						delete(s.sockets, c)
					}
				}
			}
		case data := <-s.broadcastToOperators:
			for c := range s.sockets {
				if c.category == OPERATOR {
					select {
					case c.send <- data:
					default:
						log.Println(DEBUG, "socket:", fmt.Sprintf("%s send channel not available, closing ..", c.raw.Id()))
						close(c.send)
						delete(s.sockets, c)
					}
				}
			}
		case c := <-s.registerOperator:
			s.sockets[c] = true
		case c := <-s.registerClient:
			s.sockets[c] = true
		case c := <-s.unregister:
			if _, ok := s.sockets[c]; ok {
				delete(s.sockets, c)
				close(c.send)
			}
		}
	}
}

func (s *Server) onConnect(raw socketio.Socket) {

	var cat Category

	query := raw.Request().URL.Query()

	connType := query.Get("type")
	accessID := query.Get("accessId")
	accessToken := query.Get("accessToken")
	sessionID := query.Get("sessionId")

	// Identify the connection type
	switch connType {
	case "client":
		cat = CLIENT
		break
	case "operator":
		cat = OPERATOR
		break
	default:
		log.Println(WARNING, "socket:", "Unknown connection type, dropping ...")
		return
	}

	log.Println(INFO, "socket:", fmt.Sprintf("Incoming %s connection %s", cat, raw.Id()))

	// Create a Socket Connection
	conn := Conn{
		server: s,

		raw:      raw,
		category: cat,

		send: make(chan *Message),
	}

	// Start listening for channel messages
	go conn.Listen()

	// Register event types
	// TODO: Do I really need to listen for both on every socket?

	conn.raw.On("client:message", func(data string) {
		go conn.onClientMessage(data)
	})

	conn.raw.On("client:typing", func(data string) {
		go conn.onClientTyping(data)
	})

	conn.raw.On("operator:message", func(data string) {
		go conn.onOperatorMessage(data)
	})

	conn.raw.On("operator:typing", func(data string) {
		go conn.onOperatorTyping(data)
	})

	conn.raw.On("disconnection", func() {
		s.unregister <- &conn
	})

	// Register the new client, depending on connection type
	switch conn.category {
	case OPERATOR:
		// Register the new Socket with the server as an Operator
		s.registerOperator <- &conn

		// TODO: This may not be the right name for this func now
		go conn.onOperatorConnection(accessID, accessToken)

		break
	case CLIENT:
		// Register the new Socket with the server as a Client
		s.registerClient <- &conn

		// TODO: This may not be the right name for this func now
		go conn.onClientConnection(sessionID)

		break
	default:
		log.Println(ERROR, "socket:", "Unknown connection type specified")
	}
}

/*
ServeHTTP serves the socket.io client script */
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.sock.ServeHTTP(w, r)
}
