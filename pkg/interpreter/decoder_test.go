package interpreter

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

func TestGetValueType(t *testing.T) {
	tests := map[string]struct {
		token         string
		userInput     string
		expectedValue string
		expectedType  parameterType
		expectedErr   error
	}{"Test Decode string literal token": {
		token:         "\"This is String Literal\"",
		userInput:     "",
		expectedValue: "This is String Literal",
		expectedType:  STRING,
		expectedErr:   nil,
	}, "Test Decode numeric literal token": {
		token:         "120",
		userInput:     "",
		expectedValue: "120",
		expectedType:  INTEGER,
		expectedErr:   nil,
	}, "Test Decode user input token": {
		token:         "input",
		userInput:     "data",
		expectedValue: "data",
		expectedType:  STRING,
		expectedErr:   nil,
	},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			i := NewInterpreter(&memory.MemoryManager{})

			reader := bufio.NewReader(strings.NewReader(test.userInput))
			actualValue, actualType, err := i.decoder.getValueType(test.token, reader)
			if err != test.expectedErr {
				t.Fatalf("Unexpected Error Mismatch expected %q found %q\n", test.expectedErr, err)
			}

			if actualType != test.expectedType {
				t.Fatalf("Unexpected Types Mismatch expected %q found %q\n", test.expectedType, actualType)
			}
			if actualValue != test.expectedValue {
				t.Fatalf("Unexpected Value Mismatch expected %q found %q\n", test.expectedValue, actualValue)
			}
		})
	}

}
func TestIsSymbol(t *testing.T) {

	tests := map[string]struct {
		token    string
		expected bool
	}{
		"numeric literal": {token: "120", expected: false},
		"string literal":  {token: "\"String Literal Value\"", expected: false},
		"correct symbol":  {token: "variable_name", expected: true},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			i := NewInterpreter(&memory.MemoryManager{})
			actual := i.decoder.isSymbol(test.token)
			if test.expected != actual {
				t.Fatalf("Unexpected result expected %t found %t\n", test.expected, actual)
			}
		})
	}

}

func TestAllocateIfNotDefined(t *testing.T) {
	tests := map[string]struct {
		symbol            string
		inputSymbolTable  symbolTable
		outputSymbolTable symbolTable
	}{
		"Test Undefined Symbol in Empty Symbol Table": {
			symbol:            "x",
			inputSymbolTable:  symbolTable{},
			outputSymbolTable: symbolTable{"x": 0},
		},
		"Test Already Existing Symbol": {
			symbol:            "x",
			inputSymbolTable:  symbolTable{"x": 0},
			outputSymbolTable: symbolTable{"x": 0},
		},
		"Test Undefined Symbol in nonEmpty Symbol Table": {
			symbol:            "x",
			inputSymbolTable:  symbolTable{"y": 0},
			outputSymbolTable: symbolTable{"y": 0, "x": 1},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			i := NewInterpreter(&memory.MemoryManager{})
			symbol := test.symbol
			symTable := test.inputSymbolTable
			i.decoder.allocateIfNotDefined(symbol, symTable)
			if !reflect.DeepEqual(symTable, test.outputSymbolTable) {
				t.Fatalf("Unexpected Mismatch expected %q found %q\n", test.outputSymbolTable, symTable)
			}
		})
	}
}
