package action

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"net"
	"psyner/common"
	"psyner/server/ctx"
	"sync"
)

var (
	executorMu sync.RWMutex
	executors  = make(map[common.FileSyncActionType]Executor)

	handlerMu sync.RWMutex
	handlers  = make(map[fsnotify.Op]Handler)
)

type Executor interface {
	Exec(ctx ctx.Context, conn net.Conn, command string) error
	Check(ctx ctx.Context, command string) error
}

type Handler interface {
	Do(event fsnotify.Event) ([]byte, error)
}

func RegisterExecutor(action common.FileSyncActionType, f Executor) {
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

func RegisterHandler(op fsnotify.Op, f Handler) {
	handlerMu.Lock()
	defer handlerMu.Unlock()

	if f == nil {
		panic("executor: Register executor is nil")
	}

	_, ok := handlers[op]
	if !ok {
		handlers[op] = f
	} else {
		panic(fmt.Sprintf("executor: Register called twice for %v", op))
	}
}
