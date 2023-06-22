package client

import (
	"context"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	tnetwork "github.com/notassigned/tweb/network"
	"github.com/notassigned/tweb/xmlnode"
	tview "github.com/rivo/tview"
)

type Libp2p struct {
	client     *Client
	base       *tview.Flex
	ctx        context.Context
	remotePeer peer.ID
	debugText  *tview.TextView
	debug      bool
	root       *tview.Pages
	close      func()
	children   map[string]func(*xmlnode.XmlNode)
	libp2pNode host.Host
	events     chan *xmlnode.XmlNode //events to send to server
}

func (client *Client) CreateLibp2p(n *xmlnode.XmlNode, ctx context.Context) (update func(*xmlnode.XmlNode), node func() *xmlnode.XmlNode, p tview.Primitive) {
	l := &Libp2p{
		client:    client,
		root:      tview.NewPages(),
		debugText: tview.NewTextView().SetText("start").SetDynamicColors(true).SetMaxLines(200),
		events:    make(chan *xmlnode.XmlNode),
		children:  map[string]func(*xmlnode.XmlNode){},
		base:      tview.NewFlex().SetDirection(tview.FlexRow),
		ctx:       ctx,
		debug:     client.debug,
	}

	l.base.AddItem(l.root, 0, 1, false)

	if l.debug {
		l.base.AddItem(l.debugText, 0, 1, false)
	}

	p = l.base

	node = func() *xmlnode.XmlNode { return n }

	priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)

	id, ok := n.Attributes["id"]
	if !ok {
		return nil, node, l.base
	}

	l.libp2pNode, _ = libp2p.New(
		libp2p.Identity(priv),
	)

	pid, err := peer.Decode(id)
	if err != nil {
		l.Error(err.Error())
		return nil, node, l.base
	}

	l.remotePeer = pid

	if addrs, ok := n.Attributes["addr"]; ok {
		//add comma separated multiaddrs to peerstore
		l.libp2pNode.Peerstore().AddAddrs(pid, func() []ma.Multiaddr {
			var split []ma.Multiaddr
			for _, v := range strings.Split(addrs, ",") {
				m, err := ma.NewMultiaddr(v)
				if err != nil {
					l.Error("Error parsing multiaddr: " + err.Error())
				} else {
					split = append(split, m)
				}
			}
			return split
		}(), time.Duration(int64(^uint64(0)>>1)))
	}

	go l.connect(n.Attributes["path"], pid)

	return nil, node, l.base
}

func (l *Libp2p) connect(path string, pid peer.ID) {
	stream, err := l.libp2pNode.NewStream(l.client.ctx, pid, protocol.ID("/tweb/1.0.0/"+path))

	if err != nil {
		l.Error("Connection failed: " + err.Error())
		return
	}

	nstream := tnetwork.CreatNodeStream(stream)

	//start routine to send events to remote
	go l.sendEventsToRemote(nstream)

	l.listen(nstream)
}

func (l *Libp2p) listen(stream tnetwork.NodeStream) {
	var err error = nil
	type read struct {
		node *xmlnode.XmlNode
		e    error
	}
	update := make(chan read)

	go func() {
		for {
			select {
			case <-l.ctx.Done():
				return
			default:
				node, e := stream.Read()
				update <- read{
					node: node,
					e:    e,
				}
			}
		}
	}()

	for err == nil {
		select {
		case <-l.ctx.Done():
			return
		case r := <-update:
			if r.e != nil {
				l.Error(r.e.Error())
				return
			}
			l.Update(r.node)
		}
	}
}

func (l *Libp2p) sendEventsToRemote(stream tnetwork.NodeStream) {
	for {
		select {
		case <-l.ctx.Done():
			return
		case x := <-l.events:
			if l.debug {
				l.Error("[blue]sent: [white]" + string(x.Marshal()))
			}
			err := stream.Write(x)
			if err != nil {
				l.Error("[red]Disconnected:[white] " + err.Error())
				return
			}
		}

	}
}

func (l *Libp2p) Error(e string) {
	if l.debug {
		l.debugText.Write(append([]byte(e), []byte("\n")...))
	}
}

func (l *Libp2p) Node() *xmlnode.XmlNode {
	return nil
}

func (l *Libp2p) control(node *xmlnode.XmlNode) {
	command := node.Attributes["command"]

	switch command {
	case "clear":
		l.client.tviewApp.QueueUpdateDraw(func() {
			l.root.AddPage("main", tview.NewFlex(), true, true)
		})
	}
}

func (l *Libp2p) Update(node *xmlnode.XmlNode) {
	debugText := node.Marshal()
	l.Error("[gold]received:[white] " + string(debugText))
	switch node.XMLName.Local {
	case "control":
		l.control(node)
	case "newpage":
		if l.close != nil {
			l.close()
		}
		l.client.tviewApp.QueueUpdateDraw(func() {
			l.children = make(map[string]func(*xmlnode.XmlNode))
			ctx, cancel := context.WithCancel(l.ctx)
			l.close = cancel
			update, _, p, _ := l.client.CreateItem(node.Nodes[0], l.events, false, l.children, ctx)

			if id := node.Attributes["id"]; id != "" {
				l.children[id] = update
			}

			l.root.AddPage("main", p, true, true)
		})
	case "update":
		id := node.Attributes["id"]

		if id == "" {
			return
		}

		update := l.children[id]
		if update != nil {
			l.client.tviewApp.QueueUpdateDraw(func() {
				update(node)
			})
		}

	default:
		l.Error("[red]unknown node in update:[white] " + node.XMLName.Local)
	}
}

func (l *Libp2p) Close() {
	if l.close != nil {
		l.close()
	}
}
