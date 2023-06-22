package client

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/notassigned/tweb/xmlnode"
	tview "github.com/rivo/tview"
)

type ListView struct {
	client *Client
	list   *tview.List
	events chan *xmlnode.XmlNode
}

func (client *Client) CreateListView(n *xmlnode.XmlNode, events chan *xmlnode.XmlNode) (update func(*xmlnode.XmlNode), node func() *xmlnode.XmlNode, p tview.Primitive) {
	listView := &ListView{
		client: client,
		list:   tview.NewList(),
		events: events,
	}

	listView.list.SetSecondaryTextStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDimGray))
	listView.list.SetMainTextStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGray))
	listView.list.SetMainTextStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	listView.list.ShowSecondaryText(false)
	listView.list.SetSelectedFocusOnly(true)

	listView.updateListView(n)

	listView.list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		events <- xmlnode.Node("selection", map[string]string{
			"id":    n.Attributes["id"],
			"index": fmt.Sprint(i),
			"key":   tcell.KeyNames[tcell.Key(r)],
		})
	})

	return listView.updateListView, node, listView.list
}

func (l *ListView) updateListView(node *xmlnode.XmlNode) {
	list := l.list

	l.AddItems(node.Nodes)

	for a, p := range node.Attributes {
		switch a {
		case "removesecondary":
			indices := list.FindItems("", p, false, false)
			for _, i := range indices {
				if _, text := list.GetItemText(i); text == p {
					list.RemoveItem(i)
				}
			}
		case "remove":
			indices := list.FindItems(p, "", false, false)
			for _, i := range indices {
				if text, _ := list.GetItemText(i); text == p {
					list.RemoveItem(i)
				}
			}
		case "color":
			color, err := colorLookup(p)
			if err == nil {
				list.SetMainTextColor(color)
			}
		case "background":
			color, err := colorLookup(p)
			if err == nil {
				list.SetBackgroundColor(color)
			}
		case "selectedbackground":
			color, err := colorLookup(p)
			if err == nil {
				list.SetSelectedBackgroundColor(color)
			}
		case "selectedcolor":
			color, err := colorLookup(p)
			if err == nil {
				list.SetSelectedTextColor(color)
			}
		case "secondarycolor":
			color, err := colorLookup(p)
			if err == nil {
				list.SetSecondaryTextColor(color)
			}
		case "title":
			list.SetTitle(p)
		case "titlecolor":
			color, err := colorLookup(p)
			if err == nil {
				list.SetTitleColor(color)
			}
		case "titlealign":
			list.SetTitleAlign(alignLookup(p))
		case "border":
			if p == "true" {
				list.SetBorder(true)
			}
		case "bordercolor":
			bcolor, err := colorLookup(p)
			if err == nil {
				list.SetBorderColor(bcolor)
			}
		}
	}
}

func (l *ListView) ClearItems() {
	l.list.Clear()
}

func (l *ListView) AddItems(nodes []*xmlnode.XmlNode) {
	list := l.list
	var count = 0
	for _, n := range nodes {
		if n.XMLName.Local == "item" {
			var r rune
			text := string(n.Content)
			var secondaryText string = ""

			for a, p := range n.Attributes {
				switch a {
				case "text":
					text = p
				case "secondary":
					secondaryText = p
				case "showsecondary":
					if p == "true" {
						list.ShowSecondaryText(true)
					} else {
						list.ShowSecondaryText(false)
					}

				case "selected":
					if p == "true" {
						list.SetCurrentItem(count)
					}
				}
			}
			list.AddItem(text, secondaryText, r, nil)

			count++
		}
	}
}
