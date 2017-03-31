package client

import (
  "fmt"

  // "github.com/wayn3h0/go-uuid" // UUID (RFC 4122)
  "github.com/googollee/go-socket.io" // Socket

  "github.com/minimalchat/mnml-daemon/person"
)


// Operator

type Client struct {
  person.Person
  Name string `json:"name"`
  Uid string `json:"id"`
  Socket socketio.Socket `json:"socket"`
}

func Create(client Client, sock socketio.Socket) *Client {
  if (client.Uid == "") {
    // uuid, _ := uuid.NewRandom()
    client.Uid = sock.Id()
  }

  client.Socket = sock

  return &client
}

// func (this *Client) Send(msg chat.Message) error {
//   this.socket.Emit("operator:message", msg.Content, func (sock socketio.Socket, data string) {
//     log.Println(DEBUG, "client:", "Sent message")
//   })
//   return nil
// }

func (this Client) String() string {
  return fmt.Sprintf("%s [%s]", this.Name, this.Name)
}

func (this Client) ID() string {
  return this.Uid
}


func (this Client) StoreKey() string {
  return fmt.Sprintf("client.%s", this.ID())
}