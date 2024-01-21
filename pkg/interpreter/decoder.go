package interpreter

import (
	"errors"
	"strconv"
	"strings"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

type symbolTable map[string]int
type processId int

var errType = errors.New("Type error")
var errUndefinedSymbol = errors.New("Undefined Symbol")

func (i *Interpreter) getSymbolTable(process memory.PCB) symbolTable {
	// if not executed before init Process
	_, isPresent := i.processToSymbolTable[processId(process.Id)]
	if !isPresent {
		i.processToSymbolTable[processId(process.Id)] = symbolTable{}
	}
	symTable, _ := i.processToSymbolTable[processId(process.Id)]
	return symTable
}

func (i *Interpreter) decodeArgs(instruction *Instruction, process memory.PCB) error {
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

func (i *Interpreter) getValueType(token string) (value string, valueType parameterType, err error) {
	if len(token) > 2 && strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		croppedToken := token[1 : len(token)-1]
		return croppedToken, STRING, nil
	} else if token == "input" {
		return "input", STRING, nil
	} else if _, conversionErr := strconv.Atoi(token); conversionErr == nil {
		return token, INTEGER, nil
	} else {
		return "", ANY, errType
	}

}
func (i *Interpreter) isSymbol(token string) bool {
	if len(token) > 2 && strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		return false
	}
	_, err := strconv.Atoi(token)

	return err != nil

}

func (i *Interpreter) allocateIfNotDefined(symbol string, symTable symbolTable) {
	_, isPresent := symTable[symbol]
	if !isPresent {
		// Set new address
		symTable[symbol] = len(symTable)
	}
}
