package foreman

// Abortable is an interface for a runner which can be aborted
type Abortable interface {
	Runner
	Abort()
}

// FnAborter is an wrapper struct for a runner that allows
// specifying an abort function at construction
type FnAborter struct {
	Runner
	AbortFn func()
}

// AbortAdapter builds a FnAborter for the specified Runner
func AbortAdapter(runner Runner, abort func()) *FnAborter {
	return &FnAborter{
		Runner:  runner,
		AbortFn: abort,
	}
}

// Abort implements Abortable for the FnAborter and will call the Fn specified
// at construction
func (f *FnAborter) Abort() {
	f.AbortFn()
}
