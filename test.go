package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	tnetwork "github.com/notassigned/tweb/network"
)

func TestNetwork() {
	priv1, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	priv2, pub2, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	ma1, _ := ma.NewMultiaddr("/ip4/127.0.0.1/udp/24440/quic")
	ma2, _ := ma.NewMultiaddr("/ip4/127.0.0.1/udp/24441/quic")
	h1, e1 := libp2p.New(
		libp2p.Identity(priv1),
		libp2p.ListenAddrs(ma1),
	)
	h2, e2 := libp2p.New(
		libp2p.Identity(priv2),
		libp2p.ListenAddrs(ma2),
	)

	if e1 != nil {
		fmt.Println(e1)
	}
	if e2 != nil {
		fmt.Println(e2)
	}

	var (
		stream1 *tnetwork.MessageStream
		stream2 *tnetwork.MessageStream
	)

	h2.SetStreamHandler("test", func(s network.Stream) {
		ms := tnetwork.NewMessageStream(s)
		stream2 = &ms
	})

	id2, _ := peer.IDFromPublicKey(pub2)
	h1.Peerstore().AddAddr(id2, ma2, 9999999999)
	str, e := h1.NewStream(context.Background(), id2, "test")
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	s1 := tnetwork.NewMessageStream(str)
	stream1 = &s1
	time.Sleep(time.Second)
	testMessageStream(*stream1, *stream2)
}

func testMessageStream(stream1 tnetwork.MessageStream, stream2 tnetwork.MessageStream) {
	sent := make([]byte, 0)

	stream1.Write(sent)
	recv, err := stream2.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	assertAreEqual(sent, recv)

	stream2.Write(sent)
	recv, err = stream1.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	assertAreEqual(sent, recv)

}

func assertAreEqual(a []byte, b []byte) {
	if len(a) != len(b) {
		fmt.Println(len(a), len(b))
		os.Exit(1)
	}

	for x := 0; x < len(a); x++ {
		if a[x] != b[x] {
			fmt.Println(x, a[x], b[x])
			os.Exit(1)
		}
	}
}
