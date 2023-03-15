package event

import (
	"psyner/common"
	"psyner/server/taskrun/action"
)

func init() {
	action.RegisterExecutor(common.GetFileSync, &GetFileExecutor{})
	action.RegisterExecutor(common.UpdateFileSync, &UpdateFileExecutor{})
	action.RegisterExecutor(common.DeleteFileSync, &DeleteFileExecutor{})

}
