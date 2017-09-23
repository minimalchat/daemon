package socket

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io" // Socket

	// "github.com/minimalchat/daemon/chat"
	// "github.com/minimalchat/daemon/client"
	// "github.com/minimalchat/daemon/operator"
	// "github.com/minimalchat/daemon/person"
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
Server is the socket.io abstraction for Minimal Chat */
type Server struct {
	store *store.InMemory
	sock  *socketio.Server

	sockets map[*Socket]bool

	registerClient       chan *Socket
	registerOperator     chan *Socket
	unregister           chan *Socket
	broadcastToOperators chan *SocketMessage
	broadcastToClient    chan *SocketMessage
}

func Create(ds *store.InMemory) (*Server, error) {
	log.Println(DEBUG, "socket:", "Starting WebSocket server ...")

	ping, _ := time.ParseDuration("5s")

	srv := &Server{
		store: ds,

		registerClient:   make(chan *Socket),
		registerOperator: make(chan *Socket),

		unregister: make(chan *Socket),

		broadcastToOperators: make(chan *SocketMessage),
		broadcastToClient:    make(chan *SocketMessage),

		sockets: make(map[*Socket]bool),
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
			for sock := range s.sockets {
				if sock.conn.Id() == data.target {
					select {
					case sock.send <- data:
					default:
						close(sock.send)
						delete(s.sockets, sock)
					}
				}
			}
		case data := <-s.broadcastToOperators:
			for sock := range s.sockets {
				if sock.connType == OPERATOR {
					select {
					case sock.send <- data:
					default:
						log.Println(DEBUG, "socket:", fmt.Sprintf("%s send channel not available, closing ..", sock.conn.Id()))
						close(sock.send)
						delete(s.sockets, sock)
					}
				}
			}
		case sock := <-s.registerOperator:
			s.sockets[sock] = true
		case sock := <-s.registerClient:
			s.sockets[sock] = true
		case sock := <-s.unregister:
			if _, ok := s.sockets[sock]; ok {
				delete(s.sockets, sock)
				close(sock.send)
			}
		}
	}
}

func (s *Server) onConnect(c socketio.Socket) {

	var t SocketType

	query := c.Request().URL.Query()
	connectionType := query.Get("type")
	sessionId := query.Get("sessionId")

	// Identify the connection type
	switch connectionType {
	case "client":
		t = CLIENT
		break
	case "operator":
		t = OPERATOR
		break
	default:
		log.Println(WARNING, "socket:", "Unknown connection type, dropping ...")
		return
	}

	log.Println(INFO, "socket:", fmt.Sprintf("Incoming %s connection %s", t, c.Id()))

	// Create a Socket Connection
	sock := Socket{
		server: s,

		conn:     c,
		connType: t,

		send: make(chan *SocketMessage),
	}

	// Start listening for channel messages
	go sock.Listen()

	// Register event types
	// TODO: Do I really need to listen for both on every socket?

	sock.conn.On("client:message", func(data string) {
		go sock.onClientMessage(data)
	})

	sock.conn.On("client:typing", func(data string) {
		go sock.onClientTyping(data)
	})

	sock.conn.On("operator:message", func(data string) {
		go sock.onOperatorMessage(data)
	})

	sock.conn.On("operator:typing", func(data string) {
		go sock.onOperatorTyping(data)
	})

	sock.conn.On("disconnection", func() {
		s.unregister <- &sock
	})

	// Register the new client, depending on connection type
	switch sock.connType {
	case OPERATOR:
		// Register the new Socket with the server as an Operator
		s.registerOperator <- &sock

		// TODO: This may not be the right name for this func now
		go sock.onOperatorConnection()

		break
	case CLIENT:
		// Register the new Socket with the server as a Client
		s.registerClient <- &sock

		// TODO: This may not be the right name for this func now
		go sock.onClientConnection(sessionId)

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
