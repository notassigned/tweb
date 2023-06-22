package client

import (
	"github.com/notassigned/tweb/xmlnode"
	tview "github.com/rivo/tview"
)

type Form struct {
}

func (client *Client) CreateForm(n *xmlnode.XmlNode) (update func(*xmlnode.XmlNode), node func() *xmlnode.XmlNode, p tview.Primitive) {
	//f := tview.NewForm()
	// var formResults []string
	// for _, c := range n.Nodes {
	// 	switch c.XMLName.Local {
	// 	case "input":
	// 		f.AddInputField(c.Attributes["label"], c.Attributes["value"],
	// 			intOrZero(c.Attributes["width"]), func(textToCheck string, lastChar rune) bool {

	// 			}, )
	// 	}
	// }
	return nil, nil, nil
}
