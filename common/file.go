package common

import (
	"encoding/json"
	"fmt"
	"io"
)

type FileSyncActionType string

const (
	GetFileSync    FileSyncActionType = "GET"
	UpdateFileSync                    = "UPDATE"
	DeleteFileSync                    = "DELETE"
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

func (f *FileSyncPayload) WriteTo(w io.Writer) (n int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *FileSyncPayload) ReadFrom(r io.Reader) (n int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *FileSyncPayload) String() string {
	//TODO implement me
	panic("implement me")
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
