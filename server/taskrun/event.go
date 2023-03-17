package taskrun

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"net"
	"os"
	"psyner/common"
)

func init() {
	RegisterHandler(fsnotify.Create, &CreateFileHandler{})
	RegisterHandler(fsnotify.Write, &ModifyFileHandler{})
}

var (
	_ Executor = (*GetFileExecutor)(nil)
	_ Executor = (*UpdateFileExecutor)(nil)
	_ Executor = (*DeleteFileExecutor)(nil)
)

type GetFileExecutor struct {
}

func (e *GetFileExecutor) Check(ctx Context, command string) error {
	p := common.GetFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}

	if p.RelPath == "" {
		return fmt.Errorf("GetFileExecutor invalid params:%v\n", p)
	}

	filePath := ctx.getFilePath(p.RelPath)
	ex := ctx.CheckFileExist(p.RelPath)
	if !ex {
		return fmt.Errorf("GetFileExecutor file_path %s not exist", filePath)
	}

	//TODO implement me
	log.Println("GetFileExecutor", "file_path", filePath, "exist")
	return nil
}

type UpdateFileExecutor struct {
}

func (e *UpdateFileExecutor) Check(ctx Context, command string) error {
	//TODO implement me
	log.Println("implement me")
	return nil
}

type DeleteFileExecutor struct {
}

func (e *DeleteFileExecutor) Check(ctx Context, command string) error {
	//TODO implement me
	log.Println("implement me")
	return nil
}

func (*GetFileExecutor) Exec(ctx Context, conn net.Conn, command string) error {
	p := common.GetFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}

	encoder := gob.NewEncoder(conn)

	filePath := ctx.getFilePath(p.RelPath)
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	defer f.Close()
	res := common.GetFileSyncPayloadRes{RelPath: p.RelPath, FileSize: fInfo.Size()}
	err = encoder.Encode(&res)
	if err != nil {
		return err
	}

	_, err = io.Copy(conn, f)
	if err != nil {
		return err
	}

	log.Println("GetFileExecutor from", conn.RemoteAddr())
	return err
}

func (*UpdateFileExecutor) Exec(ctx Context, conn net.Conn, command string) error {
	p := common.UpdateFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}
	return nil
}

func (*DeleteFileExecutor) Exec(ctx Context, conn net.Conn, command string) error {
	p := common.DeleteFileSyncPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}
	return nil
}
