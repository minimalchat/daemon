package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	// "errors"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-plus/uuid" // UUID (RFC 4122)
	"github.com/golang/protobuf/jsonpb"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"

	"github.com/minimalchat/daemon/pkg/store"
)

/*
CreateWebhook a new Webhook object based on the given event strings. */
func CreateWebhook(uri string, events []string) *Webhook {
	w := Webhook{
		EventTypes: events,
		Uri:        uri,
		Enabled:    true,
	}

	// Generate ID
	id, _ := uuid.NewRandom()
	w.Id = id.String() // fmt.Sprintf("wh-%s", id.String())

	// // Generate Secret
	// h := sha256.New()
	// h.Write([]byte(w.Id))

	t := time.Now()
	// b := make([]byte, 8)
	// binary.LittleEndian.PutUint64(b, uint64(t.Unix()))
	// h.Write(b)

	// h.Write([]byte(salt))

	secret, _ := uuid.NewRandom()
	w.Secret = secret.String() // fmt.Sprintf("whsec-%s", secret.String())

	// w.Secret = fmt.Sprintf("%x", h.Sum(nil))

	// Set Create/Update timestamps
	seconds := t.Unix()
	nanos := int32(t.Sub(time.Unix(seconds, 0)))

	ts := &timestamp.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}

	w.CreationTime = ts
	w.UpdatedTime = ts

	return &w
}

/*
GetByEventType returns one or more Webhooks that are for a specified event */
func GetByEventType(ds *store.InMemory, event string) ([]*Webhook, error) {
	var result []*Webhook
	// TODO: Figure out if this is the best place for this function

	// Get all webhooks
	ws, err := ds.Search("webhook.")
	if err != nil {
		return nil, err
	}

	// Iterate over each and check its EventTypes to see if there is a match
	//  with event
	for i := 0; i < len(ws); i++ {
		eventTypes := ws[i].(*Webhook).EventTypes

		// TODO: Not sure if this is more/less expensive than looping
		//  through each event type individually and doing a comparison.
		eventTypeString := strings.Join(eventTypes, "")

		if strings.Contains(eventTypeString, event) {
			result = append(result, ws[i].(*Webhook))
		}
	}

	return result, nil
	// return nil, errors.New("not implemented")
}

/*
Run executes the Webhook sending an HTTP request with an Event to the defined
URI. We generate a Mnml-Signature header that looks something along the lines
of:

Mnml-Signature: t=1492774577,
    v1=5257a869e7ecebeda32affa62cdca3fa51cad7e77a0e56ff536d0ce8e108d8bd

Where t is the unix timestamp in seconds, and v1 is an HMAC signature using
SHA-256. This can be used to validate that the signature is coming from the
daemon service and not anywhere else.

To manually verify the signature, take t and concat it with the JSON payload
received, separated by a '.' dot. Use the Webhook's secret to sign the
concatenated timestamp.payload and it should equal v1. */
func (w *Webhook) Run(t string, d []byte, n string) error {
	c := &http.Client{}

	e := CreateEvent(t)
	e.Data = string(d)
	e.SourceId = n

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", w.Uri, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	sigPayload := fmt.Sprintf("%d.%s", e.CreationTime.GetSeconds(), b)

	// Generate Signature based off of the webhook secret
	hmac := hmac.New(sha256.New, []byte(w.Secret))
	hmac.Write([]byte(sigPayload))

	sig := fmt.Sprintf("t=%d,v1=%x", e.CreationTime.GetSeconds(), hmac.Sum(nil))

	req.Header.Set("Mnml-Signature", sig)

	resp, err := c.Do(req)

	// TODO: Do something with the response, record the request?
	if resp != nil {
		log.Println(DEBUG, "webhook:", fmt.Sprintf("Sent event '%s' to %s and got: %v", t, w.Uri, resp.Status))
	}

	return err
}

/*
UnmarshalJSON converts a JSON string (as a byte array) into a Webhook object. */
func (w *Webhook) UnmarshalJSON(data []byte) error {
	u := jsonpb.Unmarshaler{}
	buf := bytes.NewBuffer(data)

	return u.Unmarshal(buf, &*w)
}

/*
MarshalJSON converts a Webhook object into a JSON string returned as a byte
array. */
func (w Webhook) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	m := jsonpb.Marshaler{}

	if err := m.Marshal(&buf, &w); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*
Key implements the Keyer interface of the Store and returns a string used for
storing the Webhook in memory. */
func (w Webhook) Key() string {
	return fmt.Sprintf("webhook.%s", w.Id)
}
