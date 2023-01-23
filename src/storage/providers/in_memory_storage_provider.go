package storage

import (
	entites "p2p-messenger/src/domain/entities"
)

type InMemoryStorageProvider struct {
	contacts []entites.Contact
}

func NewInMemoryStorageProvider() *InMemoryStorageProvider {
	return &InMemoryStorageProvider{contacts: make([]entites.Contact, 0)}
}

func (p *InMemoryStorageProvider) GetContacts() []entites.Contact {
	return p.contacts
}

func (p *InMemoryStorageProvider) AddNewContact(contact entites.Contact) {
	p.contacts = append(p.contacts, contact)
}
