package tnetwork

import (
	"encoding/binary"
	"errors"
	"io"
)

type MessageStream struct {
	Read  func() ([]byte, error)
	Write func([]byte) error
	Close func() error
}

func NewMessageStream(stream io.ReadWriteCloser) MessageStream {
	nreader := NewNReader(stream)
	incoming := make(chan []byte)
	outgoing := make(chan []byte)
	closed := false

	go func() {
		for {
			m, e := readMessage(nreader)
			if e != nil {
				close(incoming)
				return
			}
			incoming <- m
		}
	}()

	go func() {
		for {
			b, ok := <-outgoing
			if !ok {
				return
			}
			if e := writeMessage(stream, b); e != nil {
				closed = true
				return
			}
		}
	}()

	return MessageStream{
		Read: func() ([]byte, error) {
			m, ok := <-incoming
			if !ok {
				return nil, errors.New("closed")
			}
			return m, nil
		},
		Write: func(b []byte) error {
			if closed {
				return errors.New("closed")
			}
			outgoing <- b
			return nil
		},
		Close: func() error {
			closed = true
			return stream.Close()
		},
	}
}

func writeMessage(w io.Writer, message []byte) (err error) {
	sizeVarint := make([]byte, 10)

	n := binary.PutUvarint(sizeVarint, uint64(len(message)))
	if n < 1 {
		return errors.New("invalid size")
	}

	_, err = w.Write(sizeVarint[:n])
	if err != nil {
		return
	}

	_, err = w.Write(message)

	return
}
