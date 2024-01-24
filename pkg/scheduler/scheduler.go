package scheduler

import (
	"errors"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

type Scheduler struct {
	readyQueue           []*memory.PCB
	blockedQueue         []*memory.PCB
	ReadyProcessIterator int
}

var (
	ErrProcessNotReady   = errors.New("process state is not ready.")
	ErrProcessNotBlocked = errors.New("process state is not blocked.")
	ErrNoReadyProcesses  = errors.New("no processes in the ready queue.")
	ErrProcessNotFound   = errors.New("process is not found in the queue.")
)

// NewScheduler factory function that creates new scheduler
func NewScheduler() *Scheduler {
	readyQueue := make([]*memory.PCB, 0)
	blockedQueue := make([]*memory.PCB, 0)
	return &Scheduler{
		readyQueue:           readyQueue,
		blockedQueue:         blockedQueue,
		ReadyProcessIterator: 0,
	}
}

// AddToReadyQueue adds the given pcb to the ready queue.
func (s *Scheduler) AddToReadyQueue(process *memory.PCB) error {
	if process.State != memory.Ready {
		return ErrProcessNotReady
	}
	s.readyQueue = append(s.readyQueue, process)
	return nil
}

// GetNextReadyProcess gets next ready process from the ready queue.
func (s *Scheduler) GetNextReadyProcess() (*memory.PCB, error) {
	if len(s.readyQueue) == 0 {
		return &memory.PCB{}, ErrNoReadyProcesses
	}

	readyProcess := s.readyQueue[s.ReadyProcessIterator]
	s.incrementIterator()
	return readyProcess, nil
}

// BlockProcess move pcb with given pid from ready queue to block queue.
func (s *Scheduler) BlockProcess(pid int) error {
	for idx, process := range s.readyQueue {
		if process.Id == pid {
			pcb := s.removeFromReadyQueue(idx)
			pcb.State = memory.Blocked
			s.addToBlockedQueue(pcb)
			return nil
		}
	}
	return ErrProcessNotFound
}

// UnBlockProcess move pcb with given pid from block queue to ready queue.
func (s *Scheduler) UnBlockProcess(pid int) error {
	for idx, process := range s.blockedQueue {
		if process.Id == pid {
			pcb := s.removeFromBlockedQueue(idx)
			pcb.State = memory.Ready
			s.AddToReadyQueue(pcb)
			return nil
		}
	}
	return ErrProcessNotFound
}

// TerminateProcess remove process with given pid from the ready queue.
func (s *Scheduler) TerminateProcess(pid int) error {
	for idx, process := range s.readyQueue {
		if process.Id == pid {
			s.removeFromReadyQueue(idx)
			return nil
		}
	}
	return ErrProcessNotFound
}

func (s *Scheduler) incrementIterator() {
	s.ReadyProcessIterator++
	if s.ReadyProcessIterator == len(s.readyQueue) {
		s.ReadyProcessIterator = 0
	}
}

func (s *Scheduler) addToBlockedQueue(process *memory.PCB) error {
	if process.State != memory.Blocked {
		return ErrProcessNotBlocked
	}
	s.blockedQueue = append(s.blockedQueue, process)
	return nil
}

func (s *Scheduler) removeFromBlockedQueue(index int) *memory.PCB {
	pcb := s.blockedQueue[index]
	s.blockedQueue = append(s.blockedQueue[:index], s.blockedQueue[index+1:]...)
	return pcb
}

func (s *Scheduler) removeFromReadyQueue(index int) *memory.PCB {
	if s.ReadyProcessIterator > index {
		s.ReadyProcessIterator--
	}
	if s.ReadyProcessIterator == index && index==len(s.readyQueue)-1{
		s.ReadyProcessIterator=0
	}

	pcb := s.readyQueue[index]
	s.readyQueue = append(s.readyQueue[:index], s.readyQueue[index+1:]...)
	return pcb
}
