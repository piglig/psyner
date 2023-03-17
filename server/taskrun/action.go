package taskrun

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"log"
	"net"
	"psyner/common"
	"sync"
)

type Executor interface {
	Exec(ctx Context, conn net.Conn, command string) error
	Check(ctx Context, command string) error
}

type Handler interface {
	Do(event fsnotify.Event) ([]byte, error)
}

var (
	executorMu sync.RWMutex
	executors  = make(map[common.FileSyncActionType]Executor)

	handlerMu sync.RWMutex
	handlers  = make(map[fsnotify.Op]Handler)
)

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

func GetExecutor(action common.FileSyncActionType) (Executor, error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Printf("panic in executor exec: %v\n", panicErr)
		}
	}()

	executorMu.RLock()
	defer executorMu.RUnlock()
	f, ok := executors[action]
	if !ok {
		return nil, errors.Errorf("executor: unknow action %v", action)
	}
	return f, nil
}

func GetHandler(op fsnotify.Op) (Handler, error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Printf("panic in handler err: %v\n", panicErr)
			return
		}
	}()

	handlerMu.RLock()
	defer handlerMu.RUnlock()
	f, ok := handlers[op]
	if !ok {
		return nil, errors.Errorf("handler: unknow operate %v", op)
	}

	return f, nil
}

func Do(ctx Context, op fsnotify.Op, conn net.Conn, event fsnotify.Event) (err error) {
	f, err := GetHandler(op)
	res, err := f.Do(event)
	if err != nil {
		return err
	}

	_ = res
	return
}

func Exec(ctx Context, act common.FileSyncActionType, conn net.Conn, command string) (err error) {
	f, err := GetExecutor(act)
	if err != nil {
		return err
	}

	err = f.Check(ctx, command)
	if err != nil {
		return err
	}

	return f.Exec(ctx, conn, command)
}
