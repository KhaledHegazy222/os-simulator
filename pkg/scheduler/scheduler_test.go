package scheduler

import (
	"reflect"
	"testing"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

func TestAddToReadyQueue(t *testing.T) {
	s := NewScheduler()
	firstReadyProcess := &memory.PCB{
		Id:    1,
		State: memory.Ready,
	}
	blockedProcess := &memory.PCB{
		Id:    2,
		State: memory.Blocked,
	}
	secondReadyProcess := &memory.PCB{
		Id:    3,
		State: memory.Ready,
	}

	if len(s.readyQueue) != 0 {
		t.Errorf("expected 0, found %v", len(s.readyQueue))
	}

	if err := s.AddToReadyQueue(firstReadyProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if err := s.AddToReadyQueue(blockedProcess); err != ErrProcessNotReady {
		t.Errorf("expected %v, found %v", ErrProcessNotFound, err)
	}

	if err := s.AddToReadyQueue(secondReadyProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if len(s.readyQueue) != 2 {
		t.Errorf("expected 2, found %v", len(s.readyQueue))
	}
}

func TestGetNextReadyProcess(t *testing.T) {
	s := NewScheduler()
	firstReadyProcess := &memory.PCB{
		Id:    1,
		State: memory.Ready,
	}
	secondReadyProcess := &memory.PCB{
		Id:    2,
		State: memory.Ready,
	}
	thirdReadyProcess := &memory.PCB{
		Id:    3,
		State: memory.Ready,
	}

	if _, err := s.GetNextReadyProcess(); err != ErrNoReadyProcesses {
		t.Errorf("expected %v, found %v", ErrNoReadyProcesses, err)
	}

	if err := s.AddToReadyQueue(firstReadyProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if err := s.AddToReadyQueue(secondReadyProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if err := s.AddToReadyQueue(thirdReadyProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	// get first process
	if process, err := s.GetNextReadyProcess(); err != nil && !reflect.DeepEqual(process, firstReadyProcess) {
		t.Errorf("expected nil, found %v", err)
		t.Errorf("expected %v, found %v", process, firstReadyProcess)
	}

	// get second process
	if process, err := s.GetNextReadyProcess(); err != nil && !reflect.DeepEqual(process, secondReadyProcess) {
		t.Errorf("expected nil, found %v", err)
		t.Errorf("expected %v, found %v", process, secondReadyProcess)
	}

	// get third process
	if process, err := s.GetNextReadyProcess(); err != nil && !reflect.DeepEqual(process, thirdReadyProcess) {
		t.Errorf("expected nil, found %v", err)
		t.Errorf("expected %v, found %v", process, thirdReadyProcess)
	}
}

func TestBlockProcess(t *testing.T) {
	t.Run("block not existing process", func(t *testing.T) {
		s := NewScheduler()
		readyProcess := &memory.PCB{
			Id:    1,
			State: memory.Ready,
		}
		s.AddToReadyQueue(readyProcess)
		err := s.BlockProcess(1000)
		if err != ErrProcessNotFound {
			t.Errorf("expected %v, found %v", ErrProcessNotFound, err)
		}
	})
	t.Run("block existing process", func(t *testing.T) {
		s := NewScheduler()
		firstReadyProcess := &memory.PCB{
			Id:    1,
			State: memory.Ready,
		}
		secondReadyProcess := &memory.PCB{
			Id:    2,
			State: memory.Ready,
		}
		// add two ready processes to the ready queue
		if err := s.AddToReadyQueue(firstReadyProcess); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if err := s.AddToReadyQueue(secondReadyProcess); err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		// block a process and check moving the process from ready queue to blocked queue.
		if err := s.BlockProcess(firstReadyProcess.Id); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if len(s.readyQueue) != 1 && !reflect.DeepEqual(s.readyQueue[0], secondReadyProcess) {
			t.Errorf("expected 1, found %v", len(s.readyQueue))
			t.Errorf("expected %v, found %v", secondReadyProcess, s.blockedQueue[0])

		}
		if len(s.blockedQueue) != 1 && !reflect.DeepEqual(s.blockedQueue[0], firstReadyProcess) {
			t.Errorf("expected 1, found %v", len(s.blockedQueue))
			t.Errorf("expected %v, found %v", firstReadyProcess, s.blockedQueue[0])
		}
	})
}

func TestUnblockProcess(t *testing.T) {
	t.Run("unblock not existing process", func(t *testing.T) {
		s := NewScheduler()
		blockedProcess := &memory.PCB{
			Id:    1,
			State: memory.Blocked,
		}
		notFoundId := 1000
		s.addToBlockedQueue(blockedProcess)
		err := s.UnBlockProcess(notFoundId)
		if err != ErrProcessNotFound {
			t.Errorf("expected %v, found %v", ErrProcessNotFound, err)
		}
	})

	t.Run("unblock existing process", func(t *testing.T) {
		s := NewScheduler()

		firstBlockedProcess := &memory.PCB{
			Id:    1,
			State: memory.Blocked,
		}
		secondBlockedProcess := &memory.PCB{
			Id:    2,
			State: memory.Blocked,
		}
		// add two blocked processes to the blocked queue
		if err := s.addToBlockedQueue(firstBlockedProcess); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if err := s.addToBlockedQueue(secondBlockedProcess); err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		// unblock a process and check moving the process from blocked queue to ready queue.
		if err := s.UnBlockProcess(firstBlockedProcess.Id); err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		if len(s.readyQueue) != 1 && !reflect.DeepEqual(s.readyQueue[0], firstBlockedProcess) {
			t.Errorf("expected 1, found %v", len(s.readyQueue))
			t.Errorf("expected %v, found %v", firstBlockedProcess, s.readyQueue[0])
		}

		if len(s.blockedQueue) != 1 && !reflect.DeepEqual(s.blockedQueue[0], secondBlockedProcess) {
			t.Errorf("expected 1, found %v", len(s.blockedQueue))
			t.Errorf("expected %v, found %v", secondBlockedProcess, s.blockedQueue[0])
		}
	})
}

func TestTerminateProcess(t *testing.T) {
	t.Run("terminate not existing process", func(t *testing.T) {
		s := NewScheduler()

		readyProcess := &memory.PCB{
			Id:    1,
			State: memory.Ready,
		}
		notFoundId := 1000

		s.AddToReadyQueue(readyProcess)
		err := s.TerminateProcess(notFoundId)
		if err != ErrProcessNotFound {
			t.Errorf("expected %v, found %v", ErrProcessNotFound, err)
		}
	})

	t.Run("terminate existing process", func(t *testing.T) {
		s := NewScheduler()

		firstReadyProcess := &memory.PCB{
			Id:    1,
			State: memory.Ready,
		}
		secondReadyProcess := &memory.PCB{
			Id:    2,
			State: memory.Ready,
		}

		// add two ready processes to the ready queue
		if err := s.AddToReadyQueue(firstReadyProcess); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if err := s.AddToReadyQueue(secondReadyProcess); err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		// terminate a process and check removing the process from ready queue.
		if err := s.TerminateProcess(firstReadyProcess.Id); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if len(s.readyQueue) != 1 && !reflect.DeepEqual(s.readyQueue[0], secondReadyProcess) {
			t.Errorf("expected 1, found %v", len(s.readyQueue))
			t.Errorf("expected %v, found %v", secondReadyProcess, s.blockedQueue[0])

		}
	})
}

func TestIncrementIterator(t *testing.T) {
	s := NewScheduler()

	firstReadyProcess := &memory.PCB{
		Id:    1,
		State: memory.Ready,
	}
	secondReadyProcess := &memory.PCB{
		Id:    2,
		State: memory.Ready,
	}
	thirdReadyProcess := &memory.PCB{
		Id:    3,
		State: memory.Ready,
	}

	if err := s.AddToReadyQueue(firstReadyProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}
	if err := s.AddToReadyQueue(secondReadyProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}
	if err := s.AddToReadyQueue(thirdReadyProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if process, _ := s.GetNextReadyProcess(); !reflect.DeepEqual(process, firstReadyProcess) {
		t.Errorf("expected %v, found %v", firstReadyProcess, process)
	}
	if process, _ := s.GetNextReadyProcess(); !reflect.DeepEqual(process, secondReadyProcess) {
		t.Errorf("expected %v, found %v", secondReadyProcess, process)
	}
	if process, _ := s.GetNextReadyProcess(); !reflect.DeepEqual(process, thirdReadyProcess) {
		t.Errorf("expected %v, found %v", thirdReadyProcess, process)
	}
	if process, _ := s.GetNextReadyProcess(); !reflect.DeepEqual(process, firstReadyProcess) {
		t.Errorf("expected %v, found %v", firstReadyProcess, process)
	}
}

func TestAddToBlockedQueue(t *testing.T) {

	t.Run("add ready process to blocked queue, should return error", func(t *testing.T) {
		readyProcess := &memory.PCB{
			Id:    2,
			State: memory.Ready,
		}
		s := NewScheduler()
		err := s.addToBlockedQueue(readyProcess)
		if err != ErrProcessNotBlocked {
			t.Errorf("expected %v, found %v", ErrProcessNotBlocked, err)
		}
	})

	t.Run("add blocked process to blocked queue", func(t *testing.T) {
		s := NewScheduler()
		blockedProcess := &memory.PCB{
			Id:    2,
			State: memory.Blocked,
		}
		if err := s.addToBlockedQueue(blockedProcess); err != nil {
			t.Errorf("expected nil, found %v", err)
		}
		if len(s.blockedQueue) != 1 && !reflect.DeepEqual(s.blockedQueue[0], blockedProcess) {
			t.Errorf("expected 1, found %v", len(s.blockedQueue))
			t.Errorf("expected %v, found %v", blockedProcess, s.blockedQueue[0])
		}
	})
}

func TestRemoveFromBlockedQueue(t *testing.T) {
	s := NewScheduler()

	firstBlockedProcess := &memory.PCB{
		Id:    1,
		State: memory.Blocked,
	}
	secondBlockedProcess := &memory.PCB{
		Id:    2,
		State: memory.Blocked,
	}
	thirdBlockedProcess := &memory.PCB{
		Id:    3,
		State: memory.Blocked,
	}

	// add three blocked processes to the blocked queue
	if err := s.addToBlockedQueue(firstBlockedProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}
	if err := s.addToBlockedQueue(secondBlockedProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}
	if err := s.addToBlockedQueue(thirdBlockedProcess); err != nil {
		t.Errorf("expected nil, found %v", err)
	}

	if found := s.removeFromBlockedQueue(1); !reflect.DeepEqual(found, secondBlockedProcess) {
		t.Errorf("expected %v, found %v", secondBlockedProcess, found)
	}
	if len(s.blockedQueue) != 2 {
		t.Errorf("expected 2, found %v", len(s.blockedQueue))
	}
}

func TestRemoveFromReadyQueue(t *testing.T) {
	s := NewScheduler()
	firstReadyProcess := &memory.PCB{
		Id:    1,
		State: memory.Ready,
	}
	secondReadyProcess := &memory.PCB{
		Id:    2,
		State: memory.Ready,
	}
	thirdReadyProcess := &memory.PCB{
		Id:    3,
		State: memory.Ready,
	}

	s.AddToReadyQueue(firstReadyProcess)
	s.AddToReadyQueue(secondReadyProcess)
	s.AddToReadyQueue(thirdReadyProcess)

	s.removeFromReadyQueue(1)

	if !reflect.DeepEqual(firstReadyProcess, s.readyQueue[0]) {
		t.Errorf("expected %v, found %v", firstReadyProcess, s.readyQueue[0])
	}
	if !reflect.DeepEqual(thirdReadyProcess, s.readyQueue[1]) {
		t.Errorf("expected %v, found %v", thirdReadyProcess, s.readyQueue[1])
	}

}

func TestNormalizeIterator(t *testing.T) {

	t.Run("iterator is less than deleted process index", func(t *testing.T) {
		s := NewScheduler()

		readyProcess := &memory.PCB{
			Id:    1,
			State: memory.Ready,
		}

		// add four ready processes to the ready queue
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)

		// iterator is zero and index deleted is two, so should not change
		s.readyProcessIterator = 0
		s.normalizeIterator(2)

		if s.readyProcessIterator != 0 {
			t.Errorf("expected 0, found %v", s.readyProcessIterator)
		}
	})

	t.Run("iterator is equal to index,not last index", func(t *testing.T) {
		s := NewScheduler()

		readyProcess := &memory.PCB{
			Id:    1,
			State: memory.Ready,
		}

		// add four ready processes to the ready queue
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)

		// iterator is one and index is one and it's not last index, so should not change
		s.readyProcessIterator = 1
		s.normalizeIterator(1)
		if s.readyProcessIterator != 1 {
			t.Errorf("expected 1, found %v", s.readyProcessIterator)
		}
	})

	t.Run("iterator is equal to index, last index", func(t *testing.T) {
		s := NewScheduler()

		readyProcess := &memory.PCB{
			Id:    1,
			State: memory.Ready,
		}

		// add four ready processes to the ready queue
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)

		// iterator is three and index is three, so should become 0
		s.readyProcessIterator = 3
		s.normalizeIterator(3)
		if s.readyProcessIterator != 0 {
			t.Errorf("expected 0, found %v", s.readyProcessIterator)
		}
	})

	t.Run("iterator is greater than deleted process index", func(t *testing.T) {
		s := NewScheduler()

		readyProcess := &memory.PCB{
			Id:    1,
			State: memory.Ready,
		}

		// add four ready processes to the ready queue
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)
		s.AddToReadyQueue(readyProcess)

		// iterator is three and index is one, so should subtract one from the iterator
		s.readyProcessIterator = 3
		s.normalizeIterator(1)
		if s.readyProcessIterator != 2 {
			t.Errorf("expected 2, found %v", s.readyProcessIterator)
		}

	})
}
