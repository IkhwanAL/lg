package core

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
)

// #define MAX_FILE 10
var MAX_FILE = 10

// TODO Need A Test File
func SearchFile(key string) ([]string, error) {

	if key == "" {
		return nil, errors.New("key is empty")
	}

	var maxPathShow []string

	filepath.WalkDir("D:/", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, _ := d.Info()

		if info.IsDir() && strings.Contains(info.Name(), "$") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		if !strings.ContainsAny(info.Name(), key) {
			return nil
		}

		if len(maxPathShow) >= MAX_FILE {
			return filepath.SkipAll
		}

		maxPathShow = append(maxPathShow, path)

		return nil
	})

	return maxPathShow, nil
}
