package taskrun

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func CheckLocalDirChecksum(localDir string, ticker *time.Ticker) {
	// TODO periodically check if local directory is consistent with server
	for range ticker.C {
		checkSum := make(map[string]string)
		err := filepath.Walk(localDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			checksum, err := generateChecksum(path)
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

func generateChecksum(filePath string) (string, error) {
	// Generate the checksum for a file given its path
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
