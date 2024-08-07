# Overview
[![Go Reference](https://pkg.go.dev/badge/github.com/gopi-frame/queue/driver/memory.svg)](https://pkg.go.dev/github.com/gopi-frame/queue/driver/memory)
[![Go](https://github.com/gopi-frame/queue/actions/workflows/go.yml/badge.svg)](https://github.com/gopi-frame/queue/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gopi-frame/queue/graph/badge.svg?token=N2LZNDNDCT&flag=memory)](https://codecov.io/gh/gopi-frame/queue?flags[0]=memory)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopi-frame/queue/driver/memory)](https://goreportcard.com/report/github.com/gopi-frame/queue/driver/memory)
[![Mit License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Package `memory` is a memory based implementation of 
the [queue](https://pkg.go.dev/github.com/gopi-frame/contract/queue)

# Installation
```shell
go get -u -v github.com/gopi-frame/queue/driver/memory
```

# Import
```go
import "github.com/gopi-frame/queue/driver/memory"
```

# Usage

```go
package main

import (
    "github.com/gopi-frame/queue"
    "github.com/gopi-frame/queue/driver/memory"
)

type CustomJob struct {
    queue.Job
}

func (c *CustomJob) Handle() error {
    // do something
    return nil
}

func (c *CustomJob) Failed(err error) {
    // handle failed job
}

func main() {
    q := memory.NewQueue("test")
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
