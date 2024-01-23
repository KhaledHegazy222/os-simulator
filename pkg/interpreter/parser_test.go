package interpreter

import (
	"reflect"
	"testing"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

func TestParser(t *testing.T) {

	t.Run("Testing Single Command no args", func(t *testing.T) {
		i := NewInterpreter(memory.MemoryManager{})
		actual := i.parse("test")
		expected := Instruction{
			Command: "test", Args: []string{},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("Unmatched Result: expected %q, found %q\n", expected, actual)
		}

	})

	t.Run("Testing Multi Command multi args", func(t *testing.T) {
		i := NewInterpreter(memory.MemoryManager{})
		actual := i.parse("assign x 1")
		expected := Instruction{
			Command: "assign", Args: []string{"x", "1"},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("Unmatched Result: expected %q, found %q\n", expected, actual)
		}

	})

	t.Run("Testing String Literal Args with no spaces", func(t *testing.T) {
		i := NewInterpreter(memory.MemoryManager{})
		actual := i.parse("assign x \"string_content\"")
		expected := Instruction{
			Command: "assign", Args: []string{"x", "\"string_content\""},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("Unmatched Result: expected %q, found %q\n", expected, actual)
		}

	})

	t.Run("Testing String Literal Args with spaces", func(t *testing.T) {
		i := NewInterpreter(memory.MemoryManager{})
		actual := i.parse("assign x \"string content test\"")
		expected := Instruction{
			Command: "assign", Args: []string{"x", "\"string content test\""},
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("Unmatched Result: expected %q, found %q\n", expected, actual)
		}

	})

}
