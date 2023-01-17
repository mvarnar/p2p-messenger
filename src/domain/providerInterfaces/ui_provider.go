package domain

import (
	domain "p2p-messenger/src/domain/entities"
)

type UIProvider interface {
	ShowNemIncomingMessage(message domain.Message)
	Run()
}
