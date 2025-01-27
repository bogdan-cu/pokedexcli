package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		output := cleanInput(c.input)
		for idx := range output {
			actual := output[idx]
			expected := c.expected[idx]
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
				t.Fail()
			}
		}
	}
}
