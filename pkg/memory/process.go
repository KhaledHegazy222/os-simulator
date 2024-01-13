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
)

var (
	EndOfInstructionsErr = errors.New("reached end of instructions")
	ProtectionErr        = errors.New("not allowed to access that part of memory")
)

type PCB struct {
	Id       int
	State    STATE
	PC       int
	Start    int
	End      int
	CodeSize int
	ram      *RAMMemory
}

func (p *PCB) getPCBAddress() int {
	return p.Start
}

func (p *PCB) getUnparsedCodeAddress() int {
	return p.getPCBAddress() + PCBSize
}

func (p *PCB) getVariablesAddress() int {
	return p.getUnparsedCodeAddress() + p.CodeSize
}

func (p *PCB) delete() {
	for i := p.Start; i <= p.End; i++ {
		p.ram[i] = ""
	}
}

// GetNextInstruction retrieve the next instruction of the current process
func (p *PCB) GetNextInstruction() (string, error) {
	if p.PC == p.getVariablesAddress() {
		return "", EndOfInstructionsErr
	}

	instruction := p.ram[p.PC]
	return instruction, nil
}

func (p *PCB) IncrementPC() error {
	_, err := p.GetNextInstruction()
	if err != nil {
		return EndOfInstructionsErr
	}

	p.PC++
	return nil
}

// SetDataWord put data in memory in the specified location
func (p *PCB) SetDataWord(virtualLocation int, data int) error {
	if virtualLocation < 0 || virtualLocation > 3 {
		return ProtectionErr
	}

	physicalLocation := virtualLocation + p.getVariablesAddress()
	p.ram[physicalLocation] = fmt.Sprint(data)
	return nil
}

// GetDataWord retrieve data from memory from the specified location
func (p *PCB) GetDataWord(virtualLocation int) (string, error) {
	if virtualLocation < 0 || virtualLocation > 2 {
		return "", ProtectionErr
	}

	physicalLocation := virtualLocation + p.getVariablesAddress()
	return p.ram[physicalLocation], nil
}
