package buruh

import (
	"context"

	"github.com/rs/xid"
)

type Task func(ctx context.Context)

type Job struct {
	ID         xid.ID
	task       Task
	finishChan chan bool
}

func NewJob(task Task, withID bool) (job Job) {

	// job = Job{
	// 	task: task,
	// }

	// job.finishChan = make(chan bool, 2)
	job.task = task

	if withID {
		job.ID = xid.New()
	}

	return
}

func (j *Job) Do(ctx context.Context) {
	j.task(ctx)
}

func (j *Job) finish() {
	j.finishChan <- true
}

func (j *Job) Wait() {
	<-j.finishChan
}
