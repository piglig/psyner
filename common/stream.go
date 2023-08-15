package common

import (
	"encoding/binary"
	"net"
	"sync/atomic"
)

type Stream struct {
	Incoming chan<- *Packet
	Outgoing <-chan *Packet
	in       chan *Packet
	out      chan *Packet
	value    atomic.Value
	OnError  func(err error)
}

func NewStream(bufferSize int) *Stream {
	s := &Stream{
		in:      make(chan *Packet, bufferSize),
		out:     make(chan *Packet, bufferSize),
		value:   atomic.Value{},
		OnError: nil,
	}

	s.Incoming = s.in
	s.Outgoing = s.out
	return s
}

func (s *Stream) SetConnection(conn net.Conn) {
	if conn == nil {
		panic("SetConnection is nil")
	}

	s.value.Store(conn)

	go s.read(conn)
	go s.write(conn)
}

func (s *Stream) Connection() net.Conn {
	return s.value.Load().(net.Conn)
}

func (s *Stream) read(conn net.Conn) {
	var op FileSyncOp

	var length int64
	for {
		err := binary.Read(conn, binary.BigEndian, &op)
		if err != nil {
			return
		}

		err = binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			return
		}

		data := make([]byte, length)
		readLen := 0
		n := 0
		for readLen < len(data) {
			n, err = conn.Read(data[readLen:])
			if err != nil {
				return
			}

			readLen += n
		}

		s.in <- NewPacket(op, data)
	}
}

func (s *Stream) write(conn net.Conn) {

}
