package event

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"psyner/common"
	"psyner/server/ctx"
	"psyner/server/taskrun/action"
)

var (
	_ action.Executor = (*GetFileExecutor)(nil)
	_ action.Executor = (*UpdateFileExecutor)(nil)
	_ action.Executor = (*DeleteFileExecutor)(nil)
)

type GetFileExecutor struct {
}

func (e *GetFileExecutor) Check(ctx ctx.Context, command string) error {
	p := common.GetFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}

	if p.RelPath == "" {
		return fmt.Errorf("GetFileExecutor invalid params:%v\n", p)
	}

	ex := ctx.CheckFileExist(p.RelPath)

	//TODO implement me
	log.Println("GetFileExecutor", "file_path", p.RelPath, "exist", ex)
	return nil
}

type UpdateFileExecutor struct {
}

func (e *UpdateFileExecutor) Check(ctx ctx.Context, command string) error {
	//TODO implement me
	log.Println("implement me")
	return nil
}

type DeleteFileExecutor struct {
}

func (e *DeleteFileExecutor) Check(ctx ctx.Context, command string) error {
	//TODO implement me
	log.Println("implement me")
	return nil
}

func (*GetFileExecutor) Exec(ctx ctx.Context, conn net.Conn, command string) error {
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

func (*UpdateFileExecutor) Exec(ctx ctx.Context, conn net.Conn, command string) error {
	p := common.UpdateFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}
	return nil
}

func (*DeleteFileExecutor) Exec(ctx ctx.Context, conn net.Conn, command string) error {
	p := common.DeleteFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}
	return nil
}
