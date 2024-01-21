package interpreter

import (
	"errors"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

type Interpreter struct {
	memory               *memory.MemoryManager
	processToSymbolTable map[processId]symbolTable
}

type Instruction struct {
	Command string
	Args    []string
}

var errBlockedProcess = errors.New("the process is currently blocked and can't execute")
var errInvalidCommand = errors.New("the process is currently blocked and can't execute")
var errInsufficientArguments = errors.New("the process is currently blocked and can't execute")
var errInvalidArgumentType = errors.New("the process is currently blocked and can't execute")

func NewInterpreter(memoryManager memory.MemoryManager) Interpreter {
	return Interpreter{
		memory:               &memoryManager,
		processToSymbolTable: map[processId]symbolTable{},
	}
}

func (i *Interpreter) Execute(process memory.PCB) error {

	// Return if Blocked
	if process.State == memory.Blocked {
		return errBlockedProcess
	}
	// Get the next Instruction
	nextLine, err := process.GetNextInstruction()
	if err != nil {
		return err
	}

	// Parse Instruction
	instruction := i.parse(nextLine)

	// Find Matched Command
	command, err := i.matchCommand(instruction)
	if err != nil {
		return err
	}

	// Decode Instruction arguments
	symTable := i.getSymbolTable(process)
	err = i.decodeArgs(&instruction, command, symTable)
	if err != nil {
		return err
	}

	// execute Instruction
	// command.run(instruction)

	return nil

}

func (i *Interpreter) matchCommand(instruction Instruction) (allowedCommand, error) {
	matchedCommand, isPresent := availableCommands[instruction.Command]
	if !isPresent {
		return allowedCommand{}, errInvalidCommand
	}

	if len(matchedCommand.parameters) != len(instruction.Args) {
		return allowedCommand{}, errInsufficientArguments
	}

	return matchedCommand, nil
}
