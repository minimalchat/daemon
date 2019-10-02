package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	// TODO: Move away from this socket library, it is no longer maintained
	"github.com/minimalchat/go-socket.io" // Socket
	// "github.com/googollee/go-socket.io" // Socket

	"github.com/minimalchat/daemon/pkg/api/chat"
	"github.com/minimalchat/daemon/pkg/api/client"
	"github.com/minimalchat/daemon/pkg/api/operator"
	"github.com/minimalchat/daemon/pkg/api/webhook"
	"github.com/minimalchat/daemon/pkg/store" // InMemory Database
)

/*
Category defines the type of socket connection. */
type Category string

/*
Categories currently allowed */
const (
	OPERATOR Category = "operator"
	CLIENT   Category = "client"
)

/*
Conn holds the basic information for a socket connection. */
type Conn struct {
	server *Server

	// Socketio connection
	raw      socketio.Socket
	category Category

	send chan *Message
}

/*
Message provides all the details to send or receive a message from a socket
connection. */
type Message struct {
	event   string
	message string
	target  string
}

/*
Listen waits for an incoming message to send out through the connection. */
func (c Conn) Listen() {
	// Defer closing the socket
	defer func() {
		log.Println(DEBUG, "socket:", fmt.Sprintf("%s disconnected", c.raw.Id()))

		c.raw.Disconnect()
	}()

	// Listen for send channel messages and emit them
	for {
		select {
		case data, ok := <-c.send:
			if !ok {
				log.Println(WARNING, "socket:", fmt.Sprintf("Server closed %s channel", c.raw.Id()))
				return
			}

			if data.message == "" {
				log.Println(DEBUG, "socket:", fmt.Sprintf("Emitting '%s' to %s", data.event, c.raw.Id()))

			} else {
				log.Println(DEBUG, "socket:", fmt.Sprintf("Emitting '%s' to %s '%s'", data.event, c.raw.Id(), data.message))
			}

			c.raw.Emit(data.event, data.message)
		}
	}
}

func (c Conn) runWebhooks(e []string, d []byte) {
	for i := 0; i < len(e); i++ {
		w, err := webhook.GetByEventType(c.server.store, e[i])
		if err != nil {
			log.Println(WARNING, "webhooks:", err)
			continue
		} else {
			log.Println(DEBUG, "webhooks:", fmt.Sprintf("Processing webhooks for '%s' (%d found)", e[i], len(w)))

			for j := 0; j < len(w); j++ {
				// Run the Webhook, sending event and data along to the
				//  Webhook's defined endpoint
				err := w[j].Run(e[i], d, c.server.ID)
				if err != nil {
					log.Println(WARNING, "webhooks:", fmt.Sprintf("%s:", e[i]), err)
				}
			}
		}
	}
}

// Client Functions

func (c Conn) onClientConnection(sid string) {

	var cl *client.Client
	var ch *chat.Chat
	var storeBuffer store.Keyer
	var event string

	// Get all Clients
	storeBuffer, _ = c.server.store.Get(fmt.Sprintf("client.%s", sid))

	if storeBuffer != nil {
		// Hijack Client object with new Socket ID
		cl = &client.Client{
			FirstName: storeBuffer.(*client.Client).FirstName,
			LastName:  storeBuffer.(*client.Client).LastName,
			Name:      storeBuffer.(*client.Client).Name,
			Uid:       storeBuffer.(*client.Client).Uid,
			Sid:       c.raw.Id(),
		}
	}

	if cl == nil {
		// Create Client Object
		cl = client.Create(c.raw.Id())

		// Save Client Object to Data Store
		c.server.store.Put(cl)

		// Create Chat Object
		ch = chat.Create(cl)

		// Save Chat Object to Data Store
		c.server.store.Put(ch)

		event = "chat:new"

		// Call any webhooks if they exist
		b, err := json.Marshal(ch)
		if err != nil {
			log.Println(WARNING, "client", err)
		} else {
			c.runWebhooks([]string{
				webhook.EventNewChat,
				webhook.EventNewClient,
			}, b)
		}
	} else {
		// Save Client Object to Data Store with updated Sid
		c.server.store.Put(cl)

		// Get Chat Object
		storeBuffer, _ = c.server.store.Get(fmt.Sprintf("chat.%s", cl.Uid))

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
		c.server.store.Put(ch)

		event = "chat:existing"
	}

	// Convert Chat to JSON
	chJSON, _ := json.Marshal(ch)
	var buffer bytes.Buffer
	buffer.Write(chJSON)
	// buffer.WriteString("\n")

	m := Message{
		event:   event,
		message: buffer.String(),
		target:  "",
	}

	// Broadcast Chat to Operators
	c.server.broadcastToOperators <- &m

	// Send Chat back to Client
	c.send <- &m
}

