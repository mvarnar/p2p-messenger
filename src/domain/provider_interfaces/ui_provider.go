package domain

import (
	domain "p2p-messenger/src/domain/entities"
)

type UIProvider interface {
	GetNewOutgoingMessages() <-chan domain.Message
	ShowNemIncomingMessage(message domain.Message)
	Run()
}
