package chat

import (
	"fmt"
	"time"
)

/*
Message struct defines the object that stores an individual communication */
type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Chat      string    `json:"chat"`
}

func (m *Message) String() string {
	// return fmt.Sprintf("%s: %s [%s %s]", this.id, this.operator.UserName, this.FirstName, this.LastName)
	return m.Content
}

/*
StoreKey defines a key for a DataStore to reference this item */
func (m Message) StoreKey() string {
	return fmt.Sprintf("message.%s-%d", m.Chat, m.Timestamp.Unix())
}
