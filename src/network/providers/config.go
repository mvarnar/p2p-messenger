package network

type Config struct {
	BootstrapPeers  addrList
	ListenAddresses addrList
	ProtocolID      string
}
