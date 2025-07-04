package core

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

// #define MAX_FILE 20
var MAX_FILE = 20

// TODO Need A Test File
func SearchFile(key string) []string {

	var maxPathShow []string

	pathRecord := 0
	filepath.WalkDir("D:/", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		info, _ := d.Info()

		if info.IsDir() && strings.Contains(info.Name(), "$") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		if pathRecord > MAX_FILE {
			return filepath.SkipAll
		}

		pathRecord++
		maxPathShow = append(maxPathShow, path)

		return nil
	})
	log.Print(maxPathShow)
	return maxPathShow
}
