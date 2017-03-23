package utils

import (
	"github.com/manveru/faker"
	"github.com/minimalchat/mnml-daemon/chat"
	"time"
)

// MakeDummy returns an array of dummy chat.Messages
func MakeDummy(numMessages int) []chat.Message {

	// set up faker
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}

	messages := make([]chat.Message, numMessages)

	for i := 0; i <= numMessages; i++ {
		newMessage := chat.Message{
			Content:   fake.Paragraph(3, true),
			Author:    fake.Name(),
			Chat:      fake.PostCode(),
			Timestamp: time.Unix(64043302, 23423502),
		}

		// running into errors here.
		// messages[i] = newMessage // < this doesn't work
		messages = append(messages, newMessage) // < works but first half of array is full of junk?
	}

	return messages
}
