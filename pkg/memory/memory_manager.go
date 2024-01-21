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

	NotEnoughSpaceErr      = errors.New("not enough space in the memory")
	ProcessIdNotFoundErr   = errors.New("process id is not found")
	InternalMemoryErrorErr = errors.New("internal memory error")
)

// MemoryManager manager for the memory component that controls allocation and de-allocation of processes and
// other information about processes state,PC,location,...
type MemoryManager struct {
	ram             RAMMemory
	processLocation map[int]int
	numberOfProcesses int
}

// Memory interface that handles addition and deletion of processes in memory
type Memory interface {
	AddProcess(unparsedCode []string) (PCBManager, error)
	DeleteProcess(processId int) error
}

// NewMemoryManager factory method that creates new memory manager
func NewMemoryManager() MemoryManager {
	ram := RAMMemory{}
	return MemoryManager{
		ram:             ram,
		processLocation: make(map[int]int),
		numberOfProcesses: 0,
	}
}

func getProcessSize(unparsedCodeSize int) int {
	return PCBSize + unparsedCodeSize + variablesSize
}

func (m *MemoryManager) getNextID() int {
	m.numberOfProcesses++
	return m.numberOfProcesses
}

// AddProcess creates process and save it in memory
func (m *MemoryManager) AddProcess(unparsedCode []string) (PCB, error) {
	neededSize := getProcessSize(len(unparsedCode))

	for i := memoryStartAddress; i <= memoryEndAddress; i++ {
		isFree := m.ram.isFree(i, i+neededSize)
		if isFree {
			pcb := m.ram.allocateProcess(i, unparsedCode,m.numberOfProcesses)
			m.processLocation[pcb.Id] = pcb.Start
			return pcb, nil
		}
	}
	return PCB{}, NotEnoughSpaceErr
}

// DeleteProcess deletes process from memory
func (m *MemoryManager) DeleteProcess(processId int) error {
	processStartLocation := m.processLocation[processId]

	if processStartLocation == 0 {
		return ProcessIdNotFoundErr
	}

	pcb, err := m.ram.getProcessPCB(processStartLocation)
	if err != nil {
		return InternalMemoryErrorErr
	}

	delete(m.processLocation, processId)
	pcb.delete()
	return nil
}
