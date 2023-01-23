package domain

import (
	entites "p2p-messenger/src/domain/entities"
)

type StorageProvider interface {
	GetContacts() []entites.Contact
	AddNewContact(contact entites.Contact)
	SaveUserId(userId string)
	GetUserId() string
}
