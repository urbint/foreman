package foreman

import "fmt"

// runnerState is a pairing of a runner and its current State
type runnerState struct {
	Name    string
	Runner  Runner
	Status  string
	Aborted bool
	Foreman *Foreman
}

// newRunnerState builds a runner state with the specified name and runner
func newRunnerState(name string, runner Runner, foreman *Foreman) *runnerState {
	return &runnerState{
		Name:    name,
		Runner:  runner,
		Status:  "idle",
		Aborted: false,
		Foreman: foreman,
	}
}

// Start starts the runner state
func (r *runnerState) Start() error {
	if r.Status == "running" {
		return &ErrRunnerAlreadyRunning{r.Name}
	}

	r.Aborted = false
	r.Status = "running"
	go func() {
		err := r.Runner.Run()
		if err != nil {
			r.Status = fmt.Sprintf("errored: %s", err.Error())
		} else if r.Aborted {
			r.Status = "aborted"
		} else {
			r.Status = "idle"
		}
		r.Foreman.broadcastDone(Done{r.Name, err})
	}()
	return nil
}

func (r *runnerState) Abort() error {
	if r.Status != "running" {
		return &ErrRunnerNotRunning{r.Name}
	}

	if asAbortable, isAbortable := r.Runner.(Abortable); isAbortable {
		r.Aborted = true
		asAbortable.Abort()
		return nil
	}
	return &ErrNotAbortable{r.Name}
}
