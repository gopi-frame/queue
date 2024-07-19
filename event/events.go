package event

const (
	dispatching    = "queue:job:dispatching"
	dispatched     = "queue:job:dispatched"
	workerLooping  = "queue:worker:looping"
	workerStopping = "queue:worker:stopping"
	workerStopped  = "queue:worker:stopped"
	jobHandling    = "queue:job:handling"
	jobHandled     = "queue:job:handled"
	jobFailed      = "queue:job:failed"
)

func Topics() []string {
	return []string{
		dispatched,
		dispatched,
		workerLooping,
		workerStopping,
		workerStopped,
		jobHandling,
		jobHandled,
		jobFailed,
	}
}
