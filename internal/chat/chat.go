package chat

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"ID,omitempty"`
	ClientID  uuid.UUID `json:"client_ID"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	Username  string    `json:"username"`
}

type Chat struct {
	Usernames map[uuid.UUID]string
}

func (c Chat) ParseMessage(message Message) Message {
	username, ok := c.Usernames[message.ClientID]
	if ok {
		message.Username = username
	}

	return message
}
