package person

import "fmt"

type Person struct {
  FirstName string
  LastName string
}

func (this *Person) String() string {
  return fmt.Sprintf("%s %s", this.FirstName, this.LastName)
}
