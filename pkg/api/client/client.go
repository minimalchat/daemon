package client

import (
	"bytes"
	"fmt"

	"github.com/golang-plus/uuid" // UUID (RFC 4122)
	"github.com/golang/protobuf/jsonpb"
)

/*
Create takes a Sid identifier string and returns a new Client. */
func Create(sid string) *Client {
	c := Client{
		FirstName: "Site",
		LastName:  "Visitor",
		Name:      "Site Visitor",
	}

	// Generate Client UID
	uuid, _ := uuid.NewRandom()
	c.Uid = uuid.String()

	// Assign Client SID
	c.Sid = sid

	return &c
}

/*
GetFullName returns the Client FirstName and LastName concatenated by a space.
Implementing the Client as a Person. */
func (c *Client) GetFullName() string {
	if c != nil {
		return fmt.Sprintf("%s %s", c.GetFirstName(), c.GetLastName())
	}
	return ""
}

/*
UnmarshalJSON converts a JSON string (as a byte array) into a Client object. */
func (c *Client) UnmarshalJSON(data []byte) error {
	u := jsonpb.Unmarshaler{}
	buf := bytes.NewBuffer(data)

	return u.Unmarshal(buf, &*c)
}

/*
MarshalJSON converts a Client object into a JSON string returned as a byte
array. */
func (c Client) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	m := jsonpb.Marshaler{}

	if err := m.Marshal(&buf, &c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*
Key implements the Keyer interface of the Store and returns a string used for
storing the Client in memory. */
func (c Client) Key() string {
	return fmt.Sprintf("client.%s", c.Uid)
}
