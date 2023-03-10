package domain

import (
	domain "p2p-messenger/src/domain/entities"
)

type NetworkProvider interface {
	GetNewIncomingMessages() <-chan domain.Message
	SendMessage(message domain.Message)
	GetUserId() string
	GetKeyBytes() []byte
	Run(keyBytes []byte)
}
