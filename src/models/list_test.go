package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ikhwanal/lg-file/src/core"
	"github.com/stretchr/testify/assert"
)

func generateTestFile(totalFile int) []core.FsEntry {
	list := make([]core.FsEntry, totalFile)
	for i := range totalFile {

		list[i] = core.FsEntry{
			Name: fmt.Sprintf("File %d", i),
			Type: core.File,
			Path: "./",
		}
	}

	return list
}

/*
A Test For Verifying Head Tail Movement
*/
func TestListPositionHeadTail(t *testing.T) {

	TEST_FILE := generateTestFile(24)

	tail := len(TEST_FILE)
	position := 0
	viewHeight := 10

	list := ListModel{
		position: position,
		list:     TEST_FILE,
		Height:   viewHeight,
		tail:     min(viewHeight, tail),
		cursor:   0,
		head:     0,
	}

	list.position++

	position += 1

	assert.Equal(t, position, list.position, "Press Key Down And Position Must Moving In The List")

	for range viewHeight {
		list.position++
		position += 1
	}

	list.Move()

	assert.Equal(t, 1, list.head, "After Key Down The Head Pointer Must Move To Show Different List")
	assert.Equal(t, 11, list.tail, "After Key Down The Tail Pointer Must Move To Show Different List")

	assert.Condition(t, func() (success bool) {
		return (list.tail-list.head)-1 < viewHeight
	}, "After Press Key Down, The Length of Two Pointer is too big for View Height Argument")

	list.position -= 1
	position -= 1

	assert.Equal(t, position, list.position, "Press Key Up And Position Must Moving In The List")

	for range viewHeight {
		list.position--
		position -= 1
	}

	list.Move()

	assert.Equal(t, 0, list.head, "After Key Up The Head Pointer Must Move To Show Different List")
	assert.Equal(t, 10, list.tail, "After Key Up The Tail Pointer Must Move To Show Different List")

	assert.Condition(t, func() (success bool) {
		return (list.tail-list.head)-1 < viewHeight
	}, "After Press Key Up, The Length of Two Pointer is too big for View Height Argument")

}

func isValidPath(path string) bool {
	if strings.HasSuffix(path, "/") {
		return false
	}

	if strings.Contains(path, "//") {
		return false
	}

	if strings.Contains(path, "/\\") {
		return false
	}

	return true
}

func TestDirectoryList(t *testing.T) {
	list := defaultList("./")

	for _, v := range list {
		if v.Type != core.Dir {
			continue
		}

		assert.False(t, isValidPath(v.Path))
		assert.DirExists(t, v.Path, "Directory not exist, something wrong with the path when open directory", v.Path)
	}
}

func TestFileList(t *testing.T) {
	list := defaultList("./")

	for _, v := range list {
		if v.Type != core.File {
			continue
		}

		assert.FileExists(t, v.Path, "file cannot open, something wrong withg the path when open the file", v.Path)
	}
}

func TestRunningListModel(t *testing.T) {
	testFile := generateTestFile(9);

	model := NewListModel(120, 35, "./", &core.UserArgs{})

	model.OverrideList(testFile)

	assert.Condition(t, func() (success bool) {
		return model.tail == len(model.list)
	}, "Tail cannot bigger than total Height")
}
