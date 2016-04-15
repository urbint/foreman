package foreman

import (
	"sync"
)

// Runner is an interface for a type which can be managed by a Foreman
type Runner interface {
	Run() error
}

// Interface is the interface of a Foreman, useful for testing
type Interface interface {
	Start(selectors ...string) error
	Abort(selectors ...string) error
	AbortAll()
	Status() map[string]string
	Register(name string, runner Runner) error

	Subscribe(chan<- Done)
}

// Foreman is a control structure to manage multiple runners
type Foreman struct {
	runners     map[string]*runnerState
	subscribers []chan<- Done

	mu sync.RWMutex
}

// New builds a Foreman for use
func New() *Foreman {
	return &Foreman{
		runners:     map[string]*runnerState{},
		subscribers: []chan<- Done{},
		mu:          sync.RWMutex{},
	}
}

// Start starts the specified runners
//
// It will error if any of the runners are already running. It does not attempt to roll-back
// started services if an error has occured
func (f *Foreman) Start(selectors ...string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, name := range selectors {
		runnerState, exists := f.runners[name]
		if !exists {
			return &ErrNoSuchRunner{name}
		}
		if err := runnerState.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Abort aborts the specified runners
//
// It will error if any of the selected runners are not currently running, or if the specified
// runners do not implement the Abortable interface
func (f *Foreman) Abort(selectors ...string) error {
	for _, name := range selectors {
		runnerState, exists := f.runners[name]
		if !exists {
			return &ErrNoSuchRunner{name}
		}
		if err := runnerState.Abort(); err != nil {
			return err
		}
	}
	return nil
}

// AbortAll will abort all currently running abortable jobs
func (f *Foreman) AbortAll() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, runnerState := range f.runners {
		asAbortable, isAbortable := runnerState.Runner.(Abortable)
		if runnerState.Status == "running" && isAbortable {
			asAbortable.Abort()
		}
	}
}

// Status reports the status of all registered
// runners in a map with a key of the runner name
// and a value of the runners status
func (f *Foreman) Status() map[string]string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := make(map[string]string)

	for name, state := range f.runners {
		result[name] = state.Status
	}

	return result
}

// Register registers a runner with the foreman
func (f *Foreman) Register(name string, runner Runner) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, alreadyExists := f.runners[name]; alreadyExists {
		return &ErrRunnerAlreadyRegistered{name}
	}

	state := newRunnerState(name, runner, f)
	f.runners[name] = state
	return nil
}
