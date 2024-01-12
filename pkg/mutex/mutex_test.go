package mutex

import (
	"testing"
)

func TestMutex(t *testing.T) {

	t.Run("test_single_process_acquire_lock", func(t *testing.T) {
		const first_process Process = 1
		target_resource := "userInput"
		mutex := NewMutex()
		success := mutex.SemWait(target_resource, first_process)
		if !success {
			t.Fatalf("Failed to acquire the lock for resource %q to process %q\n", target_resource, first_process)
		}
		resource, isPresent := mutex.resources[target_resource]
		if !isPresent {
			t.Fatalf("Expected to get process %q\n", first_process)
		}
		if resource.ownerProcess != first_process {
			t.Fatalf("Expected process %q, found %q\n", first_process, resource.ownerProcess)
		}
	})
	t.Run("test_multiple_processes_acquire_lock_return_value", func(t *testing.T) {
		mutex := NewMutex()
		const first_process Process = 1
		const second_process Process = 2
		target_resource := "userInput"
		success := mutex.SemWait(target_resource, first_process)
		if !success {
			t.Fatalf("Failed to acquire the lock for resource %q to process %q\n", target_resource, first_process)
		}
		success = mutex.SemWait(target_resource, second_process)
		if success {
			t.Fatalf("Expected to be blocked when acquire the lock for resource %q to process %q\n", target_resource, second_process)
		}

	})
	t.Run("test_multiple_processes_acquire_lock_blocked_queues", func(t *testing.T) {
		mutex := NewMutex()
		const first_process Process = 1
		const second_process Process = 2
		target_resource := "userInput"
		mutex.SemWait(target_resource, first_process)
		mutex.SemWait(target_resource, second_process)

		resource, isPresent := mutex.resources[target_resource]
		if !isPresent {
			t.Fatalf("Expected to get process %q\n", first_process)
		}

		expected_process := first_process
		actual_process := resource.ownerProcess
		if resource.ownerProcess != first_process {
			t.Fatalf("Expected process %q, found %q\n", expected_process, actual_process)
		}

		expected_length := 1
		actual_length := len(resource.blocked_processes)
		if expected_length != actual_length {
			t.Fatalf("Expected to have %d blocked processes, but found %d\n", expected_length, actual_length)
		}
		expected_process = second_process
		actual_process = resource.blocked_processes[0]
		if expected_process != actual_process {
			t.Fatalf("Expected process %q, found %q\n", expected_process, actual_process)
		}
	})
	t.Run("test_single_process_release_lock", func(t *testing.T) {
		const first_process Process = 1
		target_resource := "userInput"
		mutex := NewMutex()
		mutex.SemWait(target_resource, first_process)
		success := mutex.SemSignal(target_resource, first_process)
		if !success {
			t.Fatalf("Failed to release the lock for resource %q to process %d\n", target_resource, first_process)
		}
	})
	t.Run("test_single_process_multiple_release_lock", func(t *testing.T) {
		const first_process Process = 1
		target_resource := "userInput"
		mutex := NewMutex()
		mutex.SemWait(target_resource, first_process)
		mutex.SemSignal(target_resource, first_process)
		success := mutex.SemSignal(target_resource, first_process)
		if success {
			t.Fatalf("Expected to fail to release the lock for resource %q to process %q\n", target_resource, first_process)
		}
	})

	t.Run("test_not_allowed_release_lock", func(t *testing.T) {
		const first_process Process = 1
		const second_process Process = 2
		target_resource := "userInput"
		mutex := NewMutex()
		mutex.SemWait(target_resource, first_process)
		success := mutex.SemSignal(target_resource, second_process)

		if success {
			t.Fatalf("Expected to fail to release the lock for resource %q to process %q\n", target_resource, second_process)
		}
	})

	t.Run("test_consecutive_locks_and_releases", func(t *testing.T) {
		mutex := NewMutex()
		const first_process Process = 1
		const second_process Process = 2
		target_resource := "userInput"
		mutex.SemWait(target_resource, first_process)
		mutex.SemSignal(target_resource, first_process)

		success := mutex.SemWait(target_resource, second_process)
		if !success {
			t.Fatalf("Failed to acquire the lock for resource %q to process %q\n", target_resource, second_process)
		}

		success = mutex.SemSignal(target_resource, second_process)
		if !success {
			t.Fatalf("Failed to release the lock for resource %q to process %d\n", target_resource, second_process)
		}

	})
}
