package common

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func GenerateChecksum(filePath string) (string, error) {
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
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
