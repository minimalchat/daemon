package operator

import (
  "fmt"

  "github.com/mihok/letschat-daemon/person"
)


// Operator

type Operator struct {
  person.Person
  UserName string `json:"username"`
}

func Create(operator Operator) *Operator {
  return &operator
}

func (this Operator) String() string {
  return fmt.Sprintf("%s [%s %s]", this.UserName, this.FirstName, this.LastName)
}

func (this Operator) StoreKey() string {
  return fmt.Sprintf("operator.%s", this.UserName)
}
