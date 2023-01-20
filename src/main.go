package main

import (
	domain "p2p-messenger/src/domain/services"
	network "p2p-messenger/src/network/providers"
	storage "p2p-messenger/src/storage/providers"
	ui "p2p-messenger/src/ui/providers"

	"flag"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func main() {
	config := network.Config{
		BootstrapPeers: dht.DefaultBootstrapPeers,
		ProtocolID:     "/chat/1.1.0"}
	flag.Var(&config.ListenAddresses, "listen", "Adds a multiaddress to the listen list")
	flag.Parse()

	networkProvider := network.NewP2PNetworkProvider(config)
	uiProvider := ui.NewFyneUIProvider()
	storageProvider := storage.NewInMemoryStorageProvider()
	m := domain.NewMessenger(&networkProvider, &uiProvider, &storageProvider)
	m.Run()
}
