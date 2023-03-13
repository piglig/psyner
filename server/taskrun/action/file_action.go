package action

import (
	"context"
	"psyner/common"
)

func FileSyncAction(ctx context.Context, action common.FileSyncActionType, command string) error {
	return Exec(ctx, action, command)
}
