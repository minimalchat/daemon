package client

import (
  "fmt"

  "github.com/mihok/letschat-daemon/person"
)


// Operator

type Client struct {
  person.Person
  Name string `json:"name"`
  Id string `json:"id"`
}

func (this Client) String() string {
  return fmt.Sprintf("%s [%s]", this.Name, this.Name)
}

func (this Client) StoreKey() string {
  return fmt.Sprintf("client.%s", this.Id)
}
