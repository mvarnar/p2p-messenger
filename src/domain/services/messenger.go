package domain

import (
	provider "p2p-messenger/src/domain/providerinterfaces"
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
	m.uiProvider.Run()
}

func (m *Messenger) readIncomingMessages(){
	for message := range m.networkProvider.GetNewIncomingMessages(){
		m.uiProvider.ShowNemIncomingMessage(message)
	}
}
