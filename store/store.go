package store

import "time"

type Store interface {
	AddOutgoingMessage(m Message) (int64, error)
	GetOutgoingMessage(id int64) (Message, error)

	GetIncomingMessages(from time.Time, to time.Time) ([]Message, error)

	AddIncomingMessage(m Message) error
	GetOutgoingMessages(s SearchRequest) ([]Message, error)
}

type SearchRequest struct {
	Limit  byte
	Offset *byte
	Phone  *string
	From   *time.Time
	To     *time.Time
}

type Message struct {
	Id    int64
	Sent  time.Time
	Text  string
	Phone string
}
