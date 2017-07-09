package operator

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/wayn3h0/go-uuid" // UUID (RFC 4122)
)

/*
Create builds a new `Operator` */
func Create(id string) *Operator {
	o := Operator{
		UserName: "steve",
	}

	if id == "" {
		uuid, _ := uuid.NewRandom()
		o.Uid = uuid.String()
	} else {
		o.Uid = id
	}

	return &o
}

func (o *Operator) GetFullName() string {
	if o != nil {
		return fmt.Sprintf("%s %s", o.GetFirstName(), o.GetLastName())
	}
	return ""
}

func (o *Operator) UnmarshalJSON(data []byte) error {
	u := jsonpb.Unmarshaler{}
	buf := bytes.NewBuffer(data)

	if err := u.Unmarshal(buf, &*o); err != nil {
		return err
	}

	return nil
}

func (o Operator) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	m := jsonpb.Marshaler{}

	if err := m.Marshal(&buf, &o); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (o Operator) StoreKey() string {
	return fmt.Sprintf("operator.%s", o.UserName)
}
