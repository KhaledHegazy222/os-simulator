package memory

import (
	"reflect"
	"testing"
)

func TestAllocateProcess(t *testing.T) {
	var ram RAMMemory

	unparsedCode := []string{
		"assign x 4",
		"assign x 4",
		"assign x 4",
		"assign x 4",
	}

	found := ram.allocateProcess(10, unparsedCode,1)

	expected := PCB{
		Start:    10,
		PC:       16,
		End:      23,
		CodeSize: 4,
	}

	if found.Start != expected.Start {
		t.Errorf("expected %v but found %v", expected.Start, found.Start)
	}

	if found.PC != expected.PC {
		t.Errorf("expected %v but found %v", expected.PC, found.PC)
	}

	if found.CodeSize != expected.CodeSize {
		t.Errorf("expected %v but found %v", expected.CodeSize, found.CodeSize)
	}

	expectedPCBInMemory := []string{
		string(Ready),
		"16",
		"10",
		"22",
		"4",
	}

	foundPCBInMemory := ram[11:16]

	if !reflect.DeepEqual(expectedPCBInMemory, foundPCBInMemory) {
		t.Errorf("expected %v but found %v", expectedPCBInMemory, foundPCBInMemory)
	}
}

func TestGetProcessPCB(t *testing.T) {
	t.Run("normal case successful retrieval of process pcb", func(t *testing.T) {
		var ram RAMMemory
		unparsedCode := []string{
			"assign x 4",
			"print x",
			"semWait file",
		}

		pcb := ram.allocateProcess(10, unparsedCode,1)

		found, err := ram.getProcessPCB(10)

		if err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		if !reflect.DeepEqual(pcb, found) {
			t.Errorf("expected %v, found %v", pcb, found)
		}
	})

	t.Run("failed retrieval of process pcb", func(t *testing.T) {
		var ram RAMMemory
		unparsedCode := []string{
			"assign x 4",
			"print x",
			"semWait file",
		}

		ram.allocateProcess(10, unparsedCode,2)

		_, err := ram.getProcessPCB(12)

		if err != UnableToRetrievePCBErr {
			t.Errorf("expected nil, found %v", err)
		}
	})
}
