package util

import "testing"

func TestShiftLeftArray(t *testing.T) {
	currentView := []string{"a", "b", "c"}
	correctResult := make([]string, 3)

	correctResult[0] = "b"
	correctResult[1] = "c"

	result := ShiftLeftArray(currentView)

	if result[0] != correctResult[0] {
		t.Errorf("Shifting Index 0 Failed Result: %s, \t Correct Result: %s", result[0], correctResult[0])
	}

	if result[1] != correctResult[1] {
		t.Errorf("Shifting Index 1 Failed Result: %s, \t Correct Result: %s", result[1], correctResult[1])
	}

	if result[2] != correctResult[2] {
		t.Errorf("Shifting Index 2 Failed Result: %v, \t Correct Result: %v", result[2], correctResult[2])
	}
}
