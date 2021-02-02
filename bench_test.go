package buruh_test

import (
	"context"
	"testing"

	"github.com/raspiantoro/buruh"
)

var testFn = func(a int, b int, ch chan<- bool) buruh.Task {
	return func(ctx context.Context) {
		// time.Sleep(5 * time.Second)
		c := a + b
		d := (a + b) * c
		ctx = context.WithValue(ctx, "calc-result", c)
		ctx = context.WithValue(ctx, "calc-result2", d)
		ch <- true
	}
}

func BenchmarkWithoutBuruh(b *testing.B) {
	ctx := context.Background()

	ch := make(chan bool)

	for i := 0; i < b.N; i++ {
		// perform the operation we're analyzing
		fn := testFn(i, i+1, ch)
		job := buruh.NewJob(fn, false)
		go job.Do(ctx)
		<-ch
		// go fn(ctx)
	}
}

func BenchmarkWithBuruh(b *testing.B) {
	dispatcher := buruh.New(&buruh.Config{
		MaxWorkerNum: 100,
		MinWorkerNum: 20,
	})

	ch := make(chan bool)

	for i := 0; i < b.N; i++ {
		// perform the operation we're analyzing
		fn := testFn(i, i+1, ch)
		job := buruh.NewJob(fn, false)
		dispatcher.Dispatch(job)
		<-ch
	}
}
