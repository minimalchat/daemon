package chat

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
)

func (m *Message) UnmarshalJSON(data []byte) error {
	u := jsonpb.Unmarshaler{}
	buf := bytes.NewBuffer(data)

	if err := u.Unmarshal(buf, &*m); err != nil {
		return err
	}

	return nil
}

func (m Message) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	mrsh := jsonpb.Marshaler{}

	if err := mrsh.Marshal(&buf, &m); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m Message) StoreKey() string {
	return fmt.Sprintf("message.%s-%d", m.Chat, m.Timestamp.Seconds)
}
