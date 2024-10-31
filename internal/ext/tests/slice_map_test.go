package tests

import (
	"auto-swaggo/internal/ext"
	"testing"
)

type TestSliceMapStruct struct {
	Name  string
	Value int
}

func TestSliceMap(t *testing.T) {

	testSliceMapStructs := []TestSliceMapStruct{
		{Name: "a", Value: 1},
		{Name: "b", Value: 2},
		{Name: "c", Value: 3},
	}

	testSliceMapStructsMapped := ext.SliceMap(testSliceMapStructs, func(d TestSliceMapStruct) string {
		return d.Name
	})

	if len(testSliceMapStructsMapped) != 3 {
		t.Error("Expected 3, got", len(testSliceMapStructsMapped))
	}

	if testSliceMapStructsMapped[0] != "a" {
		t.Error("Expected a, got", testSliceMapStructsMapped[0])
	}

	if testSliceMapStructsMapped[1] != "b" {
		t.Error("Expected b, got", testSliceMapStructsMapped[1])
	}

	if testSliceMapStructsMapped[2] != "c" {
		t.Error("Expected c, got", testSliceMapStructsMapped[2])
	}

	testSliceMapStructsMappedInt := ext.SliceMap(testSliceMapStructs, func(d TestSliceMapStruct) int {
		return d.Value
	})

	if len(testSliceMapStructsMappedInt) != 3 {
		t.Error("Expected 3, got", len(testSliceMapStructsMappedInt))
	}

	if testSliceMapStructsMappedInt[0] != 1 {
		t.Error("Expected 1, got", testSliceMapStructsMappedInt[0])
	}

	if testSliceMapStructsMappedInt[1] != 2 {
		t.Error("Expected 2, got", testSliceMapStructsMappedInt[1])
	}
}
