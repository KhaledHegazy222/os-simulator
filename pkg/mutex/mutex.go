// Package mutex provides a simple implementation of a mutex using semaphores.
package mutex

// Process represents a process identifier.
type Process int

// Resource represents a resource with an owner process and a list of blocked processes.
type Resource struct {
	ownerProcess     Process
	blockedProcesses []Process
}

// Mutex represents a mutex that manages resources and provides mutual exclusion.
type Mutex struct {
	resources map[string]Resource
}

// NewMutex creates a new Mutex instance.
func NewMutex() Mutex {
	return Mutex{resources: map[string]Resource{}}
}

// SemWait acquires a lock on the specified resource.
// If the resource is already locked, the calling process is added to the list of blocked processes.
// Returns true if the lock is acquired, false otherwise.
func (m *Mutex) SemWait(targetResource string, process Process) bool {
	// Lock
	resource, isPresent := m.resources[targetResource]
	if !isPresent {
		// Lock Resource
		m.resources[targetResource] = Resource{ownerProcess: process, blockedProcesses: []Process{}}
		return true
	} else {
		// TODO:
		// Notify os to add to blocked queue

		resource.blockedProcesses = append(resource.blockedProcesses, process)
		m.resources[targetResource] = resource
		return false
	}
}

// SemSignal releases the lock on the specified resource.
// If the calling process is the owner of the resource, it releases the lock and updates the state of blocked processes.
// Returns true if the lock is released, false otherwise.
func (m *Mutex) SemSignal(targetResource string, process Process) bool {
	// Release
	resource, isPresent := m.resources[targetResource]
	if isPresent && resource.ownerProcess == process {
		for range resource.blockedProcesses {
			// Update Process State from block to ready
		}
		// Remove Used Resource
		delete(m.resources, targetResource)
		return true
	} else {
		return false
	}
}
