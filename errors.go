package foreman

import "fmt"

// ErrNoSuchRunner is an error that is emitted
// when a runner is requested that does not exist
type ErrNoSuchRunner struct {
	name string
}

func (e *ErrNoSuchRunner) Error() string {
	return fmt.Sprintf("No such runner: Runner %s does not exist", e.name)
}

// ErrRunnerAlreadyRunning is an error that is emitted
// when attempting to start an already running runner
type ErrRunnerAlreadyRunning struct {
	name string
}

func (e *ErrRunnerAlreadyRunning) Error() string {
	return fmt.Sprintf("Already running: Runner %s has already been started", e.name)
}

// ErrRunnerNotRunning is an error that is emitted
// when attempting to abort a non-running runner
type ErrRunnerNotRunning struct {
	name string
}

func (e *ErrRunnerNotRunning) Error() string {
	return fmt.Sprintf("Not running: Runner %s is not running", e.name)
}

// ErrNotAbortable is an error that is emitted
// when attempting to abort a non-abortable runner
type ErrNotAbortable struct {
	name string
}

func (e *ErrNotAbortable) Error() string {
	return fmt.Sprintf("Not abortable: Runner %s does not implement the Abortable interface", e.name)
}

// ErrRunnerAlreadyRegistered is an error that is emitted
// when attempting to register a runner with a name
// that has already been used
type ErrRunnerAlreadyRegistered struct {
	name string
}

func (e *ErrRunnerAlreadyRegistered) Error() string {
	return fmt.Sprintf("Already registered: the name %s is already in use by another runner", e.name)
}
