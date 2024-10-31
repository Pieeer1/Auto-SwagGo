package tests

import (
	"auto-swaggo/internal/ext"
	"testing"
)

type FlattenMapTestStruct struct {
	ArrayValues []int
}

func TestFlattenMap(t *testing.T) {

	flattenMapTestStructs := []FlattenMapTestStruct{
		{ArrayValues: []int{1, 2, 3}},
		{ArrayValues: []int{4, 5, 6}},
	}

	flattenMapTestStructsFlattened := ext.FlattenMap(flattenMapTestStructs, func(f FlattenMapTestStruct) []int {
		return f.ArrayValues
	})

	if len(flattenMapTestStructsFlattened) != 6 {
		t.Error("Expected 6, got", len(flattenMapTestStructsFlattened))
	}

	for i := 0; i < 6; i++ {
		if flattenMapTestStructsFlattened[i] != i+1 {
			t.Error("Expected true, got false")
		}
	}

}
