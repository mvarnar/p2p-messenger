package network

import (
	"bufio"
	"context"
	"fmt"
	entities "p2p-messenger/src/domain/entities"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"

	"github.com/libp2p/go-libp2p/core/protocol"
 
	"encoding/json"

	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	maddr "github.com/multiformats/go-multiaddr"
)

type P2PNetworkProvider struct {
	config          Config
	messagesChannel chan entities.Message
	host            host.Host
	readWriters     map[entities.Contact]*bufio.ReadWriter
	kademliaDht     *dht.IpfsDHT
	keyBytes        []byte
}

func NewP2PNetworkProvider(Config Config) *P2PNetworkProvider {
	return &P2PNetworkProvider{
		config:          Config,
		messagesChannel: make(chan entities.Message, 100),
		readWriters:     make(map[entities.Contact]*bufio.ReadWriter),
	}
}

func (p *P2PNetworkProvider) GetNewIncomingMessages() <-chan entities.Message {
	return p.messagesChannel
}

func (p *P2PNetworkProvider) SendMessage(message entities.Message) {
	readWriter, ok := p.readWriters[message.ReceiverContact]
	if !ok {
		pid, err := peer.Decode(message.ReceiverContact.UserId)
		if err != nil {
			panic(err)
		}

		var streamErr error
		for try := 0; try < 5; try++ {
			receiverAddrInfo, _ := p.kademliaDht.FindPeer(context.Background(), pid)
			fmt.Printf("Sending message to %s\n", receiverAddrInfo)
			stream, streamErr := p.host.NewStream(context.Background(), pid, protocol.ID(p.config.ProtocolID))

			if streamErr != nil {
				fmt.Println(streamErr)
			} else {
				readWriter = bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
				go p.readData(readWriter)
				p.readWriters[message.ReceiverContact] = readWriter
				break
			}

			time.Sleep(1000 * time.Millisecond)
		}
		if streamErr != nil {
			panic(err)
		}
	}

	jsonMessage, _ := json.Marshal(message)
	_, err := readWriter.WriteString(string(jsonMessage) + "\n")
	if err != nil {
		fmt.Println("Error writing to buffer")
		panic(err)
	}
	err = readWriter.Flush()
	if err != nil {
		fmt.Println("Error flushing buffer")
		panic(err)
	}
}

func (p *P2PNetworkProvider) GetUserId() string {
	for p.host == nil {
		time.Sleep(1 * time.Second)
	}
	fmt.Println(p.host.ID())
	return p.host.ID().String()
}

func (p *P2PNetworkProvider) Run(keyBytes []byte) {
	var key crypto.PrivKey
	var err error
	if keyBytes == nil {
		key, _, err = crypto.GenerateKeyPair(crypto.RSA, 2048)
	} else {
		key, err = crypto.UnmarshalPrivateKey(keyBytes)
	}
	if err != nil {
		panic(err)
	}

	p.keyBytes, err = crypto.MarshalPrivateKey(key)
	if err != nil {
		panic(err)
	}
	identify.ActivationThresh = 1
	p.host, err = libp2p.New(
		libp2p.ListenAddrs([]maddr.Multiaddr(p.config.ListenAddresses)...),
		libp2p.Identity(key),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.EnableHolePunching(),
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.EnableAutoRelay(autorelay.WithStaticRelays(dht.GetDefaultBootstrapPeerAddrInfos())),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Host created. We are:", p.host.ID())

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	p.host.SetStreamHandler(protocol.ID(p.config.ProtocolID), p.handleStream)

	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	ctx := context.Background()
	p.kademliaDht, err = dht.New(ctx, p.host)
	if err != nil {
		panic(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	if err = p.kademliaDht.Bootstrap(ctx); err != nil {
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
			if err := p.host.Connect(ctx, *peerinfo); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	fmt.Println(p.host.Addrs())
}

func (p *P2PNetworkProvider) GetKeyBytes() []byte {
	for p.host == nil {
		time.Sleep(1 * time.Second)
	}
	return p.keyBytes
}

func (p *P2PNetworkProvider) handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	contact := entities.Contact{UserId: string(stream.Conn().RemotePeer())}
	p.readWriters[contact] = readWriter

	go p.readData(readWriter)
}

func (p *P2PNetworkProvider) readData(readWriter *bufio.ReadWriter) {
	for {
		bytes, err := readWriter.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}
		var message entities.Message
		json.Unmarshal(bytes, &message)
		p.messagesChannel <- message
	}
}
