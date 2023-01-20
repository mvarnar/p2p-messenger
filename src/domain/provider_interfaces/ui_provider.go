package domain

import (
	domain "p2p-messenger/src/domain/entities"
)

type UIProvider interface {
	GetNewOutgoingMessages() <-chan domain.Message
	ShowNemIncomingMessage(message domain.Message)
	ShowUserId(userId string)
	GetNewContacts() <-chan domain.Contact
	ShowNewContact(domain.Contact)
	Run()
}
