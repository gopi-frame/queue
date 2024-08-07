# Overview
[![Go Reference](https://pkg.go.dev/badge/github.com/gopi-frame/queue.svg)](https://pkg.go.dev/github.com/gopi-frame/queue)
[![Go](https://github.com/gopi-frame/queue/actions/workflows/go.yml/badge.svg)](https://github.com/gopi-frame/queue/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gopi-frame/queue/graph/badge.svg?token=N2LZNDNDCT&flag=queue)](https://codecov.io/gh/gopi-frame/queue?flags[0]=queue)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopi-frame/queue)](https://goreportcard.com/report/github.com/gopi-frame/queue)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Package queue provides a task queue implementation.

# Installation
```shell
go get -u -v github.com/gopi-frame/queue
```

# Import
```go
import "github.com/gopi-frame/queue"
```

# Usage

```go
package main

import (
    "github.com/gopi-frame/queue"
    "sync"
    _ "github.com/gopi-frame/queue/driver/memory"
    "time"

    //_ "github.com/gopi-frame/queue/driver/database"
    //_ "github.com/gopi-frame/queue/driver/redis"
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
    driver, err := queue.Open("memory", map[string]any{
        "name": "test",
    })
    if err != nil {
        panic(err)
    }
    q := queue.NewQueue(driver, queue.WorkerNum(3))
    go q.Run()
    for i := 0; i < 1000; i++ {
        q.Enqueue(new(CustomJob))
    }
    for {
        if q.Empty() {
            q.Stop()
            break
        }
        time.Sleep(time.Second)
    }
}
```

# Drivers
- [database](./driver/database/README.md)
- [memory](./driver/memory/README.md)
- [redis](./driver/redis/README.md)

# How to create custom driver
Just implement the [queue.Driver](https://pkg.go.dev/github.com/gopi-frame/contract/queue#Driver) interface
and register it with `queue.Register`.
Then you can use your driver in the following way:
```go
driver, err := queue.Open("custom", map[string]any{
    // options
})
```