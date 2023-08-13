package common

import (
	"encoding/binary"
	"io"
)

type Packet struct {
	Op      FileSyncOp
	Length  uint64
	Payload []byte
}

func NewPacket(op FileSyncOp, payload []byte) *Packet {
	return &Packet{Op: op, Length: uint64(len(payload)), Payload: payload}
}

func (p *Packet) Write(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, p.Op)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, p.Length)
	if err != nil {
		return err
	}

	n := 0
	writtenLen := 0
	writeUtil := 0
	for writtenLen < len(p.Payload) {
		writeUtil = writtenLen + DataBufferSize
		if writeUtil > DataBufferSize {
			writeUtil = DataBufferSize
		}

		n, err = w.Write(p.Payload[writtenLen:writeUtil])
		if err != nil {
			return err
		}

		writtenLen += n
	}
	return
}
