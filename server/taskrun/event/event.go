package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"net"
	"os"
	"psyner/common"
	"psyner/server/taskrun/runner"
)

func init() {
	runner.RegisterHandler(fsnotify.Create, &CreateFileHandler{})
	runner.RegisterHandler(fsnotify.Write, &ModifyFileHandler{})
}

var (
	_ runner.Executor = (*GetFileExecutor)(nil)
	_ runner.Executor = (*UpdateFileExecutor)(nil)
	_ runner.Executor = (*DeleteFileExecutor)(nil)
)

type GetFileExecutor struct {
}

func (e *GetFileExecutor) Check(ctx context.Context, command string) error {
	p := common.GetFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}

	if p.RelPath == "" {
		return fmt.Errorf("GetFileExecutor invalid params:%v\n", p)
	}

	//ex := ctx.CheckFileExist(p.RelPath)

	//TODO implement me
	log.Println("GetFileExecutor", "file_path", p.RelPath, "exist")
	return nil
}

type UpdateFileExecutor struct {
}

func (e *UpdateFileExecutor) Check(ctx context.Context, command string) error {
	//TODO implement me
	log.Println("implement me")
	return nil
}

type DeleteFileExecutor struct {
}

func (e *DeleteFileExecutor) Check(ctx context.Context, command string) error {
	//TODO implement me
	log.Println("implement me")
	return nil
}

func (*GetFileExecutor) Exec(ctx context.Context, conn net.Conn, command string) error {
	p := common.GetFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}

	f, err := os.Open(p.RelPath)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = io.Copy(conn, f)

	log.Println("GetFileExecutor from", conn.RemoteAddr())
	return err
}

func (*UpdateFileExecutor) Exec(ctx context.Context, conn net.Conn, command string) error {
	p := common.UpdateFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}
	return nil
}

func (*DeleteFileExecutor) Exec(ctx context.Context, conn net.Conn, command string) error {
	p := common.DeleteFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}
	return nil
}
