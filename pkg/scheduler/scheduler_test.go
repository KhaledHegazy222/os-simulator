package scheduler

import (
	"reflect"
	"testing"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

var (
	processA = &memory.PCB{
		Id:    1,
		State: memory.Ready,
	}
	processB = &memory.PCB{
		Id:    2,
		State: memory.Blocked,
	}
	processC = &memory.PCB{
		Id:    3,
		State: memory.Ready,
	}
	processD = &memory.PCB{
		Id:    4,
		State: memory.Ready,
	}
	processE = &memory.PCB{
		Id:    5,
		State: memory.Ready,
	}
	processF = &memory.PCB{
		Id:    5,
		State: memory.Blocked,
	}
	processG = &memory.PCB{
		Id:    5,
		State: memory.Blocked,
	}
)

func TestAddToReadyQueue(t *testing.T) {
	s := NewScheduler()

	if len(s.readyQueue) != 0 {
		t.Errorf("expected 0, found %v", len(s.readyQueue))
	}

	if err := s.AddToReadyQueue(processA); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if err := s.AddToReadyQueue(processB); err != ErrProcessNotReady {
		t.Errorf("expected %v, found %v", ErrProcessNotFound, err)
	}

	if err := s.AddToReadyQueue(processC); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if len(s.readyQueue) != 2 {
		t.Errorf("expected 2, found %v", len(s.readyQueue))
	}
}

func TestGetNextReadyProcess(t *testing.T) {
	s := NewScheduler()

	if _, err := s.GetNextReadyProcess(); err != ErrNoReadyProcesses {
		t.Errorf("expected %v, found %v", ErrNoReadyProcesses, err)
	}

	if err := s.AddToReadyQueue(processA); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if err := s.AddToReadyQueue(processC); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if err := s.AddToReadyQueue(processD); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	// get first process
	if process, err := s.GetNextReadyProcess(); err != nil && !reflect.DeepEqual(process, processA) {
		t.Errorf("expected nil, found %v", err)
		t.Errorf("expected %v, found %v", process, processA)
	}

	// get second process
	if process, err := s.GetNextReadyProcess(); err != nil && !reflect.DeepEqual(process, processC) {
		t.Errorf("expected nil, found %v", err)
		t.Errorf("expected %v, found %v", process, processC)
	}

	// get third process
	if process, err := s.GetNextReadyProcess(); err != nil && !reflect.DeepEqual(process, processD) {
		t.Errorf("expected nil, found %v", err)
		t.Errorf("expected %v, found %v", process, processD)
	}
}

func TestBlockProcess(t *testing.T) {
	t.Run("block not existing process", func(t *testing.T) {
		s := NewScheduler()

		s.AddToReadyQueue(processA)
		err := s.BlockProcess(1000)
		if err != ErrProcessNotFound {
			t.Errorf("expected %v, found %v", ErrProcessNotFound, err)
		}
	})
	t.Run("block existing process", func(t *testing.T) {
		s := NewScheduler()
		testProcess := &memory.PCB{
			Id:    1,
			State: memory.Ready,
		}
		// add two ready processes to the ready queue
		if err := s.AddToReadyQueue(testProcess); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if err := s.AddToReadyQueue(processC); err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		// block a process and check moving the process from ready queue to blocked queue.
		if err := s.BlockProcess(testProcess.Id); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if len(s.readyQueue) != 1 && !reflect.DeepEqual(s.readyQueue[0], processC) {
			t.Errorf("expected 1, found %v", len(s.readyQueue))
			t.Errorf("expected %v, found %v", processC, s.blockedQueue[0])

		}
		if len(s.blockedQueue) != 1 && !reflect.DeepEqual(s.blockedQueue[0], testProcess) {
			t.Errorf("expected 1, found %v", len(s.blockedQueue))
			t.Errorf("expected %v, found %v", testProcess, s.blockedQueue[0])
		}
	})
}

func TestUnblockProcess(t *testing.T) {
	t.Run("unblock not existing process", func(t *testing.T) {
		s := NewScheduler()

		s.addToBlockedQueue(processB)
		err := s.UnBlockProcess(1000)
		if err != ErrProcessNotFound {
			t.Errorf("expected %v, found %v", ErrProcessNotFound, err)
		}
	})

	t.Run("unblock existing process", func(t *testing.T) {
		s := NewScheduler()

		testProcess := &memory.PCB{
			Id:    1,
			State: memory.Blocked,
		}
		// add two blocked processes to the blocked queue
		if err := s.addToBlockedQueue(testProcess); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if err := s.addToBlockedQueue(processF); err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		// unblock a process and check moving the process from blocked queue to ready queue.
		if err := s.UnBlockProcess(testProcess.Id); err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		if len(s.readyQueue) != 1 && !reflect.DeepEqual(s.readyQueue[0], testProcess) {
			t.Errorf("expected 1, found %v", len(s.readyQueue))
			t.Errorf("expected %v, found %v", testProcess, s.readyQueue[0])
		}

		if len(s.blockedQueue) != 1 && !reflect.DeepEqual(s.blockedQueue[0], processE) {
			t.Errorf("expected 1, found %v", len(s.blockedQueue))
			t.Errorf("expected %v, found %v", processE, s.blockedQueue[0])
		}
	})
}

func TestTerminateProcess(t *testing.T) {
	t.Run("terminate not existing process", func(t *testing.T) {
		s := NewScheduler()

		s.AddToReadyQueue(processA)
		err := s.TerminateProcess(1000)
		if err != ErrProcessNotFound {
			t.Errorf("expected %v, found %v", ErrProcessNotFound, err)
		}
	})

	t.Run("terminate existing process", func(t *testing.T) {
		s := NewScheduler()

		// add two ready processes to the ready queue
		if err := s.AddToReadyQueue(processA); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if err := s.AddToReadyQueue(processC); err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		// terminate a process and check removing the process from ready queue.
		if err := s.TerminateProcess(processA.Id); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if len(s.readyQueue) != 1 && !reflect.DeepEqual(s.readyQueue[0], processC) {
			t.Errorf("expected 1, found %v", len(s.readyQueue))
			t.Errorf("expected %v, found %v", processC, s.blockedQueue[0])

		}
	})
}

func TestIncrementIterator(t *testing.T) {
	s := NewScheduler()

	if err := s.AddToReadyQueue(processA); err != nil {
		t.Errorf("expected nil, found %v", err)
	}
	if err := s.AddToReadyQueue(processC); err != nil {
		t.Errorf("expected nil, found %v", err)
	}
	if err := s.AddToReadyQueue(processD); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if process, _ := s.GetNextReadyProcess(); !reflect.DeepEqual(process, processA) {
		t.Errorf("expected %v, found %v", processA, process)
	}
	if process, _ := s.GetNextReadyProcess(); !reflect.DeepEqual(process, processC) {
		t.Errorf("expected %v, found %v", processC, process)
	}
	if process, _ := s.GetNextReadyProcess(); !reflect.DeepEqual(process, processD) {
		t.Errorf("expected %v, found %v", processD, process)
	}
	if process, _ := s.GetNextReadyProcess(); !reflect.DeepEqual(process, processA) {
		t.Errorf("expected %v, found %v", processA, process)
	}
}

func TestAddToBlockedQueue(t *testing.T) {

	t.Run("add ready process to blocked queue, should return error", func(t *testing.T) {
		s := NewScheduler()
		err := s.addToBlockedQueue(processA)
		if err != ErrProcessNotBlocked {
			t.Errorf("expected %v, found %v", ErrProcessNotBlocked, err)
		}
	})

	t.Run("add blocked process to blocked queue", func(t *testing.T) {
		s := NewScheduler()

		if err := s.addToBlockedQueue(processB); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if len(s.blockedQueue) != 1 && !reflect.DeepEqual(s.blockedQueue[0], processB) {
			t.Errorf("expected 1, found %v", len(s.blockedQueue))
			t.Errorf("expected %v, found %v", processB, s.blockedQueue[0])
		}
	})
}

func TestRemoveFromBlockedQueue(t *testing.T) {
	s := NewScheduler()

	// add three blocked processes to the blocked queue
	if err := s.addToBlockedQueue(processB); err != nil {
		t.Errorf("expected nil, found %v", err)
	}
	if err := s.addToBlockedQueue(processF); err != nil {
		t.Errorf("expected nil, found %v", err)
	}
	if err := s.addToBlockedQueue(processG); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if found := s.removeFromBlockedQueue(1); !reflect.DeepEqual(found, processF) {
		t.Errorf("expected %v, found %v", processF, found)
	}
	if len(s.blockedQueue) != 2 {
		t.Errorf("expected 2, found %v", len(s.blockedQueue))
	}
}

func TestRemoveFromReadyQueue(t *testing.T) {

	t.Run("iterator is less than deleted process index", func(t *testing.T) {
		s := NewScheduler()

		// add four ready processes to the ready queue
		s.AddToReadyQueue(processA)
		s.AddToReadyQueue(processC)
		s.AddToReadyQueue(processD)
		s.AddToReadyQueue(processE)

		// iterator is zero and index is two, so should not change
		if process := s.removeFromReadyQueue(2); !reflect.DeepEqual(process, processD) {
			t.Errorf("expected %v, found %v", processD, process)
		}

		if expected := []*memory.PCB{processA, processC, processE}; !reflect.DeepEqual([]*memory.PCB(s.readyQueue), expected) {
			t.Errorf("expected %v, found %v", expected, s.readyQueue)
		}

		if s.ReadyProcessIterator != 0 {
			t.Errorf("expected 1, found %v", s.ReadyProcessIterator)
		}
	})

	t.Run("iterator is equal to index,not last index", func(t *testing.T) {
		s := NewScheduler()

		// add four ready processes to the ready queue
		s.AddToReadyQueue(processA)
		s.AddToReadyQueue(processC)
		s.AddToReadyQueue(processD)
		s.AddToReadyQueue(processE)

		// iterator is one and index is one, so should not change
		if process, _ := s.GetNextReadyProcess(); process != processA {
			t.Errorf("expected %v, found %v", processA, process)
		}
		if process := s.removeFromReadyQueue(1); !reflect.DeepEqual(process, processC) {
			t.Errorf("expected %v, found %v", processC, process)
		}

		if expected := []*memory.PCB{processA, processD, processE}; !reflect.DeepEqual([]*memory.PCB(s.readyQueue), expected) {
			t.Errorf("expected %v, found %v", expected, s.readyQueue)
		}

		if s.ReadyProcessIterator != 1 {
			t.Errorf("expected 1, found %v", s.ReadyProcessIterator)
		}
	})

	t.Run("iterator is equal to index, last index", func(t *testing.T) {
		s := NewScheduler()

		// add four ready processes to the ready queue
		s.AddToReadyQueue(processA)
		s.AddToReadyQueue(processC)
		s.AddToReadyQueue(processD)
		s.AddToReadyQueue(processE)

		// iterator is three and index is three, so should not change
		if process, _ := s.GetNextReadyProcess(); process != processA {
			t.Errorf("expected %v, found %v", processA, process)
		}
		if process, _ := s.GetNextReadyProcess(); process != processC {
			t.Errorf("expected %v, found %v", processC, process)
		}
		if process, _ := s.GetNextReadyProcess(); process != processD {
			t.Errorf("expected %v, found %v", processD, process)
		}

		if s.ReadyProcessIterator != 3 {
			t.Errorf("expected 3, found %v", s.ReadyProcessIterator)
		}

		if process := s.removeFromReadyQueue(3); !reflect.DeepEqual(process, processE) {
			t.Errorf("expected %v, found %v", processE, process)
		}

		if expected := []*memory.PCB{processA, processC, processD}; !reflect.DeepEqual([]*memory.PCB(s.readyQueue), expected) {
			t.Errorf("expected %v, found %v", expected, s.readyQueue)
		}

		if s.ReadyProcessIterator != 0 {
			t.Errorf("expected 0, found %v", s.ReadyProcessIterator)
		}
	})

	t.Run("iterator is greater than deleted process index", func(t *testing.T) {
		s := NewScheduler()

		// add four ready processes to the ready queue
		s.AddToReadyQueue(processA)
		s.AddToReadyQueue(processC)
		s.AddToReadyQueue(processD)
		s.AddToReadyQueue(processE)

		// iterator is three and index is one, so should subtract one from the iterator
		if process, _ := s.GetNextReadyProcess(); process != processA {
			t.Errorf("expected %v, found %v", processA, process)
		}
		if process, _ := s.GetNextReadyProcess(); process != processC {
			t.Errorf("expected %v, found %v", processC, process)
		}
		if process, _ := s.GetNextReadyProcess(); process != processD {
			t.Errorf("expected %v, found %v", processD, process)
		}

		if s.ReadyProcessIterator != 3 {
			t.Errorf("expected 3, found %v", s.ReadyProcessIterator)
		}

		if process := s.removeFromReadyQueue(1); !reflect.DeepEqual(process, processC) {
			t.Errorf("expected %v, found %v", processC, process)
		}

		if expected := []*memory.PCB{processA, processD, processE}; !reflect.DeepEqual([]*memory.PCB(s.readyQueue), expected) {
			t.Errorf("expected %v, found %v", expected, s.readyQueue)
		}

		if s.ReadyProcessIterator != 2 {
			t.Errorf("expected 2, found %v", s.ReadyProcessIterator)
		}
	})
}
