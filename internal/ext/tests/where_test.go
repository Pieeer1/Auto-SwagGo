package tests

import (
	"auto-swaggo/internal/ext"
	"testing"
)

type WhereTestStruct struct {
	Name  string
	Value int
}

func TestWhere(t *testing.T) {
	testWhereStructs := []WhereTestStruct{
		{Name: "a", Value: 1},
		{Name: "b", Value: 2},
		{Name: "c", Value: 3},
	}

	testWhereStructsFiltered := ext.Where(testWhereStructs, func(d WhereTestStruct) bool {
		return d.Value > 1
	})

	if len(testWhereStructsFiltered) != 2 {
		t.Error("Expected 2, got", len(testWhereStructsFiltered))
	}

	if testWhereStructsFiltered[0].Name != "b" {
		t.Error("Expected b, got", testWhereStructsFiltered[0].Name)
	}

	if testWhereStructsFiltered[1].Name != "c" {
		t.Error("Expected c, got", testWhereStructsFiltered[1].Name)
	}

	testWhereStructsFiltered = ext.Where(testWhereStructs, func(d WhereTestStruct) bool {
		return d.Value > 2
	})

	if len(testWhereStructsFiltered) != 1 {
		t.Error("Expected 1, got", len(testWhereStructsFiltered))
	}

	if testWhereStructsFiltered[0].Name != "c" {
		t.Error("Expected c, got", testWhereStructsFiltered[0].Name)
	}

	testWhereStructsFiltered = ext.Where(testWhereStructs, func(d WhereTestStruct) bool {
		return d.Value > 3
	})

	if len(testWhereStructsFiltered) != 0 {
		t.Error("Expected 0, got", len(testWhereStructsFiltered))
	}

	testWhereStructsFiltered = ext.Where(testWhereStructs, func(d WhereTestStruct) bool {
		return d.Name == "a"
	})

	if len(testWhereStructsFiltered) != 1 {
		t.Error("Expected 1, got", len(testWhereStructsFiltered))
	}

	if testWhereStructsFiltered[0].Name != "a" {
		t.Error("Expected a, got", testWhereStructsFiltered[0].Name)
	}

}
