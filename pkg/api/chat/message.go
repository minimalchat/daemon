package chat

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

/*
CreateMessage constructs a new `Message` with a default timestamp of now */
func CreateMessage() *Message {
	now := time.Now()

	// Get the unix timestamp (seconds since epoch)
	seconds := now.Unix()
	nanos := int32(now.Sub(time.Unix(seconds, 0)))

	ts := &timestamp.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}

	m := Message{
		Timestamp: ts,
	}

	return &m
}

/*
UnmarshalJSON converts a JSON string (as a byte array) into a Message object */
func (m *Message) UnmarshalJSON(data []byte) error {
	u := jsonpb.Unmarshaler{}
	buf := bytes.NewBuffer(data)

	return u.Unmarshal(buf, &*m)
}

/*
MarshalJSON converts a Message object into a JSON string returned as a byte
array */
func (m Message) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	mrsh := jsonpb.Marshaler{}

	if err := mrsh.Marshal(&buf, &m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*
Key implements the Keyer interface of the Store and returns a string used for
storing the Message in memory. */
func (m Message) Key() string {
	return fmt.Sprintf("message.%s-%d", m.Chat, m.Timestamp.Seconds)
}
