package client

import (
	"strconv"

	"github.com/notassigned/tweb/xmlnode"
	tview "github.com/rivo/tview"
)

type Text struct {
	client *Client
	node   *xmlnode.XmlNode
	text   *tview.TextView
}

func (client *Client) CreateText(n *xmlnode.XmlNode) (update func(*xmlnode.XmlNode), node func() *xmlnode.XmlNode, p tview.Primitive) {
	t := &Text{}
	t.client = client
	t.text = tview.NewTextView()
	t.Update(n)

	return t.Update, t.Node, t.text
}

func (t *Text) Update(node *xmlnode.XmlNode) {
	if len(node.Content) > 0 {
		t.text.SetText(string(node.Content))
	}

	for a, p := range node.Attributes {
		switch a {
		case "dynamiccolors":
			if p == "true" {
				t.text.SetDynamicColors(true)
			} else {
				t.text.SetDynamicColors(false)
			}
		case "text":
			t.text.SetText(p)
		case "write":
			t.text.Write([]byte(p))
		case "maxlines":
			if max, e := strconv.Atoi(p); e == nil {
				t.text.SetMaxLines(max)
			}
		case "wrap":
			if p == "true" {
				t.text.SetWrap(true)
			} else {
				t.text.SetWrap(false)
			}
		case "wordwrap":
			if p == "true" {
				t.text.SetWordWrap(true)
			} else {
				t.text.SetWordWrap(false)
			}
		case "align":
			t.text.SetTextAlign(alignLookup(p))
		case "color":
			if color, err := colorLookup(p); err == nil {
				t.text.SetTextColor(color)
			}
		case "border":
			if p == "true" {
				t.text.SetBorder(true)
			} else {
				t.text.SetBorder(false)
			}
		case "bordercolor":
			bcolor, err := colorLookup(p)
			if err == nil {
				t.text.SetBorderColor(bcolor)
			}
		case "title":
			t.text.SetTitle(p)
		case "titlecolor":
			if color, err := colorLookup(p); err == nil {
				t.text.SetTitleColor(color)
			}
		case "titlealign":
			t.text.SetTitleAlign(alignLookup(p))
		case "content":
			t.text.SetText(p)
		}
	}
}

func (t *Text) Primitive() tview.Primitive {
	return t.text
}

func (t *Text) Node() *xmlnode.XmlNode {
	return t.node
}
