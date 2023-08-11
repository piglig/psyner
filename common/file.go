package common

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type FileSyncOp uint32

const (
	GetFileOp FileSyncOp = iota + 1
	UpdateFileOp
	DeleteFileOp
	DataOp
)

const (
	DataBufferSize = 1020
	BlockSize      = DataBufferSize - 8
)

type FileSyncPayload struct {
	OpType        FileSyncOp
	ActionPayload []byte
}

type GetFileOpPayload struct {
	RelPath string
}

func (d *GetFileOpPayload) MarshalBinary() ([]byte, error) {
	if d.RelPath == "" {
		return nil, errors.New("invalid path")
	}

	c := 4 + len(d.RelPath) + 1

	b := new(bytes.Buffer)
	b.Grow(c)

	err := binary.Write(b, binary.BigEndian, GetFileOp)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(d.RelPath)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *GetFileOpPayload) UnmarshalBinary(p []byte) error {
	if l := len(p); l < 5 || l > DataBufferSize {
		return errors.New("invalid GetFileOp len")
	}

	r := bytes.NewBuffer(p)

	var op FileSyncOp
	err := binary.Read(r, binary.BigEndian, &op)
	if err != nil {
		return err
	}

	if op != GetFileOp {
		return errors.New("invalid GetFileOp")
	}

	fileName, err := r.ReadString(0)
	if err != nil {
		return err
	}

	d.RelPath = strings.TrimRight(fileName, "\x00")
	return nil
}

type DataPayload struct {
	Block   uint32
	Payload io.Reader
}

func (d *DataPayload) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	b.Grow(DataBufferSize)

	d.Block++

	err := binary.Write(b, binary.BigEndian, DataOp)
	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(b, d.Payload, BlockSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *DataPayload) UnmarshalBinary(p []byte) error {
	if l := len(p); l < 8 || l > DataBufferSize {
		return errors.New("invalid Data")
	}

	var op FileSyncOp
	err := binary.Read(bytes.NewReader(p[:4]), binary.BigEndian, &op)
	if err != nil {
		return err
	}

	if op != DataOp {
		return errors.New("invalid DATA")
	}

	err = binary.Read(bytes.NewReader(p[4:8]), binary.BigEndian, &d.Block)
	if err != nil {
		return err
	}

	d.Payload = bytes.NewBuffer(p[8:])
	return nil
}

type GetFileSyncPayloadRes struct {
	RelPath  string
	FileSize int64
}

type UpdateFileSyncPayload struct {
	RelPath  string
	FileHash string
}

type DeleteFileSyncPayload struct {
	RelPath string
}

type FsWatcherCreateFilePayload struct {
	FileName string
	RelPath  string
	MD5      string
}
