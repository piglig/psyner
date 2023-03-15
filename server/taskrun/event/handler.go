package event

import (
	"github.com/fsnotify/fsnotify"
	"psyner/common"
	"psyner/server/taskrun/runner"
)

var (
	_ runner.Handler = (*DeleteFileHandler)(nil)
	_ runner.Handler = (*ModifyFileHandler)(nil)
)

func init() {
	runner.RegisterExecutor(common.GetFileSync, &GetFileExecutor{})
	runner.RegisterExecutor(common.UpdateFileSync, &UpdateFileExecutor{})
	runner.RegisterExecutor(common.DeleteFileSync, &DeleteFileExecutor{})

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
