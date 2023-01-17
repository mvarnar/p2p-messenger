package network

import (
	"math/rand"
	entities "p2p-messenger/src/domain/entities"
	"strconv"
	"time"
)

type P2PNetworkProvider struct {
}

func (p *P2PNetworkProvider) GetNewIncomingMessages() <-chan entities.Message {
	out := make(chan entities.Message)
	go func() {
		for {
			time.Sleep(3 * time.Second)

			out <- entities.Message{Text: strconv.Itoa(rand.Int())}
		}
	}()
	return out
}

func (p *P2PNetworkProvider) SendMessage(message entities.Message) {

}

func (p *P2PNetworkProvider) Run() {

}
