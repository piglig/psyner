package action

import (
	"context"
	"net"
	"psyner/common"
)

func FileSyncAction(ctx context.Context, action common.FileSyncActionType, conn net.Conn, command string) error {
	return Exec(ctx, action, conn, command)
}
