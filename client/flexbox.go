package client

import (
	"context"

	"github.com/notassigned/tweb/xmlnode"
	tview "github.com/rivo/tview"
)

type FlexBoxItem struct {
	flex       *FlexBox
	update     func(*xmlnode.XmlNode)
	node       func() *xmlnode.XmlNode
	primitive  tview.Primitive
	size       int
	proportion int
	focus      bool
}

type FlexBox struct {
	client      *Client
	parent      *FlexBox
	flex        *tview.Flex
	items       []*FlexBoxItem
	children    map[string]func(*xmlnode.XmlNode)
	allowRemote bool
}

func (f *FlexBox) Update(n *xmlnode.XmlNode) {
	if p, ok := n.Attributes["orient"]; ok {
		switch p {
		case "rows":
			f.flex.SetDirection(tview.FlexRow)
		case "cols":
			f.flex.SetDirection(tview.FlexColumn)
		}
	} else {
		f.flex.SetDirection(tview.FlexRow)
	}

	for a, p := range n.Attributes {
		switch a {
		case "border":
			if p == "true" {
				f.flex.SetBorder(true)
			} else {
				f.flex.SetBorder(false)
			}
		case "bordercolor":
			bcolor, err := colorLookup(p)
			if err == nil {
				f.flex.SetBorderColor(bcolor)
			}
		}
	}
}

func (client *Client) CreateFlexBox(n *xmlnode.XmlNode, allowRemote bool, events chan *xmlnode.XmlNode, children map[string]func(*xmlnode.XmlNode)) (update func(*xmlnode.XmlNode), node func() *xmlnode.XmlNode, p tview.Primitive, f *FlexBox) {
	if children == nil {
		children = make(map[string]func(*xmlnode.XmlNode))
	}

	f = &FlexBox{
		client:      client,
		flex:        tview.NewFlex(),
		allowRemote: allowRemote,
		children:    children,
	}

	if n.XMLName.Local != "flex" {
		n.SetAttr("orient", n.XMLName.Local)
	}
	f.Update(n)

	for _, c := range n.Nodes {
		f.addElement(c, events, client.ctx)
	}

	return f.Update, f.Node, f.flex, f
}

func (f FlexBox) Node() *xmlnode.XmlNode {
	var n xmlnode.XmlNode
	n.XMLName.Local = "flex"

	return &n
}

func (f *FlexBox) clear() {
	for _, c := range f.items {
		if c.flex != nil {
			c.flex.clear()
		}
	}
	f.flex.Clear()
	f.items = nil
	f.children = map[string]func(*xmlnode.XmlNode){}
}

func (f *FlexBox) addElement(n *xmlnode.XmlNode, events chan *xmlnode.XmlNode, ctx context.Context) {
	var (
		item      *FlexBoxItem = nil
		flex      *FlexBox     = nil
		update    func(*xmlnode.XmlNode)
		node      func() *xmlnode.XmlNode
		primitive tview.Primitive
	)

	size, proportion, focus := SizeProportionFocus(n)

	update, node, primitive, flex = f.client.CreateItem(n, events, f.allowRemote, f.children, ctx)
	if primitive == nil {
		return
	}

	item = makeFlexBoxItem(update, node, size, proportion, focus)
	item.flex = flex

	if id := n.Attributes["id"]; id != "" {
		f.children[id] = update
	}

	f.items = append(f.items, item)

	f.flex.AddItem(primitive, size, proportion, focus)
	if focus {
		primitive.Focus(nil)
	}
}

func makeFlexBoxItem(
	update func(*xmlnode.XmlNode),
	node func() *xmlnode.XmlNode,
	size int,
	proportion int,
	focus bool) *FlexBoxItem {
	item := &FlexBoxItem{
		update:     update,
		node:       node,
		size:       size,
		proportion: proportion,
		focus:      focus,
	}
	return item
}
