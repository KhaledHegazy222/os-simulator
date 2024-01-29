package interpreter

import (
	"testing"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

func TestMatchCommand(t *testing.T) {
	t.Run("Test Match Existing Command", func(t *testing.T) {
		i := NewInterpreter(&memory.MemoryManager{})
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
		i := NewInterpreter(&memory.MemoryManager{})
		expected := allowedCommand{}
		actual, err := i.matchCommand(Instruction{Command: "assign", Args: []string{"x"}})

		if err != ErrInsufficientArguments {
			t.Fatalf("Expected %q, Found %q\n", ErrInsufficientArguments, err)
		}

		if expected.command != actual.command {
			t.Fatalf("Unexpected mismatch: expected %q, found %q\n", expected.command, actual.command)
		}

	})
	t.Run("Test Invalid Command", func(t *testing.T) {
		i := NewInterpreter(&memory.MemoryManager{})
		expected := allowedCommand{}
		actual, err := i.matchCommand(Instruction{Command: "RandomCommand", Args: []string{"x"}})

		if err != ErrInvalidCommand {
			t.Fatalf("Expected %q, Found %q\n", ErrInsufficientArguments, err)
		}

		if expected.command != actual.command {
			t.Fatalf("Unexpected mismatch: expected %q, found %q\n", expected.command, actual.command)
		}

	})
}

func TestMatchTypes(t *testing.T) {
	t.Run("Test Matching Command types", func(t *testing.T) {

		i := Interpreter{}
		instruction := &Instruction{
			Command: "assign",
			Args:    []string{"42", "\"string data\""},
		}
		command := allowedCommand{
			command:    "assign",
			parameters: []parameterType{INTEGER, ANY},
			run:        nil,
		}

		err := i.matchTypes(instruction, command)
		if err != nil {
			t.Errorf("Error: expected nil, got %v", err)
		}

	})
	t.Run("Test Mismatching Command types", func(t *testing.T) {
		i := Interpreter{}
		instruction := &Instruction{
			Command: "assign",
			Args:    []string{"42", "32"},
		}
		command := allowedCommand{
			command:    "assign",
			parameters: []parameterType{STRING, INTEGER},
			run:        nil,
		}

		err := i.matchTypes(instruction, command)
		if err != ErrInvalidArgumentType {
			t.Errorf("Error: expected %v, got %v", ErrInvalidArgumentType, err)
		}

	})

}

func TestTypeCheck(t *testing.T) {
	tests := map[string]struct {
		tokenType   parameterType
		checkedType parameterType
		result      bool
	}{
		"integer integer match": {tokenType: INTEGER, checkedType: INTEGER, result: true},
		"integer any match":     {tokenType: INTEGER, checkedType: ANY, result: true},
		"any Integer match":     {tokenType: ANY, checkedType: INTEGER, result: false},
		"string integer match":  {tokenType: STRING, checkedType: INTEGER, result: false},
		"any any match":         {tokenType: ANY, checkedType: ANY, result: true},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			i := Interpreter{}
			actual := i.typeCheck(test.tokenType, test.checkedType)
			if test.result != actual {
				t.Fatalf("Unexpected Mismatch expected %t, found %t\n", test.result, actual)
			}
		})
	}
}
