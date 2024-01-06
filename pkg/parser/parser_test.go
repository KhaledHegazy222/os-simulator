package parser

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {

	t.Run("Testing Single Command no args", func(t *testing.T) {
		parser := NewParser()
		instructions, err := parser.Parse("../../scripts/sample1")
		expected := []Instruction{
			{Command: "test", Args: []string{}},
		}
		if err != nil {
			t.Fatalf("Unexpected error: %q\n", err)
		}

		if !reflect.DeepEqual(instructions, expected) {
			t.Fatalf("Unmatched Result: expected %q, found %q\n", expected, instructions)
		}

	})

	t.Run("Testing Multi Command multi args", func(t *testing.T) {
		parser := NewParser()
		instructions, err := parser.Parse("../../scripts/sample2")
		expected := []Instruction{
			{Command: "assign", Args: []string{"x", "1"}},
			{Command: "print", Args: []string{"x"}},
			{Command: "writeFile", Args: []string{"filename", "data"}},
			{Command: "semWait", Args: []string{"mutex1"}},
		}
		if err != nil {
			t.Fatalf("Unexpected error: %q\n", err)
		}

		if !reflect.DeepEqual(instructions, expected) {
			t.Fatalf("Unmatched Result: expected %q, found %q\n", expected, instructions)
		}

	})
}
