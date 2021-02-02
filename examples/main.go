package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/raspiantoro/buruh"
)

var fn = func(id int, wg *sync.WaitGroup) buruh.Task {
	return func(ctx context.Context) {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(1)
		time.Sleep(time.Duration(n) * time.Second)

		wKey := ctx.Value(buruh.CtxWorkerKey).(string)
		jKey := ctx.Value(buruh.CtxJobKey).(string)

		fmt.Printf("Method #%d is executing by job: %s with worker: %s\n", id, jKey, wKey)
		wg.Done()
	}
}

func main() {
	dispatcher := buruh.New(&buruh.Config{
		MaxWorkerNum: 1000,
		MinWorkerNum: 100,
		Debug:        false,
	})

	numOfJob := 5000
	wg := sync.WaitGroup{}
	wg.Add(numOfJob)

	for i := 1; i <= numOfJob; i++ {
		job := buruh.NewJob(fn(i, &wg), false)
		dispatcher.Dispatch(job)
	}

	wg.Wait()

	dispatcher.Stop()
}
