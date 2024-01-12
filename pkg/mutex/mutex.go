package mutex

type Process int

type Resource struct {
	ownerProcess      Process
	blocked_processes []Process
}

type Mutex struct {
	resources map[string]Resource
}

func NewMutex() Mutex {
	return Mutex{resources: map[string]Resource{}}
}

func (m *Mutex) SemWait(target_resource string, process Process) bool {
	// Lock
	resource, isPresent := m.resources[target_resource]
	if !isPresent {
		// Lock Resource
		m.resources[target_resource] = Resource{ownerProcess: process, blocked_processes: []Process{}}
		return true
	} else {
		// TODO:
		// Notify os to add to blocked queue

		resource.blocked_processes = append(resource.blocked_processes, process)
		m.resources[target_resource] = resource
		return false
	}
}

func (m *Mutex) SemSignal(target_resource string, process Process) bool {
	// Release
	resource, isPresent := m.resources[target_resource]
	if isPresent && resource.ownerProcess == process {
		for process := range resource.blocked_processes {
			// Update Process State from block to ready
		}
		// Remove Used Resource
		delete(m.resources, target_resource)
		return true
	} else {
		return false
	}

}
