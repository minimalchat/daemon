package chat

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/minimalchat/daemon/client"
)

/*
Create builds a new `Chat` session*/
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

func (c *Chat) UnmarshaJSON(data []byte) error {
	u := jsonpb.Unmarshaler{}
	buf := bytes.NewBuffer(data)

	if err := u.Unmarshal(buf, &*c); err != nil {
		return err
	}

	return nil
}

func (c Chat) MarshalJSON() ([]byte, error) {
	m := jsonpb.Marshaler{}
	var buf bytes.Buffer

	if err := m.Marshal(&buf, &c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c Chat) StoreKey() string {
	return fmt.Sprintf("chat.%s", c.Uid)
}
