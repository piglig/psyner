package common

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
