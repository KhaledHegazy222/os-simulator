package memory

import "fmt"

type PCB struct {
	Id       int
	State    STATE
	PC       int
	Start    int
	End      int
	codeSize int
}

func (p *PCB) getPCBAddress() int {
	return p.Start
}

func (p *PCB) getUnparsedCodeAddress() int {
	return p.getPCBAddress() + PCBSize
}

func (p *PCB) getVariablesAddress() int {
	return p.getUnparsedCodeAddress() + p.codeSize
}

func (p *PCB) GetNextInstruction() (string, error) {
	if p.PC == p.getVariablesAddress() {
		return "", EndOfInstructionsErr
	}
	instruction := RAMMemory[p.PC]
	p.PC++
	return instruction, nil
}

func (p *PCB) SetDataWord(virtualLocation int, data int) error {
	if virtualLocation < 0 || virtualLocation > 3 {
		return ProtectionErr
	}
	physicalLocation := virtualLocation + p.getVariablesAddress()
	RAMMemory[physicalLocation] = fmt.Sprint(data)
	return nil
}

func (p *PCB) GetDataWord(virtualLocation int) (string, error) {
	if virtualLocation < 0 || virtualLocation > 3 {
		return "", ProtectionErr
	}
	physicalLocation := virtualLocation + p.getVariablesAddress()
	return RAMMemory[physicalLocation], nil
}
