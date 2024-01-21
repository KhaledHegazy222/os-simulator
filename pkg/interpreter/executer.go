package interpreter

import (
	"fmt"
	"strconv"
)

type parameterType int8

const (
	INTEGER parameterType = 0
	STRING  parameterType = 1
	ANY     parameterType = 2
)

type allowedCommand struct {
	command    string
	parameters []parameterType
	run        func(instruction Instruction)
}

var availableCommands = map[string]allowedCommand{
	"assign":      {command: "assign", parameters: []parameterType{ANY, ANY}, run: runAssign},
	"print":       {command: "print", parameters: []parameterType{ANY}, run: runPrint},
	"semWait":     {command: "semWait", parameters: []parameterType{STRING}, run: runSemWait},
	"semSignal":   {command: "semSignal", parameters: []parameterType{STRING}, run: runSemSignal},
	"writeFile":   {command: "writeFile", parameters: []parameterType{STRING, ANY}, run: runWriteFile},
	"readFile":    {command: "readFile", parameters: []parameterType{STRING}, run: runReadFile},
	"printFromTo": {command: "printFromTo", parameters: []parameterType{INTEGER, INTEGER}, run: runPrintFromTo},
}

func runAssign(instruction Instruction) {

}
func runPrint(instruction Instruction) {
	printedValue := instruction.Args[0]
	fmt.Println(printedValue)
}
func runSemWait(instruction Instruction) {

}
func runSemSignal(instruction Instruction) {

}
func runWriteFile(instruction Instruction) {

}
func runReadFile(instruction Instruction) {
	
}
func runPrintFromTo(instruction Instruction) {
	start, _ := strconv.Atoi(instruction.Args[0])
	end, _ := strconv.Atoi(instruction.Args[1])
	for i := start; i <= end; i++ {
		fmt.Println(i)
	}
}
