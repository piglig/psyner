package event

import (
	"context"
	"psyner/common"
	"psyner/server/taskrun/action"
)

var (
	_ action.Executor = (*GetFileExecutor)(nil)
	_ action.Executor = (*UpdateFileExecutor)(nil)
	_ action.Executor = (*DeleteFileExecutor)(nil)
)

type GetFileExecutor struct {
}

type UpdateFileExecutor struct {
}

type DeleteFileExecutor struct {
}

func (*GetFileExecutor) Exec(ctx context.Context, action common.FileSyncActionType, command string) error {

	return nil
}

func (*UpdateFileExecutor) Exec(ctx context.Context, action common.FileSyncActionType, command string) error {

	return nil
}

func (*DeleteFileExecutor) Exec(ctx context.Context, action common.FileSyncActionType, command string) error {

	return nil
}
