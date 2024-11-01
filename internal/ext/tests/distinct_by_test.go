package tests

import (
	"testing"

	"github.com/Pieeer1/Auto-SwagGo/internal/ext"
)

type DistinctByTestStruct struct {
	Name  string
	Value int
}

func TestDistinctBy(t *testing.T) {

	distinctByTestStructs := []DistinctByTestStruct{
		{Name: "a", Value: 1},
		{Name: "b", Value: 2},
		{Name: "a", Value: 3},
		{Name: "b", Value: 4},
		{Name: "c", Value: 4},
	}

	distinctByTestStructsDistinct := ext.DistinctBy(distinctByTestStructs, func(d DistinctByTestStruct) string {
		return d.Name
	})

	if len(distinctByTestStructsDistinct) != 3 {
		t.Error("Expected 3, got", len(distinctByTestStructsDistinct))
	}

	if distinctByTestStructsDistinct[0].Name != "a" {
		t.Error("Expected a, got", distinctByTestStructsDistinct[0].Name)
	}

	if distinctByTestStructsDistinct[1].Name != "b" {
		t.Error("Expected b, got", distinctByTestStructsDistinct[1].Name)
	}

	if distinctByTestStructsDistinct[2].Name != "c" {
		t.Error("Expected c, got", distinctByTestStructsDistinct[2].Name)
	}

	distinctByTestStructsDistinct = ext.DistinctBy(distinctByTestStructs, func(d DistinctByTestStruct) int {
		return d.Value
	})

	if len(distinctByTestStructsDistinct) != 4 {
		t.Error("Expected 4, got", len(distinctByTestStructsDistinct))
	}

	if distinctByTestStructsDistinct[0].Name != "a" {
		t.Error("Expected a, got", distinctByTestStructsDistinct[0].Name)
	}

	if distinctByTestStructsDistinct[1].Name != "b" {
		t.Error("Expected b, got", distinctByTestStructsDistinct[1].Name)
	}

	if distinctByTestStructsDistinct[2].Name != "a" {
		t.Error("Expected a, got", distinctByTestStructsDistinct[2].Name)
	}

	if distinctByTestStructsDistinct[3].Name != "b" {
		t.Error("Expected c, got", distinctByTestStructsDistinct[3].Name)
	}
}
