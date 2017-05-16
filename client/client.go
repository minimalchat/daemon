package client

import (
	"fmt"

	// "github.com/wayn3h0/go-uuid" // UUID (RFC 4122)
	"github.com/googollee/go-socket.io" // Socket

	"github.com/minimalchat/daemon/person"
)

/*
Client struct defines a web visitor */
type Client struct {
	person.Person
	Name   string          `json:"name"`
	UID    string          `json:"id"`
	Socket socketio.Socket `json:"socket"`
}

/*
Create builds a new `Client` */
func Create(c Client, sock socketio.Socket) *Client {
	if c.UID == "" {
		// uuid, _ := uuid.NewRandom()
		c.UID = sock.Id()
	}

	c.Socket = sock

	return &c
}

// func (this *Client) Send(msg chat.Message) error {
//   this.socket.Emit("operator:message", msg.Content, func (sock socketio.Socket, data string) {
//     log.Println(DEBUG, "client:", "Sent message")
//   })
//   return nil
// }

func (c Client) String() string {
	return fmt.Sprintf("%s [%s]", c.Name, c.Name)
}

// func (this Client) ID() string {
//  return this.UID
// }

/*
StoreKey defines a key for a DataStore to reference this item */
func (c Client) StoreKey() string {
	return fmt.Sprintf("client.%s", c.UID)
}
