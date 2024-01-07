package memory

import (
	"errors"
	"fmt"
)

type STATE string

const (
	Running    STATE = "running"
	Terminated STATE = "terminated"
	Ready      STATE = "ready"
	Blocked    STATE = "blocked"

	variablesSize = 3
	PCBSize       = 6
	memorySize    = 40
)

var (
	nextId = 0

	NotEnoughSpaceErr    = errors.New("not enough space in the memory")
	EndOfInstructionsErr = errors.New("reached end of instructions")
	ProtectionErr        = errors.New("not allowed to access that part of memory")

	RAMMemory [memorySize]string
)


func getNextID() int {
	nextId++
	return nextId
}

func getProcessSize(unparsedCodeSize int) int {
	return unparsedCodeSize + variablesSize + PCBSize
}

func isRamFree(from int, to int) (bool) {
	if from >= 40 || from < 0 || to >= 40 || to < 0 {
		return false
	}

	for i := from; i <= to; i++ {
		if RAMMemory[i] != "" {
			return false
		}
	}
	return true
}

func allocateProcess(start int, unparsedCode []string) (PCB) {
	end := start + getProcessSize(len(unparsedCode))

	pcb := PCB{
		Id:       getNextID(),
		State:    Running,
		PC:       start + PCBSize,
		Start:    start,
		End:      end,
		codeSize: len(unparsedCode),
	}

	// allocate pcb in the first five words
	RAMMemory[start] = fmt.Sprint(pcb.Id)
	RAMMemory[start+1] = fmt.Sprint(pcb.State)
	RAMMemory[start+2] = fmt.Sprint(pcb.PC)
	RAMMemory[start+3] = fmt.Sprint(pcb.Start)
	RAMMemory[start+4] = fmt.Sprint(pcb.End)
	RAMMemory[start+5] = fmt.Sprint(pcb.codeSize)

	// allocate unparsed code
	unparsedCodeStart := start + PCBSize
	for i := 0; i < len(unparsedCode); i++ {
		RAMMemory[unparsedCodeStart+i] = unparsedCode[i]
	}

	// allocate variables in the last three words. initial value is zero
	variablesStart := unparsedCodeStart + len(unparsedCode)
	RAMMemory[variablesStart] = fmt.Sprint(0)
	RAMMemory[variablesStart+1] = fmt.Sprint(0)
	RAMMemory[variablesStart+2] = fmt.Sprint(0)

	return pcb
}

func AddProcess(unparsedCode []string) (PCB, error) {
	neededSize := len(unparsedCode) + PCBSize + variablesSize

	for i := 0; i < 40; i++ {
		isFree := isRamFree(i, i+neededSize)
		if isFree {
			return allocateProcess(i, unparsedCode), nil
		}
	}
	return PCB{}, NotEnoughSpaceErr
}

