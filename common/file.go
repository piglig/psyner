package common

type FileSyncActionType string

const (
	GetFileSync    FileSyncActionType = "GET"
	UpdateFileSync                    = "UPDATE"
	DeleteFileSync                    = "DELETE"
)

type GetFileSyncPayload struct {
	relPath string
}

type UpdateFileSyncPayload struct {
	relPath  string
	fileHash string
}

type DeleteFileSyncPayload struct {
	relPath string
}
