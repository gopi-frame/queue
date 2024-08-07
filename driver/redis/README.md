# Overview

[![Go Reference](https://pkg.go.dev/badge/github.com/gopi-frame/queue/driver/redis.svg)](https://pkg.go.dev/github.com/gopi-frame/queue/driver/redis)
[![Go](https://github.com/gopi-frame/queue/actions/workflows/go.yml/badge.svg)](https://github.com/gopi-frame/queue/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gopi-frame/queue/graph/badge.svg?token=N2LZNDNDCT&flag=redis)](https://codecov.io/gh/gopi-frame/queue?flags[0]=redis)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopi-frame/queue/driver/redis)](https://goreportcard.com/report/github.com/gopi-frame/queue/driver/redis)
[![Mit License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Package `redis` is a redis backed implementation of the [queue](https://pkg.go.dev/github.com/gopi-frame/contract/queue)
interface

# Installation

```shell
go get -u -v github.com/gopi-frame/queue/driver/redis
```

# Import

```go
import "github.com/gopi-frame/queue/redis"
```

# Usage

```go
package main

import (
    "context"
    "github.com/gopi-frame/queue"
    "github.com/gopi-frame/queue/driver/redis"
    redislib "github.com/redis/go-redis/v9"
)

type CustomJob struct {
    queue.Job `json:"-"`
}

func (c *CustomJob) Handle() error {
    // do something
    return nil
}

func (c *CustomJob) Failed(err error) {
    // handle failed job
}

func main() {
    db := redislib.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    if err := db.Ping(context.Background()).Err(); err != nil {
        panic(err)
    }
    q := redis.NewQueue(&redis.Config{
        DB: db,
        Name: "queue",
        Job: new(CustomJob),
    })
    q.Enqueue(new(CustomJob))
    q.Enqueue(new(CustomJob))
    q.Enqueue(new(CustomJob))
    fmt.Println("count:", q.Count()) // Output: count: 3
    for {
        if job, ok := q.Dequeue(); ok {
            if err := job.Handle(); err != nil {
                if job.GetQueueable().GetAttempts < job.GetMaxAttempts() {
                    q.Release(job)
                } else {
                    job.Failed(err)
                }
            } else {
                q.Ack(job)
            }
        } else {
            time.Sleep(time.Millisecond * 100)
        }
    }
}
```