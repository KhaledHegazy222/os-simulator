package scheduler

import "github.com/KhaledHegazy222/os-simulator/pkg/memory"

type queue []*memory.PCB

func (q *queue) append(pcb *memory.PCB) {
	*q = append(*q, pcb)
}

func (q *queue) delete(index int) *memory.PCB {
	pcb := (*q)[index]
	*q = append((*q)[:index], (*q)[index+1:]...)
	return pcb
}
