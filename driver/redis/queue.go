package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	queuecontract "github.com/gopi-frame/contract/queue"
	rediscontract "github.com/gopi-frame/contract/redis"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strings"
	"time"
)

// QueueJobIDKeyFormat is the redis key format for storing job id in a stored set
const QueueJobIDKeyFormat = "GOPI:QUEUE:INDEX:%s"

// QueueJobItemKeyFormat is the redis key format for storing job hash in a hash table
const QueueJobItemKeyFormat = "GOPI:QUEUE:ITEMS:%s"

// Queue is a redis queue implementation
type Queue struct {
	name   string
	client rediscontract.Client
	job    reflect.Type
}

// NewQueue creates a new redis queue
func NewQueue(cfg *Config, opts ...Option) *Queue {
	if err := cfg.Apply(opts...); err != nil {
		panic(err)
	}
	if err := cfg.Valid(); err != nil {
		panic(err)
	}
	return &Queue{
		name:   cfg.Name,
		client: cfg.DB,
		job:    reflect.Indirect(reflect.ValueOf(cfg.Job)).Type(),
	}
}

func (q *Queue) idKey() string {
	return fmt.Sprintf(QueueJobIDKeyFormat, strings.ToUpper(q.name))
}

func (q *Queue) itemKey() string {
	return fmt.Sprintf(QueueJobItemKeyFormat, strings.ToUpper(q.name))
}

// Name returns the queue name
func (q *Queue) Name() string {
	return q.name
}

// Empty returns true if queue is empty
func (q *Queue) Empty() bool {
	return q.Count() == 0
}

// Count returns the number of jobs in the queue
func (q *Queue) Count() int64 {
	value, err := q.client.ZCard(context.Background(), q.idKey()).Uint64()
	if errors.Is(err, redis.Nil) {
		return 0
	}
	if err != nil {
		panic(err)
	}
	return int64(value)
}

// Enqueue adds a job to the queue
func (q *Queue) Enqueue(job queuecontract.Job) (queuecontract.Job, bool) {
	model := NewJob(q.name, job)
	script := redis.NewScript(new(LuaScript).Enqueue())
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		panic(err)
	}
	if err := script.Run(
		context.Background(),
		q.client,
		[]string{q.idKey(), q.itemKey()},
		model.AvailableAt.UnixNano(),
		model.ID.String(),
		string(jsonBytes),
	).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return model.GetPayload(), true
		}
		defer func() {
			if err := q.client.ZRem(context.Background(), q.idKey(), model.ID.String()).Err(); err != nil {
				panic(err)
			}
		}()
		panic(err)
	}
	return model.GetPayload(), true
}

// Dequeue removes and returns a job from the queue
func (q *Queue) Dequeue() (queuecontract.Job, bool) {
	script := redis.NewScript(new(LuaScript).Dequeue())
	result, err := script.Run(
		context.Background(),
		q.client,
		[]string{q.idKey(), q.itemKey()},
		1,
		time.Now().UnixNano(),
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false
		}
		panic(err)
	}
	var dest struct {
		ID          uuid.UUID       `json:"id"`
		Queue       string          `json:"queue"`
		Payload     json.RawMessage `json:"payload"`
		Attempts    int             `json:"attempts"`
		AvailableAt time.Time       `json:"available_at"`
		CreatedAt   time.Time       `json:"created_at"`
	}
	if err := json.Unmarshal([]byte(result.(string)), &dest); err != nil {
		panic(err)
	}
	var model = new(Job)
	model.ID = dest.ID
	model.Queue = dest.Queue
	model.Attempts = dest.Attempts
	model.AvailableAt = dest.AvailableAt
	var payload = reflect.New(q.job).Interface()
	if err := json.Unmarshal(dest.Payload, payload); err != nil {
		panic(err)
	}
	model.Payload = payload.(queuecontract.Job)
	model.Payload.SetQueueable(model)
	return model.Payload, true
}

// Remove removes a job from the queue
func (q *Queue) Remove(job queuecontract.Job) {
	if job.GetQueueable() == nil {
		return
	}
	script := redis.NewScript(new(LuaScript).Remove())
	if err := script.Run(
		context.Background(),
		q.client,
		[]string{q.idKey(), q.itemKey()},
		job.GetQueueable().GetID(),
	).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return
		}
		panic(err)
	}
}

// Ack acknowledges a job
func (q *Queue) Ack(_ queuecontract.Job) {}

// Release releases a job back to the queue
func (q *Queue) Release(job queuecontract.Job) {
	if job.GetQueueable() == nil {
		return
	}
	model := job.GetQueueable().(*Job)
	model.Attempts += 1
	model.AvailableAt = time.Now()
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		panic(err)
	}
	script := redis.NewScript(new(LuaScript).Enqueue())
	if err := script.Run(
		context.Background(),
		q.client,
		[]string{q.idKey(), q.itemKey()},
		model.AvailableAt.UnixNano(),
		model.ID.String(),
		string(jsonBytes),
	).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return
		}
		panic(err)
	}
}
