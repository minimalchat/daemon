package operator

import (
	"fmt"

	// "github.com/wayn3h0/go-uuid" // UUID (RFC 4122)
	"github.com/googollee/go-socket.io" // Socket

	"github.com/minimalchat/mnml-daemon/person"
)

/*
Operator struct defines a site owner */
type Operator struct {
	person.Person
	UserName string          `json:"username"`
	UID      string          `json:"id"`
	Socket   socketio.Socket `json:"socket"`
}

/*
Create builds a new `Operator` */
func Create(o Operator, sock socketio.Socket) *Operator {
	if o.UID == "" {
		// uuid, _ := uuid.NewRandom()
		o.UID = sock.Id()
	}

	o.Socket = sock

	return &o
}

func (o Operator) String() string {
	return fmt.Sprintf("%s [%s %s]", o.UserName, o.FirstName, o.LastName)
}

// func (this Operator) ID() string {
// 	return this.UserName
// }

/*
StoreKey defines a key for a DataStore to reference this item */
func (o Operator) StoreKey() string {
	return fmt.Sprintf("operator.%s", o.UserName)
	/*
	   StoreKey defines a key for a DataStore to reference this item */
}
