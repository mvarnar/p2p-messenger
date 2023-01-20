package domain

import (
	"time"
)

type Message struct{
	SenderContact Contact
	ReceiverContact Contact
	Text string
	SentDatetime time.Time
}