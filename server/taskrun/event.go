package taskrun

import (
	"encoding/gob"
	"encoding/json"
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

type UpdateFileExecutor struct {
}

type DeleteFileExecutor struct {
}

func (*GetFileExecutor) Exec(ctx Context, conn net.Conn, command string) error {
	p := common.GetFileOpPayload{}
	if err := json.Unmarshal([]byte(command), &p); err != nil {
		return err
	}

	encoder := gob.NewEncoder(conn)

	filePath := ctx.getFilePath(p.RelPath)
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	//fInfo, err := f.Stat()
	//if err != nil {
	//	return err
	//}
	res := common.GetFileSyncPayloadRes{RelPath: p.RelPath, FileSize: 0}
	err = encoder.Encode(&res)
	if err != nil {
		return err
	}

	_, err = io.Copy(conn, f)
	if err != nil {
		return err
	}

	//conn.Close()

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
