package event

// event topics
const (
	TopicQueueBeforeRun  = "event:queue:before:run"
	TopicQueueBeforeStop = "event:queue:before:stop"
	TopicQueueAfterStop  = "event:queue:after:stop"
	TopicJobBeforeHandle = "event:job:before:handle"
	TopicJobAfterHandle  = "event:job:after:handle"
	TopicJobAfterRelease = "event:job:after:release"
)
