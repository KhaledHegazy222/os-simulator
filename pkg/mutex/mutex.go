package mutex

type Process int

type Resource struct {
	owner_process     Process
	blocked_processes []Process
}

var resources = map[string]Resource{}

func semWait(target_resource string, process Process) {
	// Lock
	resource, ok := resources[target_resource]
	if !ok {
		// Lock Resource
		resources[target_resource] = Resource{owner_process: process, blocked_processes: []Process{}}
	} else {
		// TODO:
		// Notify os to add to blocked queue

		resource.blocked_processes = append(resource.blocked_processes, process)
	}
}

func semSignal(target_resource string) {
	// Release
	rr := resources[target_resource]
	for process := range rr.blocked_processes {
		// Update Process State from block to ready
	}
	// Remove Used Resource
	delete(resources, target_resource)
}
