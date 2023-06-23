package xmlnode

import (
	"encoding/xml"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
)

type XmlNode struct {
	XMLName    xml.Name
	Attrs      []xml.Attr        `xml:",any,attr"`
	Content    []byte            `xml:",innerxml"`
	Nodes      []*XmlNode        `xml:",any"`
	Attributes map[string]string `xml:"-"`
}

// Create an XmlNode with name
func New(name string) *XmlNode {
	n := &XmlNode{}
	n.SetName(name)
	return n
}

// Create an XmlNode with name and attributes
func Node(name string, attributes map[string]string) *XmlNode {
	n := &XmlNode{}
	n.SetName(name)
	n.SetAttrs(attributes)
	return n
}

// Return XML representing a XmlNode struct
func (n *XmlNode) Marshal() []byte {
	if n.Attributes != nil {
		n.Attrs = make([]xml.Attr, len(n.Attributes))
		for name, value := range n.Attributes {
			x := xml.Attr{Name: xml.Name{Local: name}, Value: value}
			n.Attrs = append(n.Attrs, x)
		}
	}

	return []byte(n.getXMLQueryNode().OutputXML(true))
}

// Create a xmlquery.Node representing a XmlNode struct (for marshaling)
func (n *XmlNode) getXMLQueryNode() *xmlquery.Node {
	q := &xmlquery.Node{
		Data: n.XMLName.Local,
		Type: xmlquery.ElementNode,
	}

	content := string(n.Content)
	if len(content) > 0 {
		q.FirstChild = &xmlquery.Node{Type: xmlquery.TextNode, Data: content}
	}

	for a, v := range n.Attributes {
		if a != "" {
			q.Attr = append(q.Attr, xmlquery.Attr{Name: xml.Name{Local: a}, Value: v})
		}
	}

	for i := len(n.Nodes); i > 0; i-- {
		child := n.Nodes[i-1]
		cq := child.getXMLQueryNode()
		if q.FirstChild != nil {
			q.FirstChild.PrevSibling = cq
		}
		cq.NextSibling = q.FirstChild
		q.FirstChild = cq
	}

	return q
}

func CreateNodeFromFile(path string) (*XmlNode, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return CreateNodeFromXML(string(content))
}

func CreateNodeFromXML(xmlStr string) (*XmlNode, error) {
	node := &XmlNode{}

	n, err := xmlquery.Parse(strings.NewReader(xmlStr))
	if err != nil {
		return nil, err
	}

	n.Attr = nil

	err = xml.NewDecoder(strings.NewReader(xmlStr)).Decode(&node)
	if err == nil {
		node.MapAttributesRecur()
	}

	return node, err
}

func (n *XmlNode) AddChild(newNode *XmlNode) *XmlNode {
	n.Nodes = append(n.Nodes, newNode)
	return n
}

func (n *XmlNode) SetName(name string) *XmlNode {
	n.XMLName.Local = name
	return n
}

func (n *XmlNode) SetContent(s string) *XmlNode {
	n.Content = []byte(s)
	return n
}

func (n *XmlNode) SetAttrs(a map[string]string) *XmlNode {
	for attr, p := range a {
		n.SetAttr(attr, p)
	}
	return n
}

func (n *XmlNode) SetAttr(name string, value string) *XmlNode {
	if n.Attributes == nil {
		n.Attributes = map[string]string{}
		if n.Attrs != nil {
			for _, a := range n.Attrs {
				n.Attributes[a.Name.Local] = a.Value
			}
		}
	}

	n.Attributes[name] = value

	return n
}

func (n *XmlNode) MapAttributesRecur() {
	n.Attributes = make(map[string]string)
	if n.Attrs != nil {
		for _, b := range n.Attrs {
			n.Attributes[b.Name.Local] = b.Value
		}
	}
	for _, n := range n.Nodes {
		n.MapAttributesRecur()
	}
}
