package memory

import (
	"testing"
)

var unparsedCode = []string{
	"assign x 4",
	"assign x 4",
	"assign x 4",
	"assign x 4",
}

func TestGetProcessSize(t *testing.T) {
	found := getProcessSize(14)
	expected := PCBSize + 14 + variablesSize

	if found != expected {
		t.Errorf("found %v, expected %v", found, expected)
	}
}

func TestAddProcess(t *testing.T) {
	var ram RAMMemory
	memoryManager := MemoryManager{
		ram:             ram,
		processLocation: make(map[int]int),
	}

	// add first process from 1 to 13
	pcb, err := memoryManager.AddProcess(unparsedCode)
	if pcb.Start != 1 {
		t.Errorf("expected 1, but found %v", pcb.Start)
	}
	expectedEnd := pcb.Start + getProcessSize(4) - 1
	if pcb.End != expectedEnd {
		t.Errorf("expected %v, but found %v", expectedEnd, pcb.End)
	}

	if err != nil {
		t.Errorf("expected nil, but found %v", err)
	}

	// add second process from  14 to 26
	pcb, err = memoryManager.AddProcess(unparsedCode)
	if pcb.Start != 14 {
		t.Errorf("expected 1, but found %v", pcb.Start)
	}

	expectedEnd = pcb.Start + getProcessSize(4) - 1
	if pcb.End != expectedEnd {
		t.Errorf("expected %v, but found %v", expectedEnd, pcb.End)
	}

	if err != nil {
		t.Errorf("expected nil, but found %v", err)
	}

	// add third process from  27 to 38
	pcb, err = memoryManager.AddProcess(unparsedCode)
	if pcb.Start != 27 {
		t.Errorf("expected 1, but found %v", pcb.Start)
	}

	expectedEnd = pcb.Start + getProcessSize(4) - 1
	if pcb.End != expectedEnd {
		t.Errorf("expected %v, but found %v", expectedEnd, pcb.End)
	}

	if err != nil {
		t.Errorf("expected nil, but found %v", err)
	}

	// fail to add any other process
	_, err = memoryManager.AddProcess(unparsedCode)

	if err != NotEnoughSpaceErr {
		t.Errorf("expected %v, but found %v", NotEnoughSpaceErr, err)
	}
}

func TestDeleteProcess(t *testing.T) {
	memoryManager := NewMemoryManager()

	pcb, err := memoryManager.AddProcess(unparsedCode)

	if err != nil {
		t.Errorf("expected nil but found %v", err)
	}

	for i := pcb.Start; i <= pcb.End; i++ {
		if memoryManager.ram[i] == "" {
			t.Errorf("at %v: expected process info but found empty line", i)
		}
	}

	if memoryManager.processLocation[pcb.Id] != pcb.Start {
		t.Errorf("memory map of process %v: expected %v, but found %v ", pcb.Id, pcb.Id, memoryManager.processLocation[pcb.Id])
	}

	err = memoryManager.DeleteProcess(pcb.Id)
	if err != nil {
		t.Errorf("expected nil but found %v", err)
	}

	if memoryManager.processLocation[pcb.Id] != 0 {
		t.Errorf("expected 0 but found %v", memoryManager.processLocation[pcb.Id])
	}

	for i := pcb.Start; i <= pcb.End; i++ {
		if memoryManager.ram[i] != "" {
			t.Errorf("at %v: expected empty line but found %v", i, memoryManager.ram[i])
		}
	}
}
