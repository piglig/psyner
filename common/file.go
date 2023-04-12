package common

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

type FileSyncActionType uint32

const (
	GetFileSync FileSyncActionType = iota + 1
	UpdateFileSync
	DeleteFileSync
)

const (
	BufferSize = 1024
)

type FileSyncPayload struct {
	ActionType    FileSyncActionType
	ActionPayload []byte
}

func (f FileSyncPayload) Byte() []byte {
	bs, _ := json.Marshal(f)
	return bs
}

func (f FileSyncPayload) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, f.ActionType)
	if err != nil {
		return 0, err
	}

	n = 4
	err = binary.Write(w, binary.BigEndian, uint32(len(f.ActionPayload)))
	if err != nil {
		return 0, err
	}
	n += 4
	o, err := w.Write(f.ActionPayload)
	if err != nil {
		return 0, err
	}

	return n + int64(o), err
}

func (f *FileSyncPayload) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, f.ActionType)
	if err != nil {
		return 0, err
	}

	n = 4
	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return 0, err
	}
	n += 4

	buf := make([]byte, size)
	o, err := r.Read(buf)
	if err != nil {
		return 0, err
	}

	f.ActionPayload = buf
	return n + int64(o), nil
}

func (f FileSyncPayload) String() string {
	return string(f.Byte())
}

type GetFileSyncPayload struct {
	RelPath string
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

type FileCommandPayload interface {
	Byte() []byte
	io.WriterTo
	io.ReaderFrom
	fmt.Stringer
}
