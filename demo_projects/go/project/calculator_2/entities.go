package calculator_2

import "fmt"

// Counter is a simple struct to demonstrate methods with receivers.
type Counter struct {
	value int
}

// Increment increases the counter's value by one (pointer receiver).
func (c *Counter) Increment() {
	c.value++
}

// Value returns the current value of the counter (value receiver).
func (c Counter) Value() int {
	return c.value
}

// Add adds a given amount to the counter.
// This method has a branch to test conditional coverage.
func (c *Counter) Add(amount int) {
	if amount < 0 {
		// Do not add negative numbers; this branch will be tested.
		return
	}
	c.value += amount
}

// MessageBuilder is another example struct for testing.
type MessageBuilder struct {
	greeting string
}

// NewMessageBuilder creates a new MessageBuilder.
// This function will not be covered by tests to show uncovered code.
func NewMessageBuilder(greeting string) *MessageBuilder {
	if greeting == "" {
		return &MessageBuilder{greeting: "Hello"}
	}
	return &MessageBuilder{greeting: greeting}
}

// Greet returns a formatted greeting string.
func (mb *MessageBuilder) Greet(name string) string {
	if name == "" {
		return fmt.Sprintf("%s, World!", mb.greeting)
	}
	return fmt.Sprintf("%s, %s!", mb.greeting, name)
}

// unexportedHelper is an unexported function to see if it appears in coverage.
// It should not be part of the public API but will be covered by tests.
func unexportedHelper() string {
	return "I am a helper"
}