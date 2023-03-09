package event

import (
	"psyner/common"
	"psyner/server/taskrun/action"
)

func Init() {
	action.Register(common.GetFileSync, &GetFileExecutor{})
	action.Register(common.UpdateFileSync, &UpdateFileExecutor{})
	action.Register(common.DeleteFileSync, &DeleteFileExecutor{})
}
