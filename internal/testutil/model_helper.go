package testutil

import (
	"testing"

	"github.com/IgorBayerl/AdlerCov/internal/model"
)

// FindMethod searches a slice of methods for one with a specific name.
// It calls t.Fatalf if the method is not found.
func FindMethod(t *testing.T, methods []model.Method, name string) model.Method {
	t.Helper() // Marks this function as a test helper.
	for _, m := range methods {
		if m.Name == name {
			return m
		}
	}
	t.Fatalf("method with name '%s' not found in provided slice", name)
	return model.Method{} // Unreachable, but satisfies the compiler.
}

// FindLine searches a slice of lines for one with a specific line number.
// It calls t.Fatalf if the line is not found.
func FindLine(t *testing.T, lines []model.Line, number int) model.Line {
	t.Helper() // Marks this function as a test helper.
	for _, l := range lines {
		if l.Number == number {
			return l
		}
	}
	t.Fatalf("line with number '%d' not found in provided slice", number)
	return model.Line{} // Unreachable, but satisfies the compiler.
}
