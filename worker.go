package buruh

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/rs/xid"
)

const (
	workerStartJob  = "start"
	workerFinishJob = "finish"
)

type CtxKey string

var (
	CtxWorkerKey CtxKey = "worker-key"
	CtxJobKey    CtxKey = "job-key"
)

type workerTracker func(state string)

type Worker struct {
	config     *Config
	ID         xid.ID
	startTime  time.Time
	stopSignal chan bool
	tracker    workerTracker
}

func NewWorker(cfg *Config, tracker workerTracker) *Worker {
	uid := xid.New()

	if cfg.Debug {
		log.Printf("Spawn new worker, id: %s", uid.String())
	}

	return &Worker{
		ID:         uid,
		config:     cfg,
		startTime:  time.Now(),
		stopSignal: make(chan bool),
		tracker:    tracker,
	}
}

func (w *Worker) Start(jobChan <-chan Job) {

	for {
		select {
		case job := <-jobChan:
			w.do(job)
		case <-w.stopSignal:
			if w.config.Debug {
				log.Printf("Terminating worker: %s", w.ID.String())
			}

			w.stopSignal <- true
			return
		default:
			continue
		}

	}

}

func (w *Worker) do(job Job) {
	w.tracker(workerStartJob)

	if w.config.Debug {
		log.Printf("Execute job: %s, with worker: %s", job.ID.String(), w.ID.String())
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxWorkerKey, w.ID.String())
	ctx = context.WithValue(ctx, CtxJobKey, job.ID.String())

	job.Do(ctx)
	// job.finish()

	if w.config.Debug {
		log.Printf("Finish job: %s, with worker: %s", job.ID.String(), w.ID.String())
	}

	w.tracker(workerFinishJob)
}

func (w *Worker) Stop(wg *sync.WaitGroup) {
	w.stopSignal <- true

	<-w.stopSignal
	wg.Done()
}
