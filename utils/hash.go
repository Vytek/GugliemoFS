package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// FileHash computes the SHA-256 hash of a file at the given path.
// It returns the hash as a hexadecimal string or an error if the file cannot be read.
func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
