package operator

import (
	"bytes"
	"fmt"

	"github.com/golang-plus/uuid" // UUID (RFC 4122)
	"github.com/golang/protobuf/jsonpb"
)

/*
Create takes an optional Id string and returns a new Operator. */
func Create(id string) *Operator {
	o := Operator{}

	if id == "" {
		uuid, _ := uuid.NewRandom()
		o.Uid = uuid.String()
	} else {
		o.Uid = id
	}

	return &o
}

/*
GetFullName returns the Operator FirstName and LastName concatenated by a
space. Implementing the Operator as a Person. */
func (o *Operator) GetFullName() string {
	if o != nil {
		return fmt.Sprintf("%s %s", o.GetFirstName(), o.GetLastName())
	}
	return ""
}

/*
UnmarshalJSON converts a JSON string (as a byte array) into a Operator object. */
func (o *Operator) UnmarshalJSON(data []byte) error {
	u := jsonpb.Unmarshaler{}
	buf := bytes.NewBuffer(data)

	return u.Unmarshal(buf, &*o)
}

/*
MarshalJSON converts a Operator object into a JSON string returned as a byte
array. */
func (o Operator) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	m := jsonpb.Marshaler{}

	if err := m.Marshal(&buf, &o); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*
Key implements the Keyer interface of the Store and returns a string used for
storing the Operator in memory. */
func (o Operator) Key() string {
	return fmt.Sprintf("operator.%s", o.Aid)
}
