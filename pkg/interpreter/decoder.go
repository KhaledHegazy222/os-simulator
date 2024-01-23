package interpreter

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

// symbolTable is a map representing the symbol table for variables.
type symbolTable map[string]int

// processId is a type representing the unique identifier for a process.
type processId int

var (
	// errType is a common error for a type error during value conversion.
	errType = errors.New("type error")
	// errUndefinedSymbol is a common error for accessing an undefined symbol in the symbol table.
	errUndefinedSymbol = errors.New("undefined symbol")
)

// getSymbolTable returns the symbol table for the given process, initializing it if not present.
func (i *Interpreter) getSymbolTable(process *memory.PCB) symbolTable {
	// if not executed before init Process
	_, isPresent := i.processToSymbolTable[processId(process.Id)]
	if !isPresent {
		i.processToSymbolTable[processId(process.Id)] = symbolTable{}
	}
	symTable, _ := i.processToSymbolTable[processId(process.Id)]
	return symTable
}

// decodeArgs decodes and updates the arguments in the instruction based on the symbol table.
func (i *Interpreter) decodeArgs(instruction *Instruction, process *memory.PCB) error {
	symTable := i.getSymbolTable(process)
	if instruction.Command == "assign" {
		// Allocate the variable if not defined
		i.allocateIfNotDefined(instruction.Args[0], symTable)
		// Replace the destination operand with its address
		instruction.Args[0] = strconv.Itoa(symTable[instruction.Args[0]])
	}
	for index, arg := range instruction.Args {

		if i.isSymbol(arg) {
			address, isPresent := symTable[arg]
			if !isPresent {
				return errUndefinedSymbol
			}
			// get data of address of data
			data, err := process.GetDataWord(address)
			if err != nil {
				return err
			}
			instruction.Args[index] = data
		}
	}

	return nil
}

// getValueType returns the value and type of a token (argument).
func (i *Interpreter) getValueType(token string, reader io.Reader) (value string, valueType parameterType, err error) {
	if len(token) > 2 && strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		croppedToken := token[1 : len(token)-1]
		return croppedToken, STRING, nil
	} else if token == "input" {
		var data string
		fmt.Fscanf(reader, "%s", &data)
		return data, STRING, nil
	} else if _, conversionErr := strconv.Atoi(token); conversionErr == nil {
		return token, INTEGER, nil
	} else {
		return "", ANY, errType
	}
}

// isSymbol checks if the token is a symbol (variable) rather than a literal value.
func (i *Interpreter) isSymbol(token string) bool {
	if token == "input"{
		return false
	}
	if len(token) > 2 && strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		return false
	}
	_, err := strconv.Atoi(token)

	return err != nil

}

// allocateIfNotDefined allocates a new address for a symbol if it is not already defined.
func (i *Interpreter) allocateIfNotDefined(symbol string, symTable symbolTable) {
	_, isPresent := symTable[symbol]
	if !isPresent {
		// Set new address
		symTable[symbol] = len(symTable)
	}
}
