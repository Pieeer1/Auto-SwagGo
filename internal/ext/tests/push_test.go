package tests

import (
	"testing"

	"github.com/Pieeer1/Auto-SwagGo/internal/ext"
)

func TestPush(t *testing.T) {
	ints := []int{2, 3, 4}

	ints = ext.Push(ints, 1)

	if ints[0] != 1 {
		t.Errorf("Expected 1, got %d", ints[0])
	}
	if ints[1] != 2 {
		t.Errorf("Expected 2, got %d", ints[1])
	}
	if ints[2] != 3 {
		t.Errorf("Expected 3, got %d", ints[2])
	}
	if ints[3] != 4 {
		t.Errorf("Expected 4, got %d", ints[3])
	}

	if len(ints) != 4 {
		t.Errorf("Expected length 4, got %d", len(ints))
	}
}
