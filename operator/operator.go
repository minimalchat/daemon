package operator

import (
  "fmt"

  "github.com/mihok/lets-chat/person"
)


// Operator

type Operator struct {
  person.Person
  UserName string
}

func (this *Operator) String() string {
  return fmt.Sprintf("%s [%s %s]", this.UserName, this.FirstName, this.LastName)
}
