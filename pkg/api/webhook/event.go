package webhook

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang-plus/uuid" // UUID (RFC 4122)
	"github.com/golang/protobuf/jsonpb"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

// Event Types
const (
	// EVENT_NEW_MESSAGE  string = "chat:new_message"
	EventNewChat     string = "chat:new"
	EventNewOperator string = "operator:new"
	// EVENT_NEW_OPERATOR_MESSAGE string = "operator:new_message"
	EventNewClient        string = "client:new"
	EventNewClientMessage string = "client:message"
)

/*
CreateEvent a new Event object based on a given Type (t) string */
func CreateEvent(t string) *Event {
	e := Event{
		Type: t,
	}

	// Generate ID
	uuid, _ := uuid.NewRandom()
	e.Id = uuid.String()

	// Set Create/Update timestamps
	now := time.Now()
	seconds := now.Unix()
	nanos := int32(now.Sub(time.Unix(seconds, 0)))

	ts := &timestamp.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}

	e.CreationTime = ts

	return &e
}

/*
UnmarshalJSON converts a JSON string (as a byte array) into a Event object. */
func (e *Event) UnmarshalJSON(data []byte) error {
	u := jsonpb.Unmarshaler{}
	buf := bytes.NewBuffer(data)

	return u.Unmarshal(buf, &*e)
}

/*
MarshalJSON converts a Event object into a JSON string returned as a byte
array. */
func (e Event) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	m := jsonpb.Marshaler{}

	if err := m.Marshal(&buf, &e); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*
Key implements the Keyer interface of the Store and returns a string used for
storing the Webhook in memory. */
func (e Event) Key() string {
	return fmt.Sprintf("event.%d.%s", e.CreationTime.GetSeconds(), e.Id)
}
