package common

type FileSyncActionType int

const (
	GetFileSync FileSyncActionType = iota + 1
	DeleteFileSync
)

type FileSyncProto struct {
	relPath  string
	fileHash string
	fileSize int64
}

func (p FileSyncProto) GetRelPath() string {
	return p.relPath
}

func (p FileSyncProto) GetFileHash() string {
	return p.fileHash
}

func (p FileSyncProto) GetFileSize() int64 {
	return p.fileSize
}
