package common

type FileSyncActionType string

const (
	GetFileSync    FileSyncActionType = "GET"
	UpdateFileSync                    = "UPDATE"
	DeleteFileSync                    = "DELETE"
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
