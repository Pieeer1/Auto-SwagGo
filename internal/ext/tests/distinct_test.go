package tests

import (
	"auto-swaggo/internal/ext"
	"testing"
)

func TestDistinct(t *testing.T) {
	ints := []int{1, 2, 3, 2, 3, 4}
	distinctInts := ext.Distinct(ints)

	if len(distinctInts) != 4 {
		t.Errorf("Expected 4, got %d", len(distinctInts))
	}

	if !ext.Contains(distinctInts, 1) {
		t.Error("Expected true, got false")
	}

	if !ext.Contains(distinctInts, 2) {
		t.Error("Expected true, got false")
	}

	if !ext.Contains(distinctInts, 3) {
		t.Error("Expected true, got false")
	}

	if !ext.Contains(distinctInts, 4) {
		t.Error("Expected true, got false")
	}

	if ext.Contains(distinctInts, 5) {
		t.Error("Expected false, got true")
	}
}
