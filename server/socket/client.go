package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/googollee/go-socket.io" // Socket

	"github.com/minimalchat/daemon/chat"
	"github.com/minimalchat/daemon/client"
	"github.com/minimalchat/daemon/operator"
)

type SocketType string

// Socket Type enum
const (
	OPERATOR SocketType = "operator"
	CLIENT   SocketType = "client"
)

type Socket struct {
	server *Server

	// Socketio connection
	conn     socketio.Socket
	connType SocketType

	send chan *SocketMessage
}

type SocketMessage struct {
	event   string
	message string
	target  string
}

func (s Socket) Listen() {
	// Defer closing the socket
	defer func() {
		log.Println(DEBUG, "socket:", fmt.Sprintf("%s disconnected", s.conn.Id()))

		s.conn.Disconnect()
	}()

	// Listen for send channel messages and emit them
	for {
		select {
		case data, ok := <-s.send:
			if !ok {
				log.Println(WARNING, "socket:", fmt.Sprintf("Server closed %s channel", s.conn.Id()))
				return
			}

			log.Println(DEBUG, "socket:", fmt.Sprintf("Emitting '%s' to %s '%s'", data.event, s.conn.Id(), data.message))
			s.conn.Emit(data.event, data.message)
		}
	}
}

// Client Functions

func (s Socket) onClientConnection() {
	// Create Client Object
	cl := client.Create(s.conn.Id())

	// Save Client Object to Data Store
	s.server.store.Put(cl)

	// Create Chat Object
	ch := chat.Create(cl)

	// Save Chat Object to Data Store
	s.server.store.Put(ch)

	// Convert to JSON
	chJson, _ := json.Marshal(ch)
	var buffer bytes.Buffer
	buffer.Write(chJson)
	// buffer.WriteString("\n")

	sm := SocketMessage{
		event:   "chat:new",
		message: buffer.String(),
		target:  "",
	}

	// Broadcast Chat to Operators
	s.server.broadcastToOperators <- &sm

	// Send Chat back to Client
	s.send <- &sm
}

func (s Socket) onClientMessage(data string) {
	// Create message from JSON
	var msg chat.Message

	json.Unmarshal([]byte(data), &msg)

	log.Println(DEBUG, "client", fmt.Sprintf("%s: %s", s.conn.Id(), msg.Content))

	//  Save Message to Data Store
	s.server.store.Put(msg)

	// TODO:
	//  Update Chat Object
	//  Save Chat Object to Data Store?

	// Broadcast to Operators
	s.server.broadcastToOperators <- &SocketMessage{
		event:   "client:message",
		message: data,
		target:  "",
	}
}

// Operator Functions

func (s Socket) onOperatorConnection() {
	// Create Operator Object
	op := operator.Create(s.conn.Id())

	// Save Operator Object to Data Store
	s.server.store.Put(op)
}

func (s Socket) onOperatorMessage(data string) {
	// Create message from JSON
	var msg chat.Message

	json.Unmarshal([]byte(data), &msg)

	log.Println(DEBUG, "operator", fmt.Sprintf("%s: %s", s.conn.Id(), msg.Content))

	// Save Message to Data Store
	s.server.store.Put(msg)

	// TODO:
	// Update Chat Object
	// Save Chat Object to Data Store?

	s.server.broadcastToClient <- &SocketMessage{
		event:   "operator:message",
		message: data,
		target:  msg.Chat,
	}
}
