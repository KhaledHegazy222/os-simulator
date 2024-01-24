package scheduler

import (
	"reflect"
	"testing"

	"github.com/KhaledHegazy222/os-simulator/pkg/memory"
)

func TestAppend(t *testing.T) {
	var q queue
	q.append(processA)
	q.append(processB)
	q.append(processC)

	if len(q) != 3 {
		t.Errorf("expected 3, found %v", len(q))
	}
}

func TestDelete(t *testing.T) {
	var q queue
	q.append(processA)
	q.append(processB)
	q.append(processC)
	q.delete(1)

	expected := []*memory.PCB{processA, processC}

	for i := 0; i < 2; i++ {
		if !reflect.DeepEqual(q[i], expected[i]) {
			t.Errorf("expected %v, found %v", expected[i], q[i])
			break
		}
	}
}
