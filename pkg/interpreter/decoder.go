package interpreter

import (
	"errors"
	"strconv"
	"strings"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
	"github.com/KhaledHegazy222/os-simulator/pkg/systemcalls"
)

// symbolTable is a map representing the symbol table for variables.
type symbolTable map[string]int

// processId is a type representing the unique identifier for a process.
type processId int

type decoderManager struct {
	processToSymbolTable map[processId]symbolTable
	os                   *systemcalls.OS
}

var (
	// ErrType is a common error for a type error during value conversion.
	ErrType = errors.New("type error")
	// ErrUndefinedSymbol is a common error for accessing an undefined symbol in the symbol table.
	ErrUndefinedSymbol = errors.New("undefined symbol")
)

func (d *decoderManager) getSymbolTable(process *memory.PCB) symbolTable {
	// if not executed before init Process
	_, isPresent := d.processToSymbolTable[processId(process.Id)]
	if !isPresent {
		d.processToSymbolTable[processId(process.Id)] = symbolTable{}
	}
	symTable := d.processToSymbolTable[processId(process.Id)]
	return symTable
}

func (d *decoderManager) decodeArgs(instruction *Instruction, process *memory.PCB) error {
	symTable := d.getSymbolTable(process)
	if instruction.Command == "assign" {
		// Allocate the variable if not defined
		d.allocateIfNotDefined(instruction.Args[0], symTable)
		// Replace the destination operand with its address
		instruction.Args[0] = strconv.Itoa(symTable[instruction.Args[0]])
	}
	for index, arg := range instruction.Args {

		if d.isSymbol(arg) {
			address, isPresent := symTable[arg]
			if !isPresent {
				return ErrUndefinedSymbol
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

func (d *decoderManager) getValueType(token string) (value string, valueType parameterType, err error) {
	if len(token) > 2 && strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		croppedToken := token[1 : len(token)-1] 
		return croppedToken, STRING, nil
	} else if token == "input" {
		data:=d.os.GetInput()
		return data, STRING, nil
	} else if _, conversionErr := strconv.Atoi(token); conversionErr == nil {
		return token, INTEGER, nil
	} else {
		return "", ANY, ErrType
	}
}

func (d *decoderManager) isSymbol(token string) bool {
	if token == "input" {
		return false
	}
	if len(token) > 2 && strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		return false
	}
	_, err := strconv.Atoi(token)

	return err != nil

}

func (d *decoderManager) allocateIfNotDefined(symbol string, symTable symbolTable) {
	_, isPresent := symTable[symbol]
	if !isPresent {
		// Set new address
		symTable[symbol] = len(symTable)
	}
}
