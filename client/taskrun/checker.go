package taskrun

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"psyner/common"
	"time"
)

func CheckLocalDirChecksum(localDir string, interval time.Duration) {
	// TODO periodically check if local directory is consistent with server
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			checkSum := make(map[string]string)
			err := filepath.Walk(localDir, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.Mode().IsRegular() {
					return nil
				}

				checksum, err := common.GenerateChecksum(path)
				if err != nil {
					return err
				}

				relPath, err := filepath.Rel(localDir, path)
				if err != nil {
					return err
				}
				checkSum[relPath] = checksum
				fmt.Printf("time:%v %s: %s\n", time.Now(), path, checksum)
				return nil
			})

			if err != nil {
				return
			}

			// TODO compare with server checksum, get not exist file from server
		}
	}
}
