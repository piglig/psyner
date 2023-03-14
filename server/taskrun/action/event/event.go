package event

import (
	"context"
	"encoding/json"
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

func (*GetFileExecutor) Exec(ctx context.Context, command string) error {
	p := common.GetFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}
	return nil
}

func (*UpdateFileExecutor) Exec(ctx context.Context, command string) error {
	p := common.UpdateFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}
	return nil
}

func (*DeleteFileExecutor) Exec(ctx context.Context, command string) error {
	p := common.DeleteFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}
	return nil
}
