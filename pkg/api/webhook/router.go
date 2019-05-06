package webhook

import (
	"errors"
)

/*
Create a new webhook object based on the given event string. */
func Create(event string) *Webhook {
	return &Webhook{}
}

/*
UnmarshalJSON converts a JSON string (as a byte array) into a Webhook object. */
func (w *Webhook) UnmarshalJSON(data []byte) error {
	return errors.New("not implemented")
}

/*
MarshalJSON converts a Webhook object into a JSON string returned as a byte
array. */
func (w Webhook) MarshalJSON() ([]byte, error) {
	return nil, errors.New("not implemented")
}

/*
Key implements the Keyer interface of the Store and returns a string used for
storing the Webhook in memory. */
func (w Webhook) Key() string {
	return ""
}
