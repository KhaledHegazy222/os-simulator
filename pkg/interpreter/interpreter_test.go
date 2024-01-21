package interpreter

import (
	"testing"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

func TestMatchCommand(t *testing.T) {
	t.Run("Test Match Existing Command", func(t *testing.T) {
		i := NewInterpreter(memory.MemoryManager{})
		expected := availableCommands["assign"]
		actual, err := i.matchCommand(Instruction{Command: "assign", Args: []string{"x", "11"}})

		if err != nil {
			t.Fatalf("Unexpected Error %q\n", err)
		}

		if expected.command != actual.command {
			t.Fatalf("Unexpected mismatch: expected %q, found %q\n", expected.command, actual.command)
		}
	})
	t.Run("Test Match Existing Command With Insufficient Args Number", func(t *testing.T) {
		i := NewInterpreter(memory.MemoryManager{})
		expected := allowedCommand{}
		actual, err := i.matchCommand(Instruction{Command: "assign", Args: []string{"x"}})

		if err != errInsufficientArguments {
			t.Fatalf("Expected %q, Found %q\n", errInsufficientArguments, err)
		}

		if expected.command != actual.command {
			t.Fatalf("Unexpected mismatch: expected %q, found %q\n", expected.command, actual.command)
		}

	})
	t.Run("Test Invalid Command", func(t *testing.T) {
		i := NewInterpreter(memory.MemoryManager{})
		expected := allowedCommand{}
		actual, err := i.matchCommand(Instruction{Command: "RandomCommand", Args: []string{"x"}})

		if err != errInvalidCommand {
			t.Fatalf("Expected %q, Found %q\n", errInsufficientArguments, err)
		}

		if expected.command != actual.command {
			t.Fatalf("Unexpected mismatch: expected %q, found %q\n", expected.command, actual.command)
		}

	})
}
