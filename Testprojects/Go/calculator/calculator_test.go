package calculator

import "testing"

func TestAdd(t *testing.T) {
	if Add(1, 2) != 3 {
		t.Error("Expected 1 + 2 to equal 3")
	}
	if Add(-1, 1) != 0 {
		t.Error("Expected -1 + 1 to equal 0")
	}
}

func TestSubtract(t *testing.T) {
	if Subtract(3, 2) != 1 {
		t.Error("Expected 3 - 2 to equal 1")
	}
}

func TestMultiplyZero(t *testing.T) {
	if Multiply(0, 5) != 0 {
		t.Error("Expected 0 * 5 to equal 0")
	}
	if Multiply(5, 0) != 0 {
		t.Error("Expected 5 * 0 to equal 0")
	}
	// Note: We are not testing the Multiply(a,b) case where a and b are non-zero.
}

func TestDivide(t *testing.T) {
	quotient, remainder := Divide(10, 3)
	if quotient != 3 || remainder != 1 {
		t.Errorf("Expected 10 / 3 to be 3 with remainder 1, got %d and %d", quotient, remainder)
	}
	// Note: We are not testing the Divide(a, 0) case.
}

// --- Tests for entities.go ---

func TestCounter(t *testing.T) {
	// Test pointer receiver method
	c := &Counter{}
	if c.Value() != 0 {
		t.Errorf("Expected initial counter value to be 0, got %d", c.Value())
	}
	c.Increment()
	if c.Value() != 1 {
		t.Errorf("Expected counter value to be 1 after Increment, got %d", c.Value())
	}

	// Test method with a branch
	c.Add(5)
	if c.Value() != 6 {
		t.Errorf("Expected counter value to be 6 after Add(5), got %d", c.Value())
	}

	// Test the other branch of the Add method
	c.Add(-10)
	if c.Value() != 6 {
		t.Errorf("Expected counter value to remain 6 after Add(-10), got %d", c.Value())
	}
}

func TestMessageBuilder_Greet(t *testing.T) {
	// Note: We are not testing NewMessageBuilder, only the Greet method.
	mb := &MessageBuilder{greeting: "Hi"} // Manually create struct

	greeting := mb.Greet("Alice")
	expected := "Hi, Alice!"
	if greeting != expected {
		t.Errorf("Expected Greet to return '%s', got '%s'", expected, greeting)
	}

	// Test the empty name branch
	greetingWorld := mb.Greet("")
	expectedWorld := "Hi, World!"
	if greetingWorld != expectedWorld {
		t.Errorf("Expected Greet with empty name to return '%s', got '%s'", expectedWorld, greetingWorld)
	}
}

func TestUnexportedHelper(t *testing.T) {
	// This test exists solely to ensure the unexported function is covered
	// and appears in the coverage profile.
	if unexportedHelper() != "I am a helper" {
		t.Error("unexportedHelper returned unexpected value")
	}
}

// --- Test for Cyclomatic Complexity ---

func TestGetGradeForScore(t *testing.T) {
	testCases := []struct {
		score    int
		expected string
	}{
		{101, "Invalid Score"}, // Test upper bound
		// Note: The score < 0 case is intentionally NOT tested
		{95, "A"},
		{90, "A"},
		{85, "B"},
		// Note: The "C" grade (70-79) is intentionally NOT tested
		{65, "D"},
		{59, "F"},
		{0, "F"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			if got := GetGradeForScore(tc.score); got != tc.expected {
				t.Errorf("GetGradeForScore(%d) = %s; want %s", tc.score, got, tc.expected)
			}
		})
	}
}