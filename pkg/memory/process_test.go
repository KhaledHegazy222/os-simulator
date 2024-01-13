package memory

import (
	"testing"
)

func TestGetUnparsedCodeAddress(t *testing.T) {
	process := PCB{
		Start: 10,
	}

	got := process.getUnparsedCodeAddress()

	if got != 10+PCBSize {
		t.Errorf("expected %v found %v", 10+PCBSize, got)
	}
}

func TestGetVariablesAddress(t *testing.T) {
	process := PCB{
		Start:    10,
		CodeSize: 6,
	}

	got := process.getVariablesAddress()

	if got != 10+PCBSize+6 {
		t.Errorf("expected %v found %v", 10+PCBSize, got)
	}
}

func TestGetNextInstruction(t *testing.T) {

	t.Run("normal case return next instruction", func(t *testing.T) {
		var ram RAMMemory
		process := PCB{
			Start:    10,
			CodeSize: 6,
			PC:       16,
			ram:      &ram,
		}
		ram[process.PC] = "word"

		instruction, err := process.GetNextInstruction()

		if err != nil {
			t.Errorf("expected nil found %v", err)
		}

		if instruction != "word" {
			t.Errorf("expected word found %v", instruction)
		}
	})

	t.Run("return error when reaching end of instructions", func(t *testing.T) {
		var ram RAMMemory
		process := PCB{
			Start:    10,
			CodeSize: 6,
			PC:       22,
			ram:      &ram,
		}

		instruction, err := process.GetNextInstruction()

		if err != EndOfInstructionsErr {
			t.Errorf("expected nil found %v", err)
		}
		if instruction != "" {
			t.Errorf("expected empty string found %v", instruction)
		}
	})
}

func TestSetDataWord(t *testing.T) {
	t.Run("normal case return data at virtual location 2", func(t *testing.T) {
		var ram RAMMemory
		process := PCB{
			Start:    10,
			CodeSize: 6,
			PC:       16,
			ram:      &ram,
		}
		
		location := process.getVariablesAddress()
		ram[location]="13"
		ram[location+1]="14"
		
		data,err:=process.GetDataWord(0)
		
		if err != nil {
			t.Errorf("expected nil found %v", err)
		}
		
		if data != "13" {
			t.Errorf("expected 13 found %v", data)
		}

		data,err=process.GetDataWord(1)
		if err != nil {
			t.Errorf("expected nil found %v", err)
		}
		
		if data != "14" {
			t.Errorf("expected 14 found %v", data)
		}
	})

	t.Run("return error when virtual address is not between 0 and 2", func(t *testing.T) {
		var ram RAMMemory
		process := PCB{
			Start:    10,
			CodeSize: 6,
			PC:       16,
			ram:      &ram,
		}
		
		data,err:=process.GetDataWord(-1)
		
		if err != ProtectionErr {
			t.Errorf("expected %v found %v", ProtectionErr,err)
		}
		
		if data != "" {
			t.Errorf("expected empty string found %v", data)
		}

		data,err=process.GetDataWord(3)
		if err != ProtectionErr {
			t.Errorf("expected %v found %v", ProtectionErr,err)
		}
		
		if data != "" {
			t.Errorf("expected empty string found %v", data)
		}
	})
}


func TestIncrementPC(t *testing.T) {
	t.Run("normal case increment pc", func(t *testing.T) {
		var ram RAMMemory
		process := PCB{
			Start:    10,
			CodeSize: 6,
			PC:       16,
			ram:      &ram,
		}

		err:=process.IncrementPC()

		if err != nil {
			t.Errorf("expected nil found %v", err)
		}
		
		if process.PC!=17{
			t.Errorf("expected 17 found %v", process.PC)
		}
	})

	t.Run("reached end of instructions returns error", func(t *testing.T) {
		var ram RAMMemory
		process := PCB{
			Start:    10,
			CodeSize: 6,
			PC:       22,
			ram:      &ram,
		}

		err:=process.IncrementPC()

		if err != EndOfInstructionsErr {
			t.Errorf("expected %v found %v",EndOfInstructionsErr, err)
		}
		
		if process.PC!=22{
			t.Errorf("expected 22 found %v", process.PC)
		}
	})
}