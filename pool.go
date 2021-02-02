package buruh

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

type Pool struct {
	config        *Config
	workers       []*Worker
	job           chan Job
	runningWorker int64
	mu            sync.Mutex
}

func NewPool(cfg *Config) *Pool {
	p := &Pool{
		config:        cfg,
		workers:       []*Worker{},
		job:           make(chan Job),
		runningWorker: 0,
		mu:            sync.Mutex{},
	}

	return p
}

func (p *Pool) Init() {
	for i := 0; i < int(p.config.MinWorkerNum); i++ {
		if p.config.Debug {
			log.Println("Init new worker")
		}

		p.addNewWorker()
	}
}

func (p *Pool) submit(job Job) {
	if p.runningWorker == int64(len(p.workers)) {
		p.addNewWorker()
	}

	if p.config.Debug {
		log.Println("Waiting for available worker")
	}

	p.job <- job
}

func (p *Pool) tracker(state string) {
	if state == workerStartJob {
		atomic.AddInt64(&p.runningWorker, 1)
	} else {
		atomic.AddInt64(&p.runningWorker, -1)
	}
}

func (p *Pool) addNewWorker() {
	if uint64(len(p.workers)) < p.config.MaxWorkerNum {
		if p.config.Debug {
			log.Println("Add new worker")
		}

		w := NewWorker(p.config, p.tracker)
		p.workers = append(p.workers, w)
		go w.Start(p.job)
	}
}

func (p *Pool) Stop() {
	wg := &sync.WaitGroup{}
	wg.Add(len(p.workers))

	fmt.Println("Total workers: ", len(p.workers))

	for _, w := range p.workers {
		w.Stop(wg)
	}

	wg.Wait()
}
