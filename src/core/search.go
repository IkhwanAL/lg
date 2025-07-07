package core

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
)

// #define MAX_FILE 10
var MAX_FILE = 10

func SearchFile(key string) ([]string, error) {

	if key == "" {
		return nil, errors.New("key is empty")
	}

	var maxPathShow []string

	// TODO
	// Improve Search

	// TODO
	// do i need to consider adding C:/
	// because C:/ has soo much private information that windows user shouldn't know
	// It means lots of filtering that other partition shouldn't do

	// TODO
	// i should highlight directory name

	// TODO
	// I need to check their what partition they have
	// or partition filter only allow to find in one partition
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

		if !strings.Contains(info.Name(), key) {
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
