package scheduler

import (
	"reflect"
	"testing"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

func TestAppend(t *testing.T) {
	var q queue
	firstProcess := &memory.PCB{
		Id:    1,
		State: memory.Ready,
	}
	secondProcess := &memory.PCB{
		Id:    2,
		State: memory.Ready,
	}
	thirdProcess := &memory.PCB{
		Id:    3,
		State: memory.Ready,
	}
	q.append(firstProcess)
	q.append(secondProcess)
	q.append(thirdProcess)

	if len(q) != 3 {
		t.Errorf("expected 3, found %v", len(q))
	}
}

func TestDelete(t *testing.T) {
	var q queue
	firstProcess := &memory.PCB{
		Id:    1,
		State: memory.Ready,
	}
	secondProcess := &memory.PCB{
		Id:    2,
		State: memory.Ready,
	}
	thirdProcess := &memory.PCB{
		Id:    3,
		State: memory.Ready,
	}
	q.append(firstProcess)
	q.append(secondProcess)
	q.append(thirdProcess)
	q.delete(1)

	expected := []*memory.PCB{firstProcess, thirdProcess}

	for i := 0; i < 2; i++ {
		if !reflect.DeepEqual(q[i], expected[i]) {
			t.Errorf("expected %v, found %v", expected[i], q[i])
			break
		}
	}
}
