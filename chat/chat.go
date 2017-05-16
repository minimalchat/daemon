package chat

import (
	"fmt"
	"time"

	"github.com/wayn3h0/go-uuid" // UUID (RFC 4122)

	"github.com/minimalchat/daemon/client"
	"github.com/minimalchat/daemon/operator"
	// "github.com/minimalchat/daemon/person"
)

/*
Chat struct defines communication session */
type Chat struct {
	UID          string             `json:"id"`
	Client       *client.Client     `json:"client"`
	Operator     *operator.Operator `json:"operator"`
	CreationTime time.Time          `json:"creation_time"`
	UpdatedTime  time.Time          `json:"update_time"`
	Open         bool               `json:"open"`
}

/*
Create builds a new `Chat` session*/
func Create(chat Chat) *Chat {
	if chat.UID == "" {
		uuid, _ := uuid.NewRandom()
		chat.UID = uuid.String()
	}

	return &chat
}

func (c *Chat) String() string {
	// return fmt.Sprintf("%s: %s [%s %s]", this.id, this.operator.UserName, this.FirstName, this.LastName)
	return c.UID
}

/*
StoreKey defines a key for a DataStore to reference this item */
func (c Chat) StoreKey() string {
	return fmt.Sprintf("chat.%s", c.UID)
}
