package server

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	tnetwork "github.com/notassigned/tweb/network"
)

type Server struct {
	host  host.Host
	Close func() error
}

func New(h host.Host) *Server {
	return &Server{
		host:  h,
		Close: h.Close,
	}
}

func (s *Server) SetHandler(path string, onConnect func(tnetwork.Peer)) {
	s.host.SetStreamHandler(protocol.ID(fmt.Sprintf("/tweb/1.0.0/%s", path)), func(s network.Stream) {
		onConnect(*tnetwork.NewPeer(s))
	})
}
