package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

// GetDirSize traverses directory and returns the total size of files in bytes
func GetDirSize(root string) (int64, error) {
	var size int64
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
