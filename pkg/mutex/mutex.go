package mutex

type Process int

type Resource struct {
	ownerProcess     Process
	blockedProcesses []Process
}

type Mutex struct {
	resources map[string]Resource
}

func NewMutex() Mutex {
	return Mutex{resources: map[string]Resource{}}
}

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
