package client

import (
	"github.com/notassigned/tweb/xmlnode"
	tview "github.com/rivo/tview"
)

type Button struct {
	client *Client
	button *tview.Button
}

func (client *Client) CreateButton(n *xmlnode.XmlNode, events chan *xmlnode.XmlNode) (update func(*xmlnode.XmlNode), node func() *xmlnode.XmlNode, p tview.Primitive) {
	label := n.Attributes["label"]

	if label == "" {
		return
	}

	button := &Button{
		client: client,
		button: tview.NewButton(label),
	}

	button.button.SetSelectedFunc(func() {
		events <- xmlnode.New("selected").SetAttr("id", n.Attributes["id"])
	})

	button.Update(n)

	return button.Update, button.Node, button.button
}

func (b *Button) Node() *xmlnode.XmlNode {
	return nil
}

func (b *Button) Update(n *xmlnode.XmlNode) {
	for a, p := range n.Attributes {
		switch a {
		case "border":
			if p == "true" {
				b.button.SetBorder(true)
			} else {
				b.button.SetBorder(false)
			}
		}
	}
}
