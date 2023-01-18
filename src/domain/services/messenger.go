package domain

import (
	provider "p2p-messenger/src/domain/provider_interfaces"
)

type Messenger struct {
	networkProvider provider.NetworkProvider
	uiProvider      provider.UIProvider
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

	// todo fyne не может работать вне главного потока
	// этот вызов нарушает логическую изоляцию
	m.uiProvider.Run()
}

func (m *Messenger) readIncomingMessages() {
	for message := range m.networkProvider.GetNewIncomingMessages() {
		m.uiProvider.ShowNemIncomingMessage(message)
	}
}

func (m *Messenger) readOutgoingMessages() {
	for message := range m.uiProvider.GetNewOutgoingMessages() {
		m.networkProvider.SendMessage(message)
	}
}
