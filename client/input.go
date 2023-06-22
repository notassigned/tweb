package client

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/notassigned/tweb/xmlnode"
	tview "github.com/rivo/tview"
)

type InputField struct {
	client     *Client
	inputField *tview.InputField
	events     chan *xmlnode.XmlNode
}

func (client *Client) CreateInputField(n *xmlnode.XmlNode, events chan *xmlnode.XmlNode) (update func(*xmlnode.XmlNode), node func() *xmlnode.XmlNode, p tview.Primitive) {
	input := &InputField{
		client:     client,
		inputField: tview.NewInputField(),
		events:     events,
	}

	input.updateInputField(n)
	input.inputField.SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.NewRGBColor(0, 0, 0)))
	input.inputField.SetFieldStyle(tcell.StyleDefault.Background(tcell.NewRGBColor(15, 15, 15)))

	return input.updateInputField, node, input.inputField
}

func (i *InputField) updateInputField(n *xmlnode.XmlNode) {
	input := i.inputField

	input.SetDoneFunc(func(key tcell.Key) {
		if i.events != nil {
			i.events <- xmlnode.Node("event", map[string]string{
				"id":   n.Attributes["id"],
				"text": input.GetText(),
				"key":  tcell.KeyNames[key],
			})
		}
	})

	for a, p := range n.Attributes {
		switch a {
		case "color":
			color, err := colorLookup(p)
			if err == nil {
				input.SetFieldTextColor(color)
			}
		case "background":
			color, err := colorLookup(p)
			if err == nil {
				input.SetFieldBackgroundColor(color)
			} else {
				input.SetFieldBackgroundColor(tcell.ColorBlack)
			}
		case "width":
			if i, err := strconv.Atoi(p); err == nil {
				input.SetFieldWidth(i)
			}
		case "text":
			input.SetText(p)
		case "label":
			input.SetLabel(p)
		case "labelcolor":
			color, err := colorLookup(p)
			if err == nil {
				input.SetLabelColor(color)
			}
		case "title":
			input.SetTitle(p)
		case "titleAlign":
			input.SetTitleAlign(alignLookup(p))
		case "titlecolor":
			color, err := colorLookup(p)
			if err == nil {
				input.SetTitleColor(color)
			}
		case "titlealign":
			input.SetTitleAlign(alignLookup(p))
		case "placeholder":
			input.SetPlaceholder(p)
		case "border":
			if p == "true" {
				input.SetBorder(true)
			}
		case "bordercolor":
			color, err := colorLookup(p)
			if err == nil {
				input.SetBorderColor(color)
			}
		}
	}
}
