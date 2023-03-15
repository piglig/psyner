package handler

import (
	"github.com/fsnotify/fsnotify"
	"psyner/server/taskrun/action"
)

var (
	_ action.Handler = (*DeleteFileHandler)(nil)
	_ action.Handler = (*ModifyFileHandler)(nil)
)

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
