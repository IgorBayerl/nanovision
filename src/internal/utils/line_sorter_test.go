package utils

import (
	"reflect"
	"testing"
)

// testItem implements SortableByLineAndName for testing.
type testItem struct {
	line int
	name string
}

func (t testItem) GetFirstLine() int       { return t.line }
func (t testItem) GetSortableName() string { return t.name }

func TestSortByLineAndName_BasicSorting(t *testing.T) {
	items := []testItem{
		{line: 10, name: "b"},
		{line: 5, name: "c"},
		{line: 5, name: "a"},
		{line: 0, name: "z"},
		{line: 10, name: "a"},
	}

	expected := []testItem{
		{line: 5, name: "a"},
		{line: 5, name: "c"},
		{line: 10, name: "a"},
		{line: 10, name: "b"},
		{line: 0, name: "z"},
	}

	SortByLineAndName(items)

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Expected sorted items: %+v, got: %+v", expected, items)
	}
}

func TestSortByLineAndName_AllZeroLines(t *testing.T) {
	items := []testItem{
		{line: 0, name: "b"},
		{line: 0, name: "a"},
		{line: 0, name: "c"},
	}

	expected := []testItem{
		{line: 0, name: "a"},
		{line: 0, name: "b"},
		{line: 0, name: "c"},
	}

	SortByLineAndName(items)

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Expected sorted items: %+v, got: %+v", expected, items)
	}
}

func TestSortByLineAndName_EmptySlice(t *testing.T) {
	var items []testItem
	SortByLineAndName(items)
	if len(items) != 0 {
		t.Errorf("Expected empty slice, got: %+v", items)
	}
}

func TestSortByLineAndName_SameLineDifferentNames(t *testing.T) {
	items := []testItem{
		{line: 7, name: "delta"},
		{line: 7, name: "alpha"},
		{line: 7, name: "charlie"},
	}

	expected := []testItem{
		{line: 7, name: "alpha"},
		{line: 7, name: "charlie"},
		{line: 7, name: "delta"},
	}

	SortByLineAndName(items)

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Expected sorted items: %+v, got: %+v", expected, items)
	}
}
