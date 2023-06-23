package client

import (
	"context"
	"encoding/xml"
	"errors"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/notassigned/tweb/xmlnode"
	"github.com/rivo/tview"
)

type Client struct {
	tviewApp *tview.Application
	elements map[string]func(*xmlnode.XmlNode) //element ids to element update
	system   *tview.Pages
	ctx      context.Context
	cancel   context.CancelFunc
	debug    bool
}

func CreateClient(source string) (c *Client, err error) {
	context, cancel := context.WithCancel(context.Background())
	c = &Client{
		elements: make(map[string]func(*xmlnode.XmlNode)),
		tviewApp: tview.NewApplication(),
		system:   tview.NewPages(),
		ctx:      context,
		cancel:   cancel,
		debug:    false,
	}

	tview.Styles.PrimitiveBackgroundColor = tcell.NewRGBColor(0, 0, 0)
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.SecondaryTextColor = tcell.ColorWhite

	c.tviewApp.SetRoot(c.system, true)
	c.tviewApp.EnableMouse(true)

	if len(source) < 1 {
		return nil, errors.New("source empty")
	}

	c.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC: // Exit the application
			c.tviewApp.Stop()
			return nil
		case tcell.KeyF1: // Switch to the system page
			if !c.debug {
				c.system.SwitchToPage("system")
				c.debug = true
			} else {
				c.system.SwitchToPage("main")
				c.debug = false
			}
		}
		return event
	})

	c.setupSystemPage()

	c.parseXML(source, context)

	return c, nil
}

func (c *Client) Start() {
	defer c.cancel()
	c.tviewApp.Run()
}

func (c *Client) CreateItem(
	n *xmlnode.XmlNode,
	events chan *xmlnode.XmlNode,
	allowRemote bool,
	children map[string]func(*xmlnode.XmlNode),
	ctx context.Context,
) (
	update func(*xmlnode.XmlNode),
	node func() *xmlnode.XmlNode,
	p tview.Primitive,
	f *FlexBox,
) {
	switch t := n.XMLName.Local; t {
	case "flex", "rows", "cols":
		return c.CreateFlexBox(n, allowRemote, events, children)
	case "text":
		update, node, p = c.CreateText(n)
	case "input":
		update, node, p = c.CreateInputField(n, events)
	case "list":
		update, node, p = c.CreateListView(n, events)
	case "button":
		update, node, p = c.CreateButton(n, events)
	case "libp2p":
		update, node, p = c.CreateLibp2p(n, ctx)
	default:
		return nil, nil, nil, nil
	}
	return update, node, p, f
}

func (c *Client) parseXML(s string, ctx context.Context) error {
	var node xmlnode.XmlNode
	d := xml.NewDecoder(strings.NewReader(s))

	root := makeRoot(c, true, nil)

	err := d.Decode(&node)
	if err != nil {
		return err
	}
	node.MapAttributesRecur()
	root.addElement(&node, nil, ctx)

	c.setRoot(root.flex)
	return nil
}

func (ui *Client) setRoot(p tview.Primitive) {
	ui.system.AddPage("main", p, true, true)
}

func (c *Client) setupSystemPage() {
	flex := tview.NewFlex().AddItem(tview.NewList().AddItem("Enable debug", "", 0, func() {
		c.Debug(true)
	}).AddItem("Disable debug", "", 0, func() {
		c.Debug(false)
	}), 0, 1, true)
	c.system.AddPage("system", flex, true, true)
}

func makeRoot(c *Client, allowRemote bool, events chan *xmlnode.XmlNode) *FlexBox {
	t := "<rows></rows>"
	node, _ := xmlnode.CreateNodeFromXML(t)
	_, _, _, f := c.CreateFlexBox(node, allowRemote, events, nil)
	return f
}

func (c *Client) Debug(state bool) {
	//not implemented
}
