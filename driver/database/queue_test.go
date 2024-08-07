package database

import (
	"fmt"
	"github.com/gopi-frame/queue"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "test.db")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	if err := testDB.AutoMigrate(&Job{}); err != nil {
		panic(err)
	}
	m.Run()
	flush("test")
}

type testJob struct {
	queue.Job `json:"-"`

	Message string        `json:"message"`
	Delay   time.Duration `json:"delay"`
}

func (t *testJob) GetDelay() time.Duration {
	return t.Delay
}

func (t *testJob) Handle() error {
	return nil
}

func (t *testJob) Failed(err error) {
	panic(err)
}

func TestNewQueue(t *testing.T) {
	t.Run("without DB", func(t *testing.T) {
		assert.Panics(t, func() {
			NewQueue(&Config{
				Name: "test",
				Job:  new(testJob),
			})
		})
	})

	t.Run("without Name", func(t *testing.T) {
		assert.Panics(t, func() {
			NewQueue(&Config{
				DB:  testDB,
				Job: new(testJob),
			})
		})
	})

	t.Run("without Job", func(t *testing.T) {
		assert.Panics(t, func() {
			NewQueue(&Config{
				DB:   testDB,
				Name: "test",
			})
		})
	})

	t.Run("without Table", func(t *testing.T) {
		q := NewQueue(&Config{
			DB:   testDB,
			Name: "test",
			Job:  new(testJob),
		})
		assert.Equal(t, "jobs", q.table)
	})

	t.Run("full config", func(t *testing.T) {
		q := NewQueue(&Config{
			DB:    testDB,
			Name:  "test",
			Job:   new(testJob),
			Table: "jobs2",
		})
		assert.Equal(t, "jobs2", q.table)
		assert.Equal(t, "test", q.Name())
		assert.Equal(t, reflect.Indirect(reflect.ValueOf(new(testJob))).Type(), q.job)
	})
}

func TestQueue_Enqueue(t *testing.T) {
	flush("test")
	q := NewQueue(&Config{
		DB:    testDB,
		Name:  "test",
		Job:   new(testJob),
		Table: DefaultJobTable,
	})
	job, ok := q.Enqueue(&testJob{Message: "Hello world"})
	assert.True(t, ok)
	assert.Equal(t, int64(1), q.Count())
	assert.False(t, q.Empty())
	assert.Equal(t, "test", job.GetQueueable().GetQueue())
	assert.Equal(t, 0, job.GetQueueable().GetAttempts())
	assert.Equal(t, job, job.GetQueueable().GetPayload())
}

func TestQueue_Dequeue(t *testing.T) {
	flush("test")
	q := NewQueue(&Config{
		DB:   testDB,
		Name: "test",
		Job:  new(testJob),
	})
	t.Run("empty queue", func(t *testing.T) {
		model, ok := q.Dequeue()
		assert.Nil(t, model)
		assert.False(t, ok)
	})

	t.Run("dequeue", func(t *testing.T) {
		q.Enqueue(&testJob{Message: "Hello world"})
		job, ok := q.Dequeue()
		assert.Equal(t, "Hello world", job.(*testJob).Message)
		assert.True(t, ok)
		assert.Equal(t, int64(0), q.Count())
		assert.True(t, q.Empty())
	})

	t.Run("dequeue delayed Job", func(t *testing.T) {
		q.Enqueue(&testJob{Message: "Hello world", Delay: time.Second})
		job, ok := q.Dequeue()
		assert.Equal(t, int64(1), q.Count())
		assert.False(t, ok)
		assert.Nil(t, job)
		time.Sleep(time.Second)
		job, ok = q.Dequeue()
		assert.Equal(t, int64(0), q.Count())
		assert.True(t, ok)
		assert.True(t, q.Empty())
	})
}

func TestQueue_Remove(t *testing.T) {
	flush("test")
	q := NewQueue(&Config{
		DB:   testDB,
		Name: "test",
		Job:  new(testJob),
	})
	job1, ok := q.Enqueue(&testJob{Message: "Hello world"})
	assert.True(t, ok)
	job2, ok := q.Enqueue(&testJob{Message: "Hello world"})
	assert.True(t, ok)
	job3, ok := q.Enqueue(&testJob{Message: "Hello world"})
	assert.True(t, ok)
	assert.Equal(t, int64(3), q.Count())
	q.Remove(job1)
	assert.Equal(t, int64(2), q.Count())
	var dest1 = make(map[string]any)
	var dest2 = make(map[string]any)
	var dest3 = make(map[string]any)
	assert.Error(t, gorm.ErrRecordNotFound, testDB.Table(q.table).Where("id =?", job1.GetQueueable().GetID()).Take(&dest1).Error)
	assert.Nil(t, testDB.Table(q.table).Where("id =?", job2.GetQueueable().GetID()).Take(&dest2).Error)
	assert.Nil(t, testDB.Table(q.table).Where("id =?", job3.GetQueueable().GetID()).Take(&dest3).Error)
}

func TestQueue_Ack(t *testing.T) {
	flush("test")
	q := NewQueue(&Config{
		DB:   testDB,
		Name: "test",
		Job:  new(testJob),
	})
	job1, ok := q.Enqueue(&testJob{Message: "Hello world"})
	assert.True(t, ok)
	job2, ok := q.Enqueue(&testJob{Message: "Hello world"})
	assert.True(t, ok)
	job3, ok := q.Enqueue(&testJob{Message: "Hello world"})
	assert.True(t, ok)
	assert.Equal(t, int64(3), q.Count())
	q.Ack(job1)
	assert.Equal(t, int64(2), q.Count())
	var dest1 = make(map[string]any)
	var dest2 = make(map[string]any)
	var dest3 = make(map[string]any)
	assert.Error(t, gorm.ErrRecordNotFound, testDB.Table(q.table).Where("id =?", job1.GetQueueable().GetID()).Take(&dest1).Error)
	assert.Nil(t, testDB.Table(q.table).Where("id =?", job2.GetQueueable().GetID()).Take(&dest2).Error)
	assert.Nil(t, testDB.Table(q.table).Where("id =?", job3.GetQueueable().GetID()).Take(&dest3).Error)
}

func TestQueue_Release(t *testing.T) {
	flush("test")
	q := NewQueue(&Config{
		DB:   testDB,
		Name: "test",
		Job:  new(testJob),
	})
	job, ok := q.Enqueue(&testJob{Message: "Hello world"})
	assert.True(t, ok)
	assert.Equal(t, int64(1), q.Count())
	originalAvailableAt := job.GetQueueable().(*Job).AvailableAt
	time.Sleep(time.Second)
	q.Release(job)
	assert.LessOrEqual(t, time.Second, job.GetQueueable().(*Job).AvailableAt.Sub(originalAvailableAt))
	assert.Equal(t, 1, job.GetQueueable().GetAttempts())
}

// For testing purpose, clear all data from database.
func flush(queue string) {
	if err := testDB.Exec(fmt.Sprintf("DELETE FROM `%s` WHERE `queue` = '%s'", DefaultJobTable, queue)).Error; err != nil {
		panic(err)
	}
}
