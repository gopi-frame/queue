# Overview
[![Go Reference](https://pkg.go.dev/badge/github.com/gopi-frame/queue/driver/database.svg)](https://pkg.go.dev/github.com/gopi-frame/queue/driver/database)
[![Go](https://github.com/gopi-frame/queue/actions/workflows/go.yml/badge.svg)](https://github.com/gopi-frame/queue/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gopi-frame/queue/graph/badge.svg?token=N2LZNDNDCT&flag=database)](https://codecov.io/gh/gopi-frame/queue?flags[0]=database)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopi-frame/queue/driver/database)](https://goreportcard.com/report/github.com/gopi-frame/queue/driver/database)
[![Mit License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Package `database` is a database backed implementation of 
the [queue](https://pkg.go.dev/github.com/gopi-frame/contract/queue) interface

# Installation

```shell
go get -u -v github.com/gopi-frame/queue/driver/database
```

# Import

```go
import "github.com/gopi-frame/queue/driver/database"
```

# Usage

```go
package main

import (
    "fmt"
    "github.com/gopi-frame/queue/database"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "time"
)

type CustomJob struct {
    queue.Job `json:"-"`
}

func (c *CustomJob) Handle() error {
    // Do something
}

func (c *CustomJob) Failed(err error) {
    // Handle failed job
}

func main() {
    db, err := gorm.Open(sqlite.Open("queue.db"))
    if err != nil {
        panic(err)
    }
    q := database.NewQueue(&database.Config{
        DB:   db,
        Name: "queue",
        Job:  new(CustomJob),
    })
    q.Enqueue(new(CustomJob))
    q.Enqueue(new(CustomJob))
    q.Enqueue(new(CustomJob))
    fmt.Println("count:", q.Count()) // Output: count: 3
    for {
        if job, ok := q.Dequeue(); ok {
            if err := job.Handle(); err != nil {
                if job.GetQueueable().GetAttempts() < job.GetMaxAttempts() {
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