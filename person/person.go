package person

import "fmt"

/*
Person struct defines an `Author` of a communication. It could be either a `Operator` or `Client` */
type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

/*
Author interfaces all have an `ID()` */
type Author interface {
	ID() string
}

func (p *Person) String() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}
