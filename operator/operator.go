package operator

import (
  "fmt"

  // "github.com/wayn3h0/go-uuid" // UUID (RFC 4122)
  "github.com/googollee/go-socket.io" // Socket

  "github.com/minimalchat/mnml-daemon/person"
)


// Operator

type Operator struct {
  person.Person
  UserName string `json:"username"`
  Uid string `json:"id"`
  Socket socketio.Socket `json:"socket"`
}

func Create(operator Operator, sock socketio.Socket) *Operator {
  if (operator.Uid == "") {
    // uuid, _ := uuid.NewRandom()
    operator.Uid = sock.Id()
  }

  operator.Socket = sock

  return &operator
}

func (this Operator) String() string {
  return fmt.Sprintf("%s [%s %s]", this.UserName, this.FirstName, this.LastName)
}

func (this Operator) ID() string {
  return this.UserName
}


func (this Operator) StoreKey() string {
  return fmt.Sprintf("operator.%s", this.ID())
}