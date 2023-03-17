package taskrun

import (
	"github.com/fsnotify/fsnotify"
	"psyner/common"
)

var (
	_ Handler = (*DeleteFileHandler)(nil)
	_ Handler = (*ModifyFileHandler)(nil)
)

func init() {
	RegisterExecutor(common.GetFileSync, &GetFileExecutor{})
	RegisterExecutor(common.UpdateFileSync, &UpdateFileExecutor{})
	RegisterExecutor(common.DeleteFileSync, &DeleteFileExecutor{})

}

type ModifyFileHandler struct {
}

type DeleteFileHandler struct {
}

type CreateFileHandler struct {
}

func (d *DeleteFileHandler) Do(event fsnotify.Event) ([]byte, error) {
	//TODO implement me
	return nil, nil
}

func (c *CreateFileHandler) Do(event fsnotify.Event) ([]byte, error) {
	//TODO implement me

	return nil, nil
}

func (m *ModifyFileHandler) Do(event fsnotify.Event) ([]byte, error) {
	//TODO implement me
	return nil, nil
}
