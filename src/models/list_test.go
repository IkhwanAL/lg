package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TEST_FILE = []string{
	"FILE 1",
	"FILE 2",
	"FILE 3",
	"FILE 4",
	"FILE 5",
	"FILE 6",
	"FILE 7",
	"FILE 8",
	"FILE 9",
	"FILE 10",
	"FILE 11",
	"FILE 12",
	"FILE 13",
	"FILE 14",
	"FILE 15",
	"FILE 16",
	"FILE 17",
	"FILE 18",
	"FILE 19",
}

func TestListPositionHeadTail(t *testing.T) {
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
	}, "After Press Key Down, The Two Pointer Not Meet Condition Of View Height Argument")

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
	}, "After Press Key Up, The Two Pointer Not Meet Condition Of View Height Argument")

}
