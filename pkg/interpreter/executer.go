package interpreter

import (
	"strconv"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
	"github.com/KhaledHegazy222/os-simulator/pkg/systemcalls"
)

type parameterType int8

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

type allowedCommand struct {
	command    string
	parameters []parameterType
	run        func(instruction Instruction, process *memory.PCB) statusCode
}

var availableCommands = map[string]allowedCommand{
	"assign":      {command: "assign", parameters: []parameterType{INTEGER, ANY}, run: runAssign},
	"print":       {command: "print", parameters: []parameterType{ANY}, run: runPrint},
	"semWait":     {command: "semWait", parameters: []parameterType{STRING}, run: runSemWait},
	"semSignal":   {command: "semSignal", parameters: []parameterType{STRING}, run: runSemSignal},
	"writeFile":   {command: "writeFile", parameters: []parameterType{STRING, ANY}, run: runWriteFile},
	"readFile":    {command: "readFile", parameters: []parameterType{STRING}, run: runReadFile},
	"printFromTo": {command: "printFromTo", parameters: []parameterType{INTEGER, INTEGER}, run: runPrintFromTo},
}

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

func runPrint(instruction Instruction, process *memory.PCB) statusCode {
	os := systemcalls.NewOS()
	data := instruction.Args[0]
	os.PrintToStdOut(data)
	return SUCCESS
}

func runSemWait(instruction Instruction, process *memory.PCB) statusCode {
	return SUCCESS
}

func runSemSignal(instruction Instruction, process *memory.PCB) statusCode {
	return SUCCESS
}

func runWriteFile(instruction Instruction, process *memory.PCB) statusCode {
	os := systemcalls.NewOS()
	path, data := instruction.Args[0], instruction.Args[1]

	err := os.WriteToFile(path, data)
	if err != nil {
		return ERROR
	}
	return SUCCESS
}

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
