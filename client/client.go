package client

import (
  "fmt"

  "github.com/wayn3h0/go-uuid" // UUID (RFC 4122)

  "github.com/mihok/letschat-daemon/person"
)


// Operator

type Client struct {
  person.Person
  Name string `json:"name"`
  ID string `json:"id"`
}

func Create(client Client) *Client {
  if (client.ID == "") {
    uuid, _ := uuid.NewRandom()
    client.ID = uuid.String()
  }

  return &client
}

func (this Client) String() string {
  return fmt.Sprintf("%s [%s]", this.Name, this.Name)
}

func (this Client) StoreKey() string {
  return fmt.Sprintf("client.%s", this.ID)
}
