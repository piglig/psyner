package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"psyner/common"
	"psyner/server/cmd"
	"psyner/server/taskrun/action"
)

var (
	_ action.Executor = (*GetFileExecutor)(nil)
	_ action.Executor = (*UpdateFileExecutor)(nil)
	_ action.Executor = (*DeleteFileExecutor)(nil)
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

	s, ok := ctx.Value("server").(*cmd.Server)
	if !ok {
		return fmt.Errorf("GetFileExecutor invalid context\n")
	}

	ex := s.CheckFileExist(p.RelPath)

	//TODO implement me
	log.Println("GetFileExecutor", "file_path", p.RelPath, "exist", ex)
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
	log.Println("GetFileExecutor from", conn.RemoteAddr())
	return nil
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
