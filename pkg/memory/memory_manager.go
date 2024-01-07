package memory

import (
	"errors"
)

const (
	variablesSize      = 3
	PCBSize            = 6
	memorySize         = 40 + 3
	memoryStartAddress = 1
	memoryEndAddress   = 40
)

var (
	nextId = 1

	NotEnoughSpaceErr    = errors.New("not enough space in the memory")
	ProcessIdNotFoundErr = errors.New("process id is not found")
	InternalMemoryErrorErr = errors.New("internal memory error")
)

// MemoryManager manager for the memory component that controls allocation and de-allocation of processes and
// other information about processes state,PC,location,...
type MemoryManager struct {
	ram             RAMMemory
	processLocation map[int]int
}


func getProcessSize(unparsedCodeSize int) int {
	return unparsedCodeSize + variablesSize + PCBSize
}

// AddProcess creates process and save it in memory
func (m *MemoryManager) AddProcess(unparsedCode []string) (PCB, error) {
	neededSize := getProcessSize(len(unparsedCode))

	for i := memoryStartAddress; i <= memoryEndAddress; i++ {
		isFree := m.ram.isFree(i, i+neededSize)
		if isFree {
			pcb := m.ram.allocateProcess(i, unparsedCode)
			m.processLocation[pcb.Id] = pcb.Start
			return pcb, nil
		}
	}
	return PCB{}, NotEnoughSpaceErr
}

// DeleteProcess deletes process from memory
func (m *MemoryManager) DeleteProcess(id int) error {
	processStartLocation:=m.processLocation[id]

	if processStartLocation == 0 {
		return ProcessIdNotFoundErr
	}

	pcb,err:=m.ram.getProcessPCB(processStartLocation)
	if err != nil {
		return InternalMemoryErrorErr
	}
	
	m.processLocation[id]=0
	pcb.delete()
	return nil
}
