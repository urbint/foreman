package foreman

// Done is the struct emitted by Foreman to Subscribers
//
// It includes the name of the runner, and the error, if any
type Done struct {
	Name  string
	Error error
}

// Subscribe subscribes to Foreman for notifications
// regarding the finishing of jobs
func (f *Foreman) Subscribe(out chan<- Done) {
	f.subscribers = append(f.subscribers, out)
}

// broadcastDone sends Done notifications to all subscribers
func (f *Foreman) broadcastDone(done Done) {
	for _, sub := range f.subscribers {
		sub <- done
	}
}
