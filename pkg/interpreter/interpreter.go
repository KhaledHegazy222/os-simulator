// Package interpreter provides an interpreter for processing instructions
// in the context of an operating system simulation.
package interpreter

import (
	"errors"
	"os"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

// Interpreter represents the interpreter for processing instructions.
type Interpreter struct {
	memory               *memory.MemoryManager
	processToSymbolTable map[processId]symbolTable
	decoder              *decoderManager
	parser               *parserManager
}

// Instruction represents a single instruction with a command and its arguments.
type Instruction struct {
	Command string
	Args    []string
}

var (
	// Common error for a blocked process.
	ErrBlockedProcess = errors.New("the process is currently blocked and can't execute")
	// Common error for an invalid command.
	ErrInvalidCommand = errors.New("invalid command")
	// Common error for insufficient arguments for a command.
	ErrInsufficientArguments = errors.New("insufficient arguments for the command")
	// Common error for an invalid argument type.
	ErrInvalidArgumentType = errors.New("invalid argument type")
	// Common error for a runtime error during instruction execution.
	ErrRunTimeError = errors.New("runtime error")
)

// NewInterpreter creates a new Interpreter instance with the provided memory manager.
func NewInterpreter(memoryManager *memory.MemoryManager) Interpreter {
	processToSymbolTable := map[processId]symbolTable{}
	decoder := &decoderManager{processToSymbolTable: processToSymbolTable}
	parser := &parserManager{}
	return Interpreter{
		memory:               memoryManager,
		processToSymbolTable: processToSymbolTable,
		decoder:              decoder,
		parser:               parser,
	}
}

// Execute executes the next instruction for the given process.
func (i *Interpreter) Execute(process *memory.PCB) error {
	// Return if Blocked
	if process.State == memory.Blocked {
		return ErrBlockedProcess
	}
	// Get the next Instruction
	nextLine, err := process.GetNextInstruction()
	if err != nil {
		return err
	}

	// Parse Instruction
	instruction := i.parser.parse(nextLine)

	// Find Matched Command
	command, err := i.matchCommand(instruction)
	if err != nil {
		return err
	}

	// Decode Instruction arguments
	if err = i.decoder.decodeArgs(&instruction, process); err != nil {
		return err
	}

	if err = i.matchTypes(&instruction, command); err != nil {
		return err
	}

	// Execute Instruction
	status := command.run(instruction, process)
	if status != SUCCESS {
		return ErrRunTimeError
	}

	process.IncrementPC()
	return nil
}

func (i *Interpreter) matchCommand(instruction Instruction) (allowedCommand, error) {
	matchedCommand, isPresent := availableCommands[instruction.Command]
	if !isPresent {
		return allowedCommand{}, ErrInvalidCommand
	}

	if len(matchedCommand.parameters) != len(instruction.Args) {
		return allowedCommand{}, ErrInsufficientArguments
	}

	return matchedCommand, nil
}

func (i *Interpreter) matchTypes(instruction *Instruction, command allowedCommand) error {
	for index, arg := range instruction.Args {
		value, valueType, err := i.decoder.getValueType(arg, os.Stdin)
		if err != nil {
			return err
		}
		instruction.Args[index] = value
		if !i.typeCheck(valueType, command.parameters[index]) {
			return ErrInvalidArgumentType
		}
	}
	return nil
}

func (i *Interpreter) typeCheck(tokenType parameterType, checkedType parameterType) bool {
	if checkedType == ANY {
		return true
	}
	return tokenType == checkedType
}
