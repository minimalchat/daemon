package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io" // Socket

	"github.com/minimalchat/daemon/chat"
	"github.com/minimalchat/daemon/client"
	"github.com/minimalchat/daemon/operator"
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
	// TODO: Poor data structure, should thing of something smarter?
	//  We need to store the active sockets rather than specific application
	//  data.
	Sockets map[string]socketio.Socket
	store   *store.InMemory
	server  *socketio.Server
}

/*
Listen creates a new Server instance and begins listening for ws://
connections. */
func Listen(ds *store.InMemory) (*Server, error) {
	log.Println(DEBUG, "socket:", "Starting WebSocket server ...")

	ping, _ := time.ParseDuration("5s")

	srv, err := socketio.NewServer(nil)
	srv.SetPingInterval(ping)

	sck := Server{
		Sockets: make(map[string]socketio.Socket),
		store:   ds,
		server:  srv,
	}

	// TODO: Return an error instead
	if err != nil {
		return nil, err
	}

	srv.On("connection", sck.onConnect)
	srv.On("error", sck.onError)

	return &sck, nil
}

/*
ServeHTTP serves the socket.io client script */
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}

func (s Server) emitToOperators(event string, data string) {
	ops, _ := s.store.Search("operator.")

	if len(ops) == 0 {
		log.Println(INFO, "socket:", fmt.Sprintf("No operators online (%d records)", len(ops)))
		return
	}

	log.Println(INFO, "socket:", fmt.Sprintf("Sending '%s' message (%d operators)", event, len(ops)))

	for _, op := range ops {
		o := op.(*operator.Operator)
		osck, ok := s.Sockets[o.Uid]

		if !ok {
			log.Println(WARNING, "socket:", "Operator went away")
			continue
		}

		osck.Emit(event, data, nil)

		log.Println(DEBUG, "socket:", fmt.Sprintf("Sent %s \"%s\" to %s", event, data, osck.Id()))
	}
}

func (s Server) onOperatorConnection(sock socketio.Socket) {
	// TODO: Operator should be created via API before a connection can
	//  be made
	// Create Operator
	op := operator.Create(sock.Id())

	// Save Operator to datastore
	s.store.Put(op)
}

func (s Server) onClientConnection(sock socketio.Socket) {
	// TODO: Try to recover previous Client/Chat

	// Create Client
	cl := client.Create(sock.Id())

	// Save Client to datastore
	s.store.Put(cl)

	// Create Chat
	ch := chat.Create(cl)

	// Save Chat to datastore
	s.store.Put(ch)

	// Convert to JSON object
	jsonChat, _ := json.Marshal(ch)
	var buffer bytes.Buffer
	buffer.Write(jsonChat)
	buffer.WriteString("\n")

	// Emit chat:new to Operators
	s.emitToOperators("chat:new", buffer.String())

	// Emit chat:new to self
	sock.Emit("chat:new", buffer.String())
}

func (s Server) onConnect(sock socketio.Socket) {

	// Get type GET parameter
	t := sock.Request().URL.Query().Get("type")

	log.Println(INFO, "socket:", fmt.Sprintf("Incoming %s connection %s", t, sock.Id()))

	// TODO: Verify that the socket connection is real
	if t == "operator" {

		s.onOperatorConnection(sock)

	} else if t == "client" {

		s.onClientConnection(sock)
	} else {
		// TODO: Write some proper error handling here, do we close the connection?
		log.Println(ERROR, "socket:", "Unknown chat type specified")
	}

	// Save Socket to use later
	s.Sockets[sock.Id()] = sock

	sock.On("client:message", s.onClientMessage(sock))
	sock.On("operator:message", s.onOperatorMessage(sock))

	sock.On("disconnection", func() {
		log.Println(DEBUG, "socket:", fmt.Sprintf("%s disconnected", sock.Id()))

		// TODO: Save chat?

		delete(s.Sockets, sock.Id())
	})

}

func (s Server) onClientMessage(sock socketio.Socket) func(string) {
	return func(msg string) {

		log.Println(DEBUG, "client", fmt.Sprintf("%s: %s", sock.Id(), msg))

		var m chat.Message

		// String to JSON
		json.Unmarshal([]byte(msg), &m)

		// Save Message to datastore
		s.store.Put(m)

		// Update Operators of the new messages
		s.emitToOperators("client:message", msg)
	}
}

func (s Server) onOperatorMessage(sock socketio.Socket) func(string) {
	return func(msg string) {

		log.Println(DEBUG, "operator", fmt.Sprintf("%s: %s", sock.Id(), msg))

		var m chat.Message

		// String to JSON
		json.Unmarshal([]byte(msg), &m)

		// Save Message to datastore
		s.store.Put(m)

		// Update Client with new message
		clsck, ok := s.Sockets[m.Chat]

		if !ok {
			log.Println(WARNING, "socket:", "Client went away")
			return
		}

		clsck.Emit("operator:message", msg, nil)
	}
}

func (s Server) onError(sock socketio.Socket, err error) {

	// TODO: Write some proper error handling here
	log.Println(ERROR, "socket:", err)
}
