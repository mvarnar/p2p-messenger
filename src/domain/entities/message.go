package domain

import (
	"time"
)

type Message struct{
	SenderId string
	ReceiverId string
	Text string
	SentDatetime time.Time
}