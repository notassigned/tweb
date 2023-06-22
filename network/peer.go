package tnetwork

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Peer struct {
	ID     peer.ID
	Stream NodeStream
}

func NewPeer(s network.Stream) *Peer {
	return &Peer{
		ID:     s.Conn().RemotePeer(),
		Stream: CreatNodeStream(s),
	}
}
