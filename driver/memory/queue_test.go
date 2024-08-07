package memory

import (
	"github.com/gopi-frame/queue"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockJob struct {
	queue.Job `json:"-"`

	delay time.Duration
}

func (m *mockJob) GetDelay() time.Duration {
	return m.delay
}

func (m *mockJob) Handle() error {
	return nil
}

func (m *mockJob) Failed(err error) {
	panic(err)
}

func TestQueue_Name(t *testing.T) {
	q := NewQueue("test")
	assert.Equal(t, "test", q.Name())
}

func TestQueue_Count(t *testing.T) {
	q := NewQueue("test")
	for i := 0; i < 5; i++ {
		_, ok := q.Enqueue(&mockJob{})
		assert.True(t, ok)
	}
	assert.Equal(t, int64(5), q.Count())
}

func TestQueue_Empty(t *testing.T) {
	q := NewQueue("test")
	assert.True(t, q.Empty())
	q.Enqueue(&mockJob{})
	assert.False(t, q.Empty())
}

func TestQueue_Enqueue(t *testing.T) {
	q := NewQueue("test")
	_, ok := q.Enqueue(&mockJob{})
	assert.True(t, ok)
	assert.Equal(t, int64(1), q.Count())
}

func TestQueue_Dequeue(t *testing.T) {
	t.Run("empty queue", func(t *testing.T) {
		q := NewQueue("test")
		model, ok := q.Dequeue()
		assert.Nil(t, model)
		assert.False(t, ok)
	})

	t.Run("dequeue", func(t *testing.T) {
		q := NewQueue("test")
		job := &mockJob{}
		_, ok := q.Enqueue(job)
		assert.True(t, ok)
		model, ok := q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, job, model)
		assert.Equal(t, int64(0), q.Count())
		assert.True(t, q.Empty())
	})

	t.Run("dequeue delayed job", func(t *testing.T) {
		q := NewQueue("test")
		job := &mockJob{delay: time.Second}
		_, ok := q.Enqueue(job)
		assert.True(t, ok)
		model, ok := q.Dequeue()
		assert.Nil(t, model)
		assert.False(t, ok)
		time.Sleep(time.Second)
		model, ok = q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, job, model)
		assert.Equal(t, int64(0), q.Count())
		assert.True(t, q.Empty())
	})
}

func TestQueue_Remove(t *testing.T) {
	q := NewQueue("test")
	job, ok := q.Enqueue(new(mockJob))
	assert.True(t, ok)
	assert.Equal(t, int64(1), q.Count())
	q.Remove(job)
	assert.Equal(t, int64(0), q.Count())
}

func TestQueue_Release(t *testing.T) {
	q := NewQueue("test")
	job, ok := q.Enqueue(new(mockJob))
	assert.True(t, ok)
	assert.Equal(t, 0, job.GetQueueable().GetAttempts())
	model, ok := q.Dequeue()
	assert.True(t, ok)
	q.Release(model)
	assert.Equal(t, 1, model.GetQueueable().GetAttempts())
}
