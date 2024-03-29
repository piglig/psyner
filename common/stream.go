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
	done     chan struct{}
	onError  func(err error)
}

func NewStream(bufferSize int) *Stream {
	s := &Stream{
		in:      make(chan *Packet, bufferSize),
		out:     make(chan *Packet, bufferSize),
		done:    make(chan struct{}),
		value:   atomic.Value{},
		onError: nil,
	}

	s.Incoming = s.in
	s.Outgoing = s.out
	return s
}

func (s *Stream) OnError(callback func(error)) {
	if callback == nil {
		panic("OnError is nil")
	}

	s.onError = callback
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
	var err error
	var op FileSyncOp
	var length uint64

	defer func() {
		if err != nil {
			s.onError(err)
			s.done <- struct{}{}
		}
	}()

	for {
		err = binary.Read(conn, binary.BigEndian, &op)
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
	for {
		select {
		case p := <-s.out:
			err := p.Write(conn)
			if err != nil {
				return
			}
		case <-s.done:
			return
		}
	}
}

func (s *Stream) Close() {
	s.Connection().Close()
	close(s.in)
}
