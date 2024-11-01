package tests

import (
	"testing"

	"github.com/Pieeer1/Auto-SwagGo/internal/ext"
)

func TestSequenceEqual(t *testing.T) {
	intArray := []int{1, 2, 3, 4, 5}
	secondIntArray := []int{1, 2, 3, 4, 5}

	if !ext.SequenceEqual(intArray, secondIntArray) {
		t.Errorf("Expected true, got false")
	}

	secondIntArray = []int{1, 2, 3, 4, 6}

	if ext.SequenceEqual(intArray, secondIntArray) {
		t.Errorf("Expected false, got true")
	}

	secondIntArray = []int{1, 2, 3, 4}

	if ext.SequenceEqual(intArray, secondIntArray) {
		t.Errorf("Expected false, got true")
	}

	secondIntArray = []int{1, 2, 3, 4, 5, 6}

	if ext.SequenceEqual(intArray, secondIntArray) {
		t.Errorf("Expected false, got true")
	}

	secondIntArray = []int{1, 2, 3, 5, 4}

	if ext.SequenceEqual(intArray, secondIntArray) {
		t.Errorf("Expected false, got true")
	}
}
