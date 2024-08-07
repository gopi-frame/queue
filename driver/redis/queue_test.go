package redis

import (
	"context"
	"fmt"
	rediscontract "github.com/gopi-frame/contract/redis"
	"github.com/gopi-frame/queue"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
	"time"
)

var testDB rediscontract.Client

func TestMain(m *testing.M) {
	testDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := testDB.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	m.Run()
}

// For testing purpose, clear all data from redis.
func flush(queue string) {
	if err := testDB.Del(context.Background(), fmt.Sprintf(QueueJobIDKeyFormat, strings.ToUpper(queue))).Err(); err != nil {
		panic(err)
	}
	if err := testDB.Del(context.Background(), fmt.Sprintf(QueueJobItemKeyFormat, strings.ToUpper(queue))).Err(); err != nil {
		panic(err)
	}
}

type testJob struct {
	queue.Job `json:"-"`

	delay time.Duration
}

func (t *testJob) GetDelay() time.Duration {
	return t.delay
}

func (t *testJob) Handle() error {
	return nil
}

func (t *testJob) Failed(err error) {
	fmt.Println("failed:", err)
}

func TestQueue_Enqueue(t *testing.T) {
	flush("test")
	q := NewQueue(&Config{}, WithDB(testDB), WithName("test"), WithJob(new(testJob)))
	assert.Equal(t, "test", q.Name())
	assert.True(t, q.Empty())
	job := new(testJob)
	result, ok := q.Enqueue(job)
	assert.True(t, ok)
	assert.Equal(t, int64(1), q.Count())
	assert.Equal(t, job, result)
	assert.Equal(t, "test", result.GetQueueable().GetQueue())
	assert.Equal(t, 0, result.GetQueueable().GetAttempts())
	assert.Equal(t, job, result.GetQueueable().GetPayload())
	assert.Equal(t, "test", job.GetQueueable().GetQueue())
	assert.Equal(t, 0, job.GetQueueable().GetAttempts())
}

func TestQueue_Dequeue(t *testing.T) {
	flush("test")

	t.Run("empty queue", func(t *testing.T) {
		q := NewQueue(&Config{
			DB:   testDB,
			Name: "test",
			Job:  new(testJob),
		})
		result, ok := q.Dequeue()
		assert.Nil(t, result)
		assert.False(t, ok)
	})

	t.Run("dequeue", func(t *testing.T) {
		q := NewQueue(&Config{
			DB:   testDB,
			Name: "test",
			Job:  new(testJob),
		})
		job := new(testJob)
		q.Enqueue(job)
		result, ok := q.Dequeue()
		assert.Equal(t, job.GetQueueable().GetID(), result.GetQueueable().GetID())
		assert.True(t, ok)
	})

	t.Run("dequeue delayed job", func(t *testing.T) {
		q := NewQueue(&Config{
			DB:   testDB,
			Name: "test",
			Job:  new(testJob),
		})
		job := new(testJob)
		job.delay = time.Second
		q.Enqueue(job)
		result, ok := q.Dequeue()
		assert.Nil(t, result)
		assert.False(t, ok)
		time.Sleep(time.Second)
		result, ok = q.Dequeue()
		assert.NotNil(t, result)
		assert.Equal(t, job.GetQueueable().GetID(), result.GetQueueable().GetID())
		assert.True(t, ok)
	})
}

func TestQueue_Remove(t *testing.T) {
	flush("test")
	q := NewQueue(&Config{
		DB:   testDB,
		Name: "test",
		Job:  new(testJob),
	})
	job := new(testJob)
	q.Enqueue(job)
	assert.Equal(t, int64(1), q.Count())
	q.Remove(job)
	assert.Equal(t, int64(0), q.Count())
}

func TestQueue_Release(t *testing.T) {
	flush("test")
	q := NewQueue(&Config{
		DB:   testDB,
		Name: "test",
		Job:  new(testJob),
	})
	job := new(testJob)
	q.Enqueue(job)
	assert.Equal(t, 0, job.GetQueueable().GetAttempts())
	q.Release(job)
	assert.Equal(t, 1, job.GetQueueable().GetAttempts())
}

func TestNewQueue(t *testing.T) {
	flush("test")
	t.Run("invalid config", func(t *testing.T) {
		assert.Panics(t, func() {
			NewQueue(new(Config))
		})
	})

	t.Run("invalid option", func(t *testing.T) {
		assert.Panics(t, func() {
			NewQueue(new(Config), WithJob(nil))
		})
	})

	t.Run("valid config", func(t *testing.T) {
		q := NewQueue(&Config{
			DB:   testDB,
			Name: "test",
			Job:  new(testJob),
		})
		assert.Equal(t, "test", q.Name())
		assert.Equal(t, testDB, q.client)
		assert.Equal(t, reflect.TypeFor[testJob](), q.job)
		assert.True(t, q.Empty())
	})
}
