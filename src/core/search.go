package core

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// #define MAX_FILE 10
var MAX_FILE = 1000

// 2.5-ish second
func SearchFile(key string) ([]string, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	var maxPathShow []string
	err := filepath.WalkDir("D:/", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if d.Name() == "System Volume Information" || d.Name() == "$RECYCLE.BIN" || d.Name() == "Recovery" {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.Contains(d.Name(), key) {
			return nil
		}

		if len(maxPathShow) >= MAX_FILE {
			return filepath.SkipAll
		}

		maxPathShow = append(maxPathShow, path)

		return nil
	})

	return maxPathShow, err
}

// 1.8 second
func SearchFileV2(key string) ([]string, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	dirs, err := os.ReadDir("D:/")

	if err != nil {
		return nil, err
	}

	var results []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, dir := range dirs {

		_, err := dir.Info()

		if err != nil {
			log.Print(err)
			continue
		}

		if strings.Contains(dir.Name(), ".") {
			continue
		}

		if dir.Name() == "System Volume Information" || dir.Name() == "Recovery" || dir.Name() == "$RECYCLE.BIN" {
			continue
		}

		wg.Add(1)
		go func(dirName string) {
			defer wg.Done()
			err = filepath.WalkDir("D:/"+dirName, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() {
					return nil
				}

				if !strings.Contains(d.Name(), key) {
					return nil
				}

				mu.Lock()
				results = append(results, path)
				mu.Unlock()

				return nil
			})
		}(dir.Name())

		if err != nil {
			log.Print(err)
			continue
		}
	}

	wg.Wait()

	return results, err
}

// TODO i got hit in C:/ Which lots of access denied
// If in one folder has lots of restriction or access denied folder
// Turn off the one layer folder
func SearchFileV3(path string, key string) ([]FsEntry, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	var wg sync.WaitGroup

	ch := make(chan FsEntry, 10)

	wg.Add(1)
	go search(path, key, ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	var result []FsEntry

	for path := range ch {
		result = append(result, path)
	}

	return result, nil
}

func search(rootDir string, fileToSearch string, result chan<- FsEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	dir, err := os.ReadDir(rootDir)

	if err != nil {
		log.Print(err)
		return
	}

	for _, d := range dir {
		if d.IsDir() {
			wg.Add(1)
			go search(filepath.FromSlash(rootDir+d.Name()+"/"), fileToSearch, result, wg)

			continue
		}

		if !strings.Contains(d.Name(), fileToSearch) {
			continue
		}

		result <- FsEntry{
			Name: rootDir + d.Name(),
			Type: File,
			Path: filepath.ToSlash(rootDir + d.Name()),
		}
	}
}