func (c Conn) onClientMessage(m string) {
	// Create message from JSON
	var msg chat.Message

	json.Unmarshal([]byte(m), &msg)

	log.Println(DEBUG, "client", fmt.Sprintf("%s: %s", c.raw.Id(), msg.Content))

	//  Save Message to Data Store
	c.server.store.Put(msg)

	// TODO:
	//  Update Chat Object
	//  Save Chat Object to Data Store?
	// Run any webhooks if they exist
	c.runWebhooks([]string{webhook.EventNewClientMessage}, []byte(m))

	// Broadcast to Operators
	c.server.broadcastToOperators <- &Message{
		event:   "client:message",
		message: m,
		target:  "",
	}
}

func (c Conn) onClientTyping(m string) {
	// Create message from JSON
	var msg chat.Message

	json.Unmarshal([]byte(m), &msg)

	log.Println(DEBUG, "client", fmt.Sprintf("%s: typing ...", c.raw.Id()))

	c.server.broadcastToOperators <- &Message{
		event:   "client:typing",
		message: m,
		target:  "",
	}
}

// Operator Functions

func (c Conn) onOperatorConnection(id string, t string) {

	var o *operator.Operator
	// TODO: Is this the best way to go about providing access controls?
	// Is id set?
	// Is token set?
	// Get operator with these variables
	// TODO: What should we do if there is no id/token (access ID,
	//  access token)?
	// TODO: This is the only way we can find the operator right now, we
	//  need to improve the InMemory store to handle querying
	operators, err := c.server.store.Search("operator.")
	if err != nil {
		// TODO: What should happen here?
		log.Println(ERROR, "operator", fmt.Sprintf("Something unexpected happened"))
	}

	for _, op := range operators {
		log.Println(DEBUG, "operator", "Does operator match", fmt.Sprintf("(%s == %s)", id, op.(*operator.Operator).Aid))
		if op.(*operator.Operator).Aid == id &&
			op.(*operator.Operator).Atoken == t {
			o = op.(*operator.Operator)
			break
		}
	}

	if o != nil {
		// TODO: Currently the whole apparatus works off of the Uid..
		//  So this feels weird.
		// Update the Uid..
		o.Uid = c.raw.Id()
	} else {
		// If there is no result from the store, create new Operator Object
		o = operator.Create(c.raw.Id())

		b, err := json.Marshal(o)
		if err != nil {
			log.Println(ERROR, "operator", err)
		} else {
			c.runWebhooks([]string{webhook.EventNewOperator}, b)
		}
	}

	// Save Operator Object to Data Store
	c.server.store.Put(o)

	b, err := json.Marshal(o)
	if err != nil {
		log.Println(ERROR, "operator", err)
	}

	// Broadcast the new Operator to all Operators
	c.server.broadcastToOperators <- &Message{
		event:   "operator:new",
		message: string(b),
	}
}

func (c Conn) onOperatorMessage(m string) {
	// Create message from JSON
	var msg chat.Message

	json.Unmarshal([]byte(m), &msg)

	log.Println(DEBUG, "operator", fmt.Sprintf("%s: %s", c.raw.Id(), msg.Content))

	// Save Message to Data Store
	c.server.store.Put(msg)

	// TODO: Update Chat Object?
	// TODO: Save Chat Object to Data Store?

	storeBuffer, _ := c.server.store.Get(fmt.Sprintf("client.%s", msg.Chat))

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

	c.server.broadcastToClient <- &Message{
		event:   "operator:message",
		message: m,
		target:  cl.Sid,
	}
}

func (c Conn) onOperatorTyping(m string) {
	// Create message from JSON
	var msg chat.Message

	json.Unmarshal([]byte(m), &msg)

	log.Println(DEBUG, "operator", fmt.Sprintf("%s: typing ...", c.raw.Id()))

	storeBuffer, _ := c.server.store.Get(fmt.Sprintf("client.%s", msg.Chat))

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

	c.server.broadcastToClient <- &Message{
		event:   "operator:typing",
		message: m,
		target:  cl.Sid,
	}
}
