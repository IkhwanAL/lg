package models

import (
	"fmt"
	"testing"

	"github.com/ikhwanal/everywhere_anywhere/src/core"
	"github.com/stretchr/testify/assert"
)

func generateTestFile() []core.FsEntry {
	list := make([]core.FsEntry, 15)

	for i := range 15 {
		list = append(list, core.FsEntry{
			Name: fmt.Sprintf("File %d", i),
			Type: core.File,
			Path: "./",
		})
	}

	return list
}

func TestListPositionHeadTail(t *testing.T) {

	TEST_FILE := generateTestFile()

	tail := len(TEST_FILE) - 1
	position := 0
	viewHeight := 10

	list := ListModel{
		position:   position,
		list:       TEST_FILE,
		viewHeight: viewHeight,
		tail:       min(viewHeight, tail),
		cursor:     0,
		head:       0,
	}

	list.position++

	position += 1

	assert.Equal(t, position, list.position, "Press Key Down And Position Must Moving In The List")

	for range 10 {
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

	for range 10 {
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
