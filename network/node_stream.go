package tnetwork

import (
	"errors"
	"io"

	"github.com/notassigned/tweb/xmlnode"
)

type NodeStream struct {
	Read          func() (*xmlnode.XmlNode, error)
	Write         func(*xmlnode.XmlNode) error
	Close         func() error
	MessageStream func() MessageStream
}

func NewNodeStream(s io.ReadWriteCloser) NodeStream {
	msgStream := NewMessageStream(s)
	return NodeStream{
		Read: func() (*xmlnode.XmlNode, error) {
			bytes, err := msgStream.Read()
			if err != nil {
				return nil, err
			}
			return xmlnode.CreateNodeFromXML(string(bytes))
		},
		Write: func(xn *xmlnode.XmlNode) error {
			if xn == nil {
				return errors.New("node was nil")
			}
			return msgStream.Write(xn.Marshal())
		},
		Close:         msgStream.Close,
		MessageStream: func() MessageStream { return msgStream },
	}
}
