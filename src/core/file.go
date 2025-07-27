package core

import (
	"os"
	"path/filepath"
)

type FileType int

const (
	File FileType = iota
	Dir
)

type FsEntry struct {
	Name string
	Type FileType
	Path string
}

func CreateNewFile(path string, filename string) error {
	completePath := filepath.Join(path + "/" + filename)

	file, err := os.Create(filepath.ToSlash(completePath))

	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}
