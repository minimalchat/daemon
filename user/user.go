package user

import (
  "fmt"

  "github.com/mihok/lets-chat/person"
)


// Operator

type User struct {
  person.Person
  UserName string
}

func (this *User) String() string {
  return fmt.Sprintf("%s [%s %s]", this.UserName, this.FirstName, this.LastName)
}
