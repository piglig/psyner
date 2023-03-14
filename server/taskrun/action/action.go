package action

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"psyner/common"
	"sync"
)

type Executor interface {
	Exec(ctx context.Context, conn net.Conn, command string) error
	Check(ctx context.Context, command string) error
}

var (
	executorMu sync.RWMutex
	executors  = make(map[common.FileSyncActionType]Executor)
)

func Register(action common.FileSyncActionType, f Executor) {
	executorMu.Lock()
	defer executorMu.Unlock()

	if f == nil {
		panic("executor: Register executor is nil")
	}

	_, ok := executors[action]
	if !ok {
		executors[action] = f
	} else {
		panic(fmt.Sprintf("executor: Register called twice for %v", action))
	}
}

func Exec(ctx context.Context, action common.FileSyncActionType, conn net.Conn, command string) (err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errors.Errorf("panic in executor exec: %v", panicErr)
		}
	}()

	executorMu.RLock()
	defer executorMu.RUnlock()
	f, ok := executors[action]
	if !ok {
		return errors.Errorf("executor: unknow action %v", action)
	}

	err = f.Check(ctx, command)
	if err != nil {
		return err
	}

	return f.Exec(ctx, conn, command)
}
