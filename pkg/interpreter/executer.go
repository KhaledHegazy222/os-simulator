package interpreter

import (
	"strconv"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
	"github.com/KhaledHegazy222/os-simulator/pkg/systemcalls"
)

// parameterType represents the type of a command parameter.
type parameterType int8

// statusCode represents the status code returned after executing a command.
type statusCode int8

const (
	// INTEGER represents the parameter type for integer values.
	INTEGER parameterType = 0
	// STRING represents the parameter type for string values.
	STRING parameterType = 1
	// ANY represents the parameter type for any data type.
	ANY parameterType = 2
)

const (
	// SUCCESS represents the successful status code after command execution.
	SUCCESS statusCode = 0
	// ERROR represents the error status code after command execution.
	ERROR statusCode = 1
)

// allowedCommand represents a command along with its expected parameters and execution function.
type allowedCommand struct {
	command    string
	parameters []parameterType
	run        func(instruction Instruction, process *memory.PCB) statusCode
}

// availableCommands is a map of command names to their corresponding allowedCommand structure.
var availableCommands = map[string]allowedCommand{
	"assign":      {command: "assign", parameters: []parameterType{INTEGER, ANY}, run: runAssign},
	"print":       {command: "print", parameters: []parameterType{ANY}, run: runPrint},
	"semWait":     {command: "semWait", parameters: []parameterType{STRING}, run: runSemWait},
	"semSignal":   {command: "semSignal", parameters: []parameterType{STRING}, run: runSemSignal},
	"writeFile":   {command: "writeFile", parameters: []parameterType{STRING, ANY}, run: runWriteFile},
	"readFile":    {command: "readFile", parameters: []parameterType{STRING}, run: runReadFile},
	"printFromTo": {command: "printFromTo", parameters: []parameterType{INTEGER, INTEGER}, run: runPrintFromTo},
}

// runAssign executes the "assign" command, setting the value of a variable in the process's memory.
func runAssign(instruction Instruction, process *memory.PCB) statusCode {
	destinationAddress, err := strconv.Atoi(instruction.Args[0])
	if err != nil {
		return ERROR
	}
	value, err := strconv.Atoi(instruction.Args[1])
	if err != nil {
		return ERROR
	}
	process.SetDataWord(destinationAddress, value)
	return SUCCESS
}

// runPrint executes the "print" command, printing the specified data to standard output.
func runPrint(instruction Instruction, process *memory.PCB) statusCode {
	os := systemcalls.NewOS()
	data := instruction.Args[0]
	os.PrintToStdOut(data)
	return SUCCESS
}

// runSemWait executes the "semWait" command, representing a semaphore wait operation.
func runSemWait(instruction Instruction, process *memory.PCB) statusCode {
	return SUCCESS
}

// runSemSignal executes the "semSignal" command, representing a semaphore signal operation.
func runSemSignal(instruction Instruction, process *memory.PCB) statusCode {
	return SUCCESS
}

// runWriteFile executes the "writeFile" command, writing data to a specified file.
func runWriteFile(instruction Instruction, process *memory.PCB) statusCode {
	os := systemcalls.NewOS()
	path, data := instruction.Args[0], instruction.Args[1]

	err := os.WriteToFile(path, data)
	if err != nil {
		return ERROR
	}
	return SUCCESS
}

// runReadFile executes the "readFile" command, reading data from a specified file.
func runReadFile(instruction Instruction, process *memory.PCB) statusCode {
	os := systemcalls.NewOS()
	path := instruction.Args[0]
	_, err := strconv.Atoi(instruction.Args[1])
	if err != nil {
		return ERROR
	}

	_, err = os.ReadFile(path)
	if err != nil {
		return ERROR
	}
	return SUCCESS
}

// runPrintFromTo executes the "printFromTo" command, printing a range of values.
func runPrintFromTo(instruction Instruction, process *memory.PCB) statusCode {
	os := systemcalls.NewOS()
	start, err := strconv.Atoi(instruction.Args[0])
	if err != nil {
		return ERROR
	}
	end, err := strconv.Atoi(instruction.Args[1])
	if err != nil {
		return ERROR
	}
	for i := start; i <= end; i++ {
		os.PrintToStdOut(strconv.Itoa(i))
	}
	return SUCCESS
}
