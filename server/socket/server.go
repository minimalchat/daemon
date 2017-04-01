package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io" // Socket

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
Server is the socket.io abstraction for Minimal Chat */
type Server struct {
	// TODO: Poor data structure, should thing of something smarter?
	Operators map[string]*operator.Operator
	Clients   map[string]*client.Client
	Chats     map[string]*chat.Chat
	Server    *socketio.Server
}

/*
Listen creates a new Server instance and begins listening for ws://
connections. */
func Listen(ds *store.InMemory) (*Server, error) {
	log.Println(DEBUG, "socket:", "Starting WebSocket server ...")

	srv, err := socketio.NewServer(nil)
	sck := Server{
		Operators: make(map[string]*operator.Operator),
		Clients:   make(map[string]*client.Client),
		Chats:     make(map[string]*chat.Chat),
		Server:    srv,
	}

	// TODO: Return an error instead
	if err != nil {
		return nil, err
	}

	srv.On("connection", sck.onConnection(ds))
	srv.On("error", sck.onError)

	return &sck, nil
}

/*
ServeHTTP serves the socket.io client script */
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Server.ServeHTTP(w, r)
}

func (s Server) emitToOperators(event string, data string) {
	// if (event == nil) {
	//   log.Println(WARNING, "Unknown event to emit")
	//   return
	// }

	// Update Operators of the new messages
	for _, op := range s.Operators {
		log.Println(DEBUG, "socket:", fmt.Sprintf("Sending %s \"%s\" to %s", event, data, op.Socket.Id()))

		op.Socket.Emit(event, data, nil)
	}
}

func (s Server) onOperatorConnection(ds *store.InMemory, sock socketio.Socket) {
	// Create Operator
	// TODO: Pull Operator from DB if we already know them
	s.Operators[sock.Id()] = operator.Create(
		operator.Operator{
			// FirstName: "Operator",
			// LastName: "Steve",
			UserName: "steve",
		},
		sock,
	)

	// Save Operator to DB
	ds.Put(*s.Operators[sock.Id()])

}

func (s Server) onClientConnection(ds *store.InMemory, sock socketio.Socket) {
	// Create Client
	// TODO: See if we can "recall" if this is a returning client?
	s.Clients[sock.Id()] = client.Create(
		client.Client{
			Name: "Site Visitor",
		},
		sock,
	)

	// Save Client to DB
	ds.Put(*s.Clients[sock.Id()])

	// Create Chat
	// TODO: See if we can "recall" the returning chat?
	s.Chats[sock.Id()] = chat.Create(chat.Chat{
		Client:       s.Clients[sock.Id()],
		Operator:     nil,
		Open:         true,
		CreationTime: time.Now(),
		UpdatedTime:  time.Now(),
	})

	// Save Chat to DB
	ds.Put(*s.Chats[sock.Id()])

	jsonChat, _ := json.Marshal(s.Chats[sock.Id()])
	var buffer bytes.Buffer
	buffer.Write(jsonChat)
	buffer.WriteString("\n")

	// Emit to Operators
	s.emitToOperators("chat:new", buffer.String())

}

func (s Server) onConnection(ds *store.InMemory) func(sock socketio.Socket) {
	return func(sock socketio.Socket) {
		log.Println(INFO, "socket:", fmt.Sprintf("Incoming connection %s %s", sock.Id(), sock.Request().URL.Query().Get("type")))

		t := sock.Request().URL.Query().Get("type")

		// TODO: Verify that the socket connection is real
		if t == "operator" {

			s.onOperatorConnection(ds, sock)

		} else if t == "client" {

			s.onClientConnection(ds, sock)

		} else {

			// TODO: Write some proper error handling here, do we close the connection?
			log.Println(ERROR, "socket:", "Unknown chat type specified")
		}

		sock.On("client:message", s.onClientMessage(ds, sock))
		sock.On("operator:message", s.onOperatorMessage(ds, sock))

		// Disconnection event
		sock.On("disconnection", func() {
			log.Println(DEBUG, "socket:", fmt.Sprintf("%s disconnected", sock.Id()))

			// TODO: Save chat?

			if t == "operator" {

				delete(s.Operators, sock.Id())
			} else if t == "client" {

				delete(s.Clients, sock.Id())
			}
		})
	}
}

func (s Server) onClientMessage(ds *store.InMemory, sock socketio.Socket) func(msg string) {
	return func(msg string) {

		log.Println(DEBUG, "client", fmt.Sprintf("%s: %s", sock.Id(), msg))

		// Create Message
		m := chat.Message{
			Timestamp: time.Now(),
			Content:   msg,
			Author:    s.Clients[sock.Id()].StoreKey(),
			Chat:      s.Chats[sock.Id()].UID,
		}

		// Save Message to DB
		ds.Put(m)

		// Update Operators of the new messages
		s.emitToOperators("client:message", msg)
	}
}

func (s Server) onOperatorMessage(ds *store.InMemory, sock socketio.Socket) func(msg string) {
	return func(msg string) {
		// TODO: Get chat from message, and then update it/send to correct client

		log.Println(DEBUG, "operator", fmt.Sprintf("%s: %s", sock.Id(), msg))

		// Create Message
		// m := chat.Message{
		//   Timestamp: time.Now(),
		//   Content: msg,
		//   Author: this.Operators[sock.Id()].StoreKey(),
		//   Chat: ch.ID,
		// }

		// Save Message to DB
		// ds.Put(m)
	}
}

func (s Server) onError(sock socketio.Socket, err error) {

	// TODO: Write some proper error handling here
	log.Println(ERROR, "socket:", err)
}
