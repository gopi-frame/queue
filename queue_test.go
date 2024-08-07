package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQueue_Run(t *testing.T) {
	q := NewQueue(&mockQueue{name: "test"}, 3)
	go q.Run()
	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, 3, len(q.workers))

	t.Run("success", func(t *testing.T) {
		job := new(mockJob)
		q.Enqueue(job)
		time.Sleep(time.Millisecond * 200)
		assert.True(t, job.handled)
	})

	t.Run("failed", func(t *testing.T) {
		job := new(mockJob)
		job.shouldFailTimes = 5
		q.Enqueue(job)
		time.Sleep(time.Millisecond * 200)
		assert.False(t, job.handled)
		assert.True(t, job.failed)
		assert.Equal(t, 3, job.GetQueueable().GetAttempts())
		assert.Equal(t, "fail", job.failedReason.Error())
	})

	t.Run("success after failed", func(t *testing.T) {
		job := new(mockJob)
		job.shouldFailTimes = 2
		q.Enqueue(job)
		for {
			if job.handled {
				assert.Equal(t, 2, job.GetQueueable().GetAttempts())
				break
			}
		}
	})

	t.Run("delayed job", func(t *testing.T) {
		job := new(mockJob)
		job.delay = time.Millisecond * 200
		start := time.Now()
		q.Enqueue(job)
		for {
			if job.handled {
				assert.GreaterOrEqual(t, time.Since(start), 200*time.Millisecond)
				break
			}
		}
	})

	t.Run("retry delay job", func(t *testing.T) {
		job := new(mockJob)
		job.retryDelay = 100 * time.Millisecond
		job.shouldFailTimes = 2
		start := time.Now()
		q.Enqueue(job)
		for {
			if job.handled {
				assert.GreaterOrEqual(t, time.Since(start), 200*time.Millisecond)
				break
			}
		}
	})
}

func TestQueue_Stop(t *testing.T) {
	q := NewQueue(&mockQueue{name: "test"}, 3)
	go q.Run()
	time.Sleep(time.Millisecond * 200)
	q.Stop()
	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, 0, len(q.workers))
}

func TestQueue_Uptime(t *testing.T) {
	q := NewQueue(&mockQueue{name: "test"}, 3)
	assert.Equal(t, time.Duration(0), q.Uptime())
	go q.Run()
	time.Sleep(time.Millisecond * 200)
	q.Stop()
	assert.True(t, q.Uptime() > time.Millisecond*200)
}
