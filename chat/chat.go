package chat

import (
  "fmt"
  "time"

  "github.com/minimalchat/mnml-daemon/operator"
  "github.com/minimalchat/mnml-daemon/client"
  // "github.com/minimalchat/mnml-daemon/person"
)

type Chat struct {
  ID string `json:"id"`
  Client *client.Client `json:"client"`
  Operator *operator.Operator `json:"operator"`
  CreationTime time.Time `json:"creation_time"`
  UpdatedTime time.Time `json:"update_time"`
  Open bool `json:"open"`
}

func (this *Chat) String() string {
  // return fmt.Sprintf("%s: %s [%s %s]", this.id, this.operator.UserName, this.FirstName, this.LastName)
  return this.ID
}

func (this Chat) StoreKey() string {
  return fmt.Sprintf("chat.%s", this.ID)
}


type Message struct {
  Timestamp time.Time `json:"timestamp"`
  Content string `json:"content"`
  Author string `json:"author"`
  Chat string `json:"chat"`
}

func (this *Message) String() string {
  // return fmt.Sprintf("%s: %s [%s %s]", this.id, this.operator.UserName, this.FirstName, this.LastName)
  return this.Content
}

func (this Message) StoreKey() string {
  return fmt.Sprintf("message.%s-%d", this.Chat, this.Timestamp.Unix())
}