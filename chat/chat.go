package chat

import (
  "fmt"
  "time"

  "github.com/mihok/letschat-daemon/operator"
  "github.com/mihok/letschat-daemon/client"
)

type Chat struct {
  ID string `json:"id"`
  Client *client.Client `json:"client"`
  Operator *operator.Operator `json:"operator"`
  CreationTime time.Time `json:"creation_time"`
  UpdatedTime time.Time `json:"update_time"`
}

func (this *Chat) String() string {
  // return fmt.Sprintf("%s: %s [%s %s]", this.id, this.operator.UserName, this.FirstName, this.LastName)
  return this.ID
}

func (this Chat) StoreKey() string {
  return fmt.Sprintf("chat.%s", this.ID)
}
