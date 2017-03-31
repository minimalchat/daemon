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

	for i := 0; i < numMessages; i++ {
		messages[i].Content = fake.Paragraph(3, true)
		messages[i].Author = fake.Name()
		messages[i].Chat = fake.PostCode()
		messages[i].Timestamp = time.Unix(64042203, 23423502)
	}

	return messages
}
