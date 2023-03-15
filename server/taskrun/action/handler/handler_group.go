package handler

import (
	"github.com/fsnotify/fsnotify"
	"psyner/server/taskrun/action"
)

func init() {
	action.RegisterHandler(fsnotify.Create, &CreateFileHandler{})
	action.RegisterHandler(fsnotify.Write, &ModifyFileHandler{})
}
