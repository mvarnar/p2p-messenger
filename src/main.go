package main

import (
	domain "p2p-messenger/src/domain/services"
	network "p2p-messenger/src/network/providers"
	ui "p2p-messenger/src/ui/providers"
)

func main() {
	m := domain.NewMessenger(&network.P2PNetworkProvider{}, &ui.FyneUIProvider{})
	m.Run()
}
