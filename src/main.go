package main

import (
	domain "p2p-messenger/src/domain/services"
	network "p2p-messenger/src/network/providers"
	ui "p2p-messenger/src/ui/providers"

	"flag"

	"github.com/ipfs/go-log/v2"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func main() {
	config := network.Config{
		RendezvousString: "rendezvous",
		BootstrapPeers:   dht.DefaultBootstrapPeers,
		ProtocolID:       "/chat/1.1.0"}
	flag.Var(&config.ListenAddresses, "listen", "Adds a multiaddress to the listen list")
	flag.Parse()
	logger := log.Logger("p2p-msgr")

	networkProvider := network.NewP2PNetworkProvider(logger, config)
	uiProvider := ui.NewFyneUIProvider()
	m := domain.NewMessenger(&networkProvider, &uiProvider)
	m.Run()
}
