package examples

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	tnetwork "github.com/notassigned/tweb/network"
	x "github.com/notassigned/tweb/xmlnode"
)

const LOBBY_FILE = "./lobby.xml"

var view []byte

type nameCheck struct {
	name      string
	available chan bool
}

type Lobby struct {
	peers   map[string]tnetwork.Peer
	all     chan string
	joined  chan *user
	left    chan *user
	getName chan nameCheck
}

type user struct {
	name string
	peer tnetwork.Peer
}

func NewLobby() *Lobby {
	l := &Lobby{
		peers:   make(map[string]tnetwork.Peer),
		all:     make(chan string, 1),
		joined:  make(chan *user),
		left:    make(chan *user),
		getName: make(chan nameCheck),
	}
	go func() {
		for {
			select {
			case msg, ok := <-l.all:
				if !ok {
					return
				}
				addText := x.New("update").SetAttr("id", "chat").SetAttr("write", msg+"\n").Marshal()
				for _, p := range l.peers {
					p.Stream.MessageStream().Write(addText)
				}
			case u := <-l.joined:
				delete(l.peers, u.name)
				l.sendAddPeer(u.name)
				l.peers[u.name] = u.peer
				l.messageAll(fmt.Sprintf("[green]%s [blue]joined.[white]", u.name))
				l.sendPeerList(&u.peer)
				l.updateCount(len(l.peers))
				fmt.Println("joined")
			case u := <-l.left:
				delete(l.peers, u.name)
				l.messageAll(fmt.Sprintf("[green]%s [red]left.[white]", u.name))
				l.sendRemovePeer(u.name)
				l.updateCount(len(l.peers))
			case x := <-l.getName:
				_, exists := l.peers[x.name]
				if !exists {
					l.peers[x.name] = tnetwork.Peer{}
				}
				x.available <- !exists
			}
		}
	}()

	return l
}

func (l *Lobby) OnJoin(p tnetwork.Peer) {
	getLobbyView()

	name, e := getName(l, p.Stream)
	if e != nil {
		fmt.Println(e)
		return
	}

	p.Stream.MessageStream().Write(view)
	p.Stream.Write(x.Node("update", map[string]string{"id": "input", "label": name + ">"}))

	u := &user{name: name, peer: p}
	l.joined <- u

	for {
		n, e := p.Stream.Read()
		if e != nil {
			l.left <- u
			return
		}
		if n.Attributes["text"] != "" {
			p.Stream.Write(x.New("update").SetAttr("id", "input").SetAttr("text", ""))
			l.messageAll(fmt.Sprintf("[green]%s:[white] %s", name, n.Attributes["text"]))
		}
	}
}

func getName(l *Lobby, ms tnetwork.NodeStream) (name string, e error) {
	getName := x.New("rows").AddChild(x.Node("text", map[string]string{
		"text":  "Lobby",
		"color": "lightgreen",
		"align": "center",
		"size":  "1",
	})).AddChild(x.Node("input", map[string]string{
		"id":          "name",
		"label":       "Enter name: ",
		"labelcolor":  "lightgreen",
		"focus":       "true",
		"size":        "3",
		"border":      "true",
		"bordercolor": "red",
	}))
	ms.Write(x.New("newpage").AddChild(getName))

	for {
		resp, e := ms.Read()
		if e != nil {
			return "", e
		}
		if name = resp.Attributes["text"]; name != "" {
			available := make(chan bool)
			l.getName <- nameCheck{name: name, available: available}
			if <-available {
				return resp.Attributes["text"], nil
			}
			ms.Write(x.Node("update", map[string]string{
				"id":          "name",
				"text":        "",
				"placeholder": fmt.Sprintf("Name %s is taken", name),
			}))
		}
	}
}

func getLobbyView() {
	getLobby := func() {
		node, err := x.CreateNodeFromFile(LOBBY_FILE)
		if err != nil {
			fmt.Println("Error reading lobby view file:\n", err)
			return
		} else {
			fmt.Println("lobby updated")
			view = x.New("newpage").AddChild(node).Marshal()
		}
	}

	getLobby()

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println("NewWatcher failed: ", err)
			return
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			defer close(done)

			for {
				select {
				case _, ok := <-watcher.Events:
					if !ok {
						return
					}
					getLobby()
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					fmt.Println("error:", err)
				}
			}

		}()
	}()
}

func (l *Lobby) sendPeerList(p *tnetwork.Peer) {
	update := x.Update("peers", nil)

	for name := range l.peers {
		update.AddChild(x.Node("item", map[string]string{
			"text": name,
		}))
	}

	p.Stream.Write(update)
}

func (l *Lobby) sendAddPeer(name string) {
	bytes := x.Update("peers", nil).AddChild(x.Node("item", map[string]string{
		"text": name,
	})).Marshal()
	for _, p := range l.peers {
		p.Stream.MessageStream().Write(bytes)
	}
}

func (l *Lobby) sendRemovePeer(name string) {
	bytes := x.Update("peers", map[string]string{
		"remove": name,
	}).Marshal()

	for _, p := range l.peers {
		p.Stream.MessageStream().Write(bytes)
	}
}

func (l *Lobby) messageAll(msg string) {
	fmt.Println(msg)
	l.all <- msg
}

func (l *Lobby) updateCount(count int) {
	update := x.Node("update", map[string]string{"id": "count", "text": fmt.Sprintf("%d", count)}).Marshal()
	for _, p := range l.peers {
		p.Stream.MessageStream().Write(update)
	}
	fmt.Println("count", count)
}
