package mutex

import (
	"testing"
)

func TestMutex(t *testing.T) {

	t.Run("test single process acquire lock", func(t *testing.T) {
		const firstProcess Process = 1
		targetResource := "userInput"
		mutex := NewMutex()
		success := mutex.SemWait(targetResource, firstProcess)
		if !success {
			t.Fatalf("Failed to acquire the lock for resource %q to process %q\n", targetResource, firstProcess)
		}
		resource, isPresent := mutex.resources[targetResource]
		if !isPresent {
			t.Fatalf("Expected to get process %q\n", firstProcess)
		}
		if resource.ownerProcess != firstProcess {
			t.Fatalf("Expected process %q, found %q\n", firstProcess, resource.ownerProcess)
		}
	})
	t.Run("test multiple processes acquire lock return value", func(t *testing.T) {
		mutex := NewMutex()
		const firstProcess Process = 1
		const secondProcess Process = 2
		targetResource := "userInput"
		success := mutex.SemWait(targetResource, firstProcess)
		if !success {
			t.Fatalf("Failed to acquire the lock for resource %q to process %q\n", targetResource, firstProcess)
		}
		success = mutex.SemWait(targetResource, secondProcess)
		if success {
			t.Fatalf("Expected to be blocked when acquire the lock for resource %q to process %q\n", targetResource, secondProcess)
		}

	})
	t.Run("test multiple processes acquire lock blocked queues", func(t *testing.T) {
		mutex := NewMutex()
		const firstProcess Process = 1
		const secondProcess Process = 2
		targetResource := "userInput"
		mutex.SemWait(targetResource, firstProcess)
		mutex.SemWait(targetResource, secondProcess)

		resource, isPresent := mutex.resources[targetResource]
		if !isPresent {
			t.Fatalf("Expected to get process %q\n", firstProcess)
		}

		expectedProcess := firstProcess
		actualProcess := resource.ownerProcess
		if resource.ownerProcess != firstProcess {
			t.Fatalf("Expected process %q, found %q\n", expectedProcess, actualProcess)
		}

		expectedLength := 1
		actualLength := len(resource.blockedProcesses)
		if expectedLength != actualLength {
			t.Fatalf("Expected to have %d blocked processes, but found %d\n", expectedLength, actualLength)
		}
		expectedProcess = secondProcess
		actualProcess = resource.blockedProcesses[0]
		if expectedProcess != actualProcess {
			t.Fatalf("Expected process %q, found %q\n", expectedProcess, actualProcess)
		}
	})
	t.Run("test single process release lock", func(t *testing.T) {
		const firstProcess Process = 1
		targetResource := "userInput"
		mutex := NewMutex()
		mutex.SemWait(targetResource, firstProcess)
		success := mutex.SemSignal(targetResource, firstProcess)
		if !success {
			t.Fatalf("Failed to release the lock for resource %q to process %d\n", targetResource, firstProcess)
		}
	})
	t.Run("test single process multiple release lock", func(t *testing.T) {
		const firstProcess Process = 1
		targetResource := "userInput"
		mutex := NewMutex()
		mutex.SemWait(targetResource, firstProcess)
		mutex.SemSignal(targetResource, firstProcess)
		success := mutex.SemSignal(targetResource, firstProcess)
		if success {
			t.Fatalf("Expected to fail to release the lock for resource %q to process %q\n", targetResource, firstProcess)
		}
	})

	t.Run("test not allowed release lock", func(t *testing.T) {
		const firstProcess Process = 1
		const secondProcess Process = 2
		targetResource := "userInput"
		mutex := NewMutex()
		mutex.SemWait(targetResource, firstProcess)
		success := mutex.SemSignal(targetResource, secondProcess)

		if success {
			t.Fatalf("Expected to fail to release the lock for resource %q to process %q\n", targetResource, secondProcess)
		}
	})

	t.Run("test consecutive locks and releases", func(t *testing.T) {
		mutex := NewMutex()
		const firstProcess Process = 1
		const secondProcess Process = 2
		targetResource := "userInput"
		mutex.SemWait(targetResource, firstProcess)
		mutex.SemSignal(targetResource, firstProcess)

		success := mutex.SemWait(targetResource, secondProcess)
		if !success {
			t.Fatalf("Failed to acquire the lock for resource %q to process %q\n", targetResource, secondProcess)
		}

		success = mutex.SemSignal(targetResource, secondProcess)
		if !success {
			t.Fatalf("Failed to release the lock for resource %q to process %d\n", targetResource, secondProcess)
		}

	})
}
