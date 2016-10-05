package person

import "fmt"

type Person struct {
  FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
}

func (this *Person) String() string {
  return fmt.Sprintf("%s %s", this.FirstName, this.LastName)
}
