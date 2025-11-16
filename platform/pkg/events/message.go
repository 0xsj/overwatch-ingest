// platform/pkg/events/message.go
package events

import (
	"encoding/json"
	"time"
)

// Message represents an event message from the event bus.
type Message struct {
	// Subject is the topic/subject the message was published to
	Subject string

	// Data is the raw message payload
	Data []byte

	// Metadata contains additional message information
	Metadata map[string]string

	// Timestamp is when the message was published
	Timestamp time.Time
}

// UnmarshalJSON unmarshals the message data into a struct.
func (m *Message) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(m.Data, v)
}

// NewMessage creates a new message.
func NewMessage(subject string, data []byte) *Message {
	return &Message{
		Subject:   subject,
		Data:      data,
		Metadata:  make(map[string]string),
		Timestamp: time.Now(),
	}
}