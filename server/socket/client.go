package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	// TODO: Move away from this socket library, it is no longer maintained
	"github.com/minimalchat/go-socket.io" // Socket
	// "github.com/googollee/go-socket.io" // Socket

	"github.com/minimalchat/daemon/chat"
	"github.com/minimalchat/daemon/client"
	"github.com/minimalchat/daemon/operator"
	"github.com/minimalchat/daemon/store" // InMemory Database
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

	// Access control
	// accessId    string
	// accessToken string

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

			if data.message == "" {
				log.Println(DEBUG, "socket:", fmt.Sprintf("Emitting '%s' to %s", data.event, s.conn.Id()))

			} else {
				log.Println(DEBUG, "socket:", fmt.Sprintf("Emitting '%s' to %s '%s'", data.event, s.conn.Id(), data.message))
			}

			s.conn.Emit(data.event, data.message)
		}
	}
}

// Client Functions

func (s Socket) onClientConnection(sessionId string) {

	var cl *client.Client
	var ch *chat.Chat
	var storeBuffer store.StoreKeyer
	var event string

	// Get all Clients
	storeBuffer, _ = s.server.store.Get(fmt.Sprintf("client.%s", sessionId))

	if storeBuffer != nil {
		// Hijack Client object with new Socket ID
		cl = &client.Client{
			FirstName: storeBuffer.(*client.Client).FirstName,
			LastName:  storeBuffer.(*client.Client).LastName,
			Name:      storeBuffer.(*client.Client).Name,
			Uid:       storeBuffer.(*client.Client).Uid,
			Sid:       s.conn.Id(),
		}
	}

	if cl == nil {
		// Create Client Object
		cl = client.Create(s.conn.Id())

		// Save Client Object to Data Store
		s.server.store.Put(cl)

		// Create Chat Object
		ch = chat.Create(cl)

		// Save Chat Object to Data Store
		s.server.store.Put(ch)

		event = "chat:new"
	} else {
		// Save Client Object to Data Store with updated Sid
		s.server.store.Put(cl)

		// Get Chat Object
		storeBuffer, _ = s.server.store.Get(fmt.Sprintf("chat.%s", cl.Uid))

		// Hijack the Chat Object
		ch = &chat.Chat{
			CreationTime: storeBuffer.(*chat.Chat).CreationTime,
			// TODO: Update UpdatedTime to now
			UpdatedTime: storeBuffer.(*chat.Chat).UpdatedTime,
			Open:        storeBuffer.(*chat.Chat).Open,
			Uid:         storeBuffer.(*chat.Chat).Uid,
			Client:      cl,
		}

		// Save Chat Object to Data Store with updated Client
		s.server.store.Put(ch)

		event = "chat:existing"
	}

	// Convert Chat to JSON
	chJson, _ := json.Marshal(ch)
	var buffer bytes.Buffer
	buffer.Write(chJson)
	// buffer.WriteString("\n")

	sm := SocketMessage{
		event:   event,
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

func (s Socket) onClientTyping(data string) {
	// Create message from JSON
	var msg chat.Message

	json.Unmarshal([]byte(data), &msg)

	log.Println(DEBUG, "client", fmt.Sprintf("%s: typing ...", s.conn.Id()))

	s.server.broadcastToOperators <- &SocketMessage{
		event:   "client:typing",
		message: data,
		target:  "",
	}
}

// Operator Functions

func (s Socket) onOperatorConnection(accessId string, accessToken string) {

	var o *operator.Operator
	// TODO: Is this the best way to go about providing access controls?
	// Is accessId set?
	// Is accessToken set?
	// Get operator with these variables
	// TODO: What should we do if there is no accessId/accessToken?
	// TODO: This is the only way we can find the operator right now, we
	//  need to improve the InMemory store to handle querying
	operators, err := s.server.store.Search("operator.")
	if err != nil {
		// TODO: What should happen here?
		log.Println(ERROR, "operator", fmt.Sprintf("Something unexpected happened"))
	}

	for _, op := range operators {
		log.Println(DEBUG, "operator", "Does operator match", fmt.Sprintf("(%s == %s)", accessId, op.(*operator.Operator).Aid))
		if op.(*operator.Operator).Aid == accessId &&
			op.(*operator.Operator).Atoken == accessToken {
			o = op.(*operator.Operator)
			break
		}
	}

	if o != nil {
		// TODO: Currently the whole apparatus works off of the Uid..
		//  So this feels weird.
		// Update the Uid..
		o.Uid = s.conn.Id()
	} else {
		// If there is no result from the store, create new Operator Object
		o = operator.Create(s.conn.Id())
	}

	// Save Operator Object to Data Store
	s.server.store.Put(o)

	b, err := json.Marshal(o)
	if err != nil {
		log.Println(ERROR, "operator", err)
	}

	// Broadcast the new Operator to all Operators
	s.server.broadcastToOperators <- &SocketMessage{
		event:   "operator:new",
		message: string(b),
	}
}

func (s Socket) onOperatorMessage(data string) {
	// Create message from JSON
	var msg chat.Message

	json.Unmarshal([]byte(data), &msg)

	log.Println(DEBUG, "operator", fmt.Sprintf("%s: %s", s.conn.Id(), msg.Content))

	// Save Message to Data Store
	s.server.store.Put(msg)

	// TODO: Update Chat Object?
	// TODO: Save Chat Object to Data Store?

	storeBuffer, _ := s.server.store.Get(fmt.Sprintf("client.%s", msg.Chat))

	if storeBuffer == nil {
		log.Println(ERROR, "operator:", fmt.Sprintf("Client %s does not exist!", msg.Chat))
		return
	}

	// TODO: This could be better, seems kinda hacky
	// Hijack the Client object we need to message to
	cl := &client.Client{
		FirstName: storeBuffer.(*client.Client).FirstName,
		LastName:  storeBuffer.(*client.Client).LastName,
		Name:      storeBuffer.(*client.Client).Name,
		Uid:       storeBuffer.(*client.Client).Uid,
		Sid:       storeBuffer.(*client.Client).Sid,
	}

	s.server.broadcastToClient <- &SocketMessage{
		event:   "operator:message",
		message: data,
		target:  cl.Sid,
	}
}

func (s Socket) onOperatorTyping(data string) {
	// Create message from JSON
	var msg chat.Message

	json.Unmarshal([]byte(data), &msg)

	log.Println(DEBUG, "operator", fmt.Sprintf("%s: typing ...", s.conn.Id()))

	storeBuffer, _ := s.server.store.Get(fmt.Sprintf("client.%s", msg.Chat))

	if storeBuffer == nil {
		log.Println(ERROR, "operator:", fmt.Sprintf("Client %s does not exist!", msg.Chat))
		return
	}

	// TODO: This could be better, seems kinda hacky
	// Hijack the Client object we need to message to
	cl := &client.Client{
		FirstName: storeBuffer.(*client.Client).FirstName,
		LastName:  storeBuffer.(*client.Client).LastName,
		Name:      storeBuffer.(*client.Client).Name,
		Uid:       storeBuffer.(*client.Client).Uid,
		Sid:       storeBuffer.(*client.Client).Sid,
	}

	s.server.broadcastToClient <- &SocketMessage{
		event:   "operator:typing",
		message: data,
		target:  cl.Sid,
	}
}
