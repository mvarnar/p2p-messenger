package network

import (
	"bufio"
	"context"
	"fmt"
	entities "p2p-messenger/src/domain/entities"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	maddr "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-log/v2"
)

type P2PNetworkProvider struct {
	logger          *log.ZapEventLogger
	config          Config
	messagesChannel chan entities.Message
	readWriter      *bufio.ReadWriter
}

func NewP2PNetworkProvider(Logger *log.ZapEventLogger, Config Config) P2PNetworkProvider {
	return P2PNetworkProvider{logger: Logger, config: Config, messagesChannel: make(chan entities.Message, 100)}
}

func (p *P2PNetworkProvider) GetNewIncomingMessages() <-chan entities.Message {
	return p.messagesChannel
}

func (p *P2PNetworkProvider) SendMessage(message entities.Message) {
	_, err := p.readWriter.WriteString(fmt.Sprintf("%s\n", message.Text))
	if err != nil {
		fmt.Println("Error writing to buffer")
		panic(err)
	}
	err = p.readWriter.Flush()
	if err != nil {
		fmt.Println("Error flushing buffer")
		panic(err)
	}
}

func (p *P2PNetworkProvider) Run() {
	host, err := libp2p.New(libp2p.ListenAddrs([]maddr.Multiaddr(p.config.ListenAddresses)...))
	if err != nil {
		panic(err)
	}
	p.logger.Info("Host created. We are:", host.ID())
	p.logger.Info(host.Addrs())

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(p.config.ProtocolID), p.handleStream)

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	ctx := context.Background()
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	p.logger.Debug("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range p.config.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				p.logger.Warning(err)
			} else {
				p.logger.Info("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	p.logger.Info("Announcing ourselves...")
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, p.config.RendezvousString)
	p.logger.Debug("Successfully announced!")

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	p.logger.Debug("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, p.config.RendezvousString)
	if err != nil {
		panic(err)
	}

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		p.logger.Debug("Found peer:", peer)

		p.logger.Debug("Connecting to:", peer)
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(p.config.ProtocolID))

		if err != nil {
			p.logger.Warning("Connection failed:", err)
			continue
		} else {
			p.readWriter = bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			go p.readData()
		}

		p.logger.Info("Connected to:", peer)
	}

	select {}
}

func (p *P2PNetworkProvider) handleStream(stream network.Stream) {
	p.logger.Info("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	p.readWriter = bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go p.readData()
}

func (p *P2PNetworkProvider) readData() {
	for {
		str, err := p.readWriter.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			p.messagesChannel <- entities.Message{Text: str}
		}

	}
}
