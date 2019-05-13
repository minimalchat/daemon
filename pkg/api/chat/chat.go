package chat

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/minimalchat/daemon/pkg/api/client"
)

/*
Create takes a Client object and returns a new Chat session. */
func Create(cl *client.Client) *Chat {
	now := time.Now()

	// Get the unix timestamp (seconds since epoch)
	seconds := now.Unix()
	nanos := int32(now.Sub(time.Unix(seconds, 0)))

	ts := &timestamp.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}

	c := Chat{
		CreationTime: ts,
		UpdatedTime:  ts,
		Open:         true,
		Uid:          cl.Uid,
		Client:       cl,
	}

	return &c
}

/*
UnmarshalJSON converts a JSON string (as a byte array) into a Chat object. */
func (c *Chat) UnmarshalJSON(data []byte) error {
	u := jsonpb.Unmarshaler{}
	buf := bytes.NewBuffer(data)

	return u.Unmarshal(buf, &*c)
}

/*
MarshalJSON converts a Chat object into a JSON string returned as a byte
array. */
func (c Chat) MarshalJSON() ([]byte, error) {
	m := jsonpb.Marshaler{}
	var buf bytes.Buffer

	if err := m.Marshal(&buf, &c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*
Key implements the Keyer interface of the Store and returns a string used for
storing the Webhook in memory. */
func (c Chat) Key() string {
	return fmt.Sprintf("chat.%s", c.Uid)
}
