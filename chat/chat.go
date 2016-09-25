package api

import (
  "fmt"

  "github.com/mihok/letschat-daemon/operator"
  "github.com/mihok/letschat-daemon/user"
)

type Chat struct {
  id string
  user user.User
  operator operator.Operator
  creationTime uint64
  updatedTime uint64
}

func Create(chat Chat) Chat {

}

func (this *Chat) String() string {
  // return fmt.Sprintf("%s: %s [%s %s]", this.id, this.operator.UserName, this.FirstName, this.LastName)
  return this.id
}