package memory

import (
	"errors"
	"fmt"
	"strconv"
)

var UnableToRetrievePCBErr = errors.New("not able to retrieve PCB from memory")



// RAMMemory represents a Fixed-sized 40 words RAM memory
type RAMMemory [memorySize]string




func (ram *RAMMemory) isFree(from int, to int) bool {
	if from > memoryEndAddress || from < memoryStartAddress || to > memoryEndAddress || to < memoryStartAddress {
		return false
	}

	for i := from; i <= to; i++ {
		if ram[i] != "" {
			return false
		}
	}
	return true
}

func (ram *RAMMemory) allocateProcess(start int, unparsedCode []string,id int) PCB {
	end := start + getProcessSize(len(unparsedCode)) -1

	pcb := PCB{
		Id:       id,
		State:    Ready,
		PC:       start + PCBSize,
		Start:    start,
		End:      end,
		CodeSize: len(unparsedCode),
		ram:ram,
	}

	// allocate pcb in the first 6 words
	pcbAddress:=pcb.getPCBAddress()
	ram[pcbAddress] = fmt.Sprint(pcb.Id)
	ram[pcbAddress+1] = fmt.Sprint(pcb.State)
	ram[pcbAddress+2] = fmt.Sprint(pcb.PC)
	ram[pcbAddress+3] = fmt.Sprint(pcb.Start)
	ram[pcbAddress+4] = fmt.Sprint(pcb.End)
	ram[pcbAddress+5] = fmt.Sprint(pcb.CodeSize)

	// allocate unparsed code
	unparsedCodeStartAddress := pcb.getUnparsedCodeAddress()
	for i := 0; i < len(unparsedCode); i++ {
		ram[unparsedCodeStartAddress+i] = unparsedCode[i]
	}

	// allocate variables in the last three words. initial value is zero
	variablesStartAddress := pcb.getVariablesAddress()
	ram[variablesStartAddress] = fmt.Sprint(0)
	ram[variablesStartAddress+1] = fmt.Sprint(0)
	ram[variablesStartAddress+2] = fmt.Sprint(0)

	return pcb
}

func (ram *RAMMemory) getProcessPCB(startLocation int) (PCB, error) {

	id, idErr := strconv.Atoi(ram[startLocation])
	state:=STATE(ram[startLocation+1])
	pc, pcErr := strconv.Atoi(ram[startLocation+2])
	start, startErr := strconv.Atoi(ram[startLocation+3])
	end, endErr := strconv.Atoi(ram[startLocation+4])
	codeSize, codeSizeErr := strconv.Atoi(ram[startLocation+5])

	if idErr != nil || pcErr != nil || startErr != nil || endErr != nil || codeSizeErr != nil {
		return PCB{}, UnableToRetrievePCBErr
	}

	pcb := PCB{
		Id:       id,
		State:    state,
		PC:       pc,
		Start:    start,
		End:      end,
		CodeSize: codeSize,
		ram: ram,
	}
	return pcb, nil
}
