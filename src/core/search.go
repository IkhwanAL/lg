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

var SKIP_KNOWN_DIRECTIRIES = []string{
	"node_modules",
	".git",
	"Recovery",
	"System Volume Information",
	"$RECYCLE.BIN",
}

// #define MAX_FILE 10
var MAX_FILE = 1000

// 2.5-ish second
func SearchFile(key string) ([]string, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	var maxPathShow []string

	// TODO
	// Improve Search
	// How to improve search if we search thousand directories
	// does goroutine help?

	// TODO
	// do i need to consider adding C:/
	// because C:/ has soo much private information that windows user shouldn't know
	// It means lots of filtering that other partition shouldn't do

	// TODO
	// i should highlight directory name

	// TODO
	// I need to check their what partition they have
	// or partition filter only allow to find in one partition
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

// 0.5s or < 1s
func SearchFileV3(key string) ([]string, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	var wg sync.WaitGroup

	ch := make(chan string, 5)

	wg.Add(1)
	go search("D:/", key, ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	var result []string

	for path := range ch {
		result = append(result, path)
	}

	return result, nil
}

func search(rootDir string, fileToSearch string, result chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	dir, err := os.ReadDir(rootDir)

	if err != nil {
		log.Print(err)
		return
	}

	for _, d := range dir {
		if d.IsDir() {
			wg.Add(1)
			go search(rootDir+d.Name()+"/", fileToSearch, result, wg)

			continue
		}

		if !strings.Contains(d.Name(), fileToSearch) {
			continue
		}

		result <- rootDir + d.Name()
	}
}
