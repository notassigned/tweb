package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/mr-tron/base58"

	"github.com/notassigned/tweb/examples"
	tnetwork "github.com/notassigned/tweb/network"
	"github.com/notassigned/tweb/server"
	"github.com/notassigned/tweb/xmlnode"
)

func StartServer() {
	fmt.Println("Starting server")
	path := os.Args[2]

	//read server id file and setup libp2p host
	content, e := ioutil.ReadFile(path)
	checkErr(e)
	decoded, _ := base58.Decode(string(content))
	priv, e := crypto.UnmarshalPrivateKey(decoded)
	checkErr(e)

	h, e := libp2p.New(
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/udp/7000/quic"),
	)
	checkErr(e)

	server1 := server.New(h)

	server1.SetHandler("timeserver", timeserver)
	server1.SetHandler("lobby", examples.NewLobby().OnJoin)

	for {
		time.Sleep(time.Second)
	}
}

func timeserver(p tnetwork.Peer) {
	spacer := xmlnode.New("text")
	text := xmlnode.New("text").SetAttr("size", "1").SetAttr("align", "center")
	flex := xmlnode.New("flex").AddChild(spacer).AddChild(text).AddChild(spacer)
	newpage := xmlnode.New("newpage").AddChild(flex)
	for {
		start := time.Now()
		text.SetAttr("text", time.Now().Format("02-Jan-2006 15:04:05"))
		if p.Stream.Write(newpage) != nil {
			return
		}

		time.Sleep(time.Second - time.Since(start))
	}
}

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
