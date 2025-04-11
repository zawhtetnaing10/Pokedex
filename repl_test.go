package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " Hello World ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Scyther Greninja",
			expected: []string{"scyther", "greninja"},
		},
		{
			input:    "THIS is Insane",
			expected: []string{"this", "is", "insane"},
		},
		{
			input:    "Flareon Vaporeon Jolteon",
			expected: []string{"flareon", "vaporeon", "jolteon"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("The actual result doesn't have the same length as expected")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Words don't match expected: %v, actual: %v", expectedWord, word)
			}
		}
	}
}
