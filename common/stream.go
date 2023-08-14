package common

import "sync/atomic"

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
