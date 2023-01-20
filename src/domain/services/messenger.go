package domain

import (
	helpers "p2p-messenger/src/collections_helpers"
	provider "p2p-messenger/src/domain/provider_interfaces"
)

type Messenger struct {
	networkProvider   provider.NetworkProvider
	uiProvider        provider.UIProvider
	storageProvider   provider.StorageProvider
}
type F struct {
}

func NewMessenger(NetworkProvider provider.NetworkProvider, UIProvider provider.UIProvider) Messenger {
	return Messenger{networkProvider: NetworkProvider, uiProvider: UIProvider}
}

func (m *Messenger) Run() {
	go m.networkProvider.Run()
	go m.readIncomingMessages()
	go m.readOutgoingMessages()
	go m.readNewContacts()
	go m.GetUserId()

	// todo fyne не может работать вне главного потока
	// этот вызов нарушает логическую изоляцию
	m.uiProvider.Run()
}

func (m *Messenger) readIncomingMessages() {
	for message := range m.networkProvider.GetNewIncomingMessages() {
		contacts := m.storageProvider.GetContacts()
		if !helpers.Contains(contacts, message.SenderContact) {
			m.uiProvider.ShowNewContact(message.SenderContact)
			m.storageProvider.AddNewContact(message.SenderContact)
		}
		m.uiProvider.ShowNemIncomingMessage(message)
	}
}

func (m *Messenger) readOutgoingMessages() {
	for message := range m.uiProvider.GetNewOutgoingMessages() {
		m.networkProvider.SendMessage(message)
	}
}

func (m *Messenger) GetUserId() {
	userId := m.networkProvider.GetUserId()
	m.uiProvider.ShowUserId(userId)
}

func (m *Messenger) readNewContacts() {
	for contact := range m.uiProvider.GetNewContacts() {
		m.uiProvider.ShowNewContact(contact)
		m.storageProvider.AddNewContact(contact)
	}
}
