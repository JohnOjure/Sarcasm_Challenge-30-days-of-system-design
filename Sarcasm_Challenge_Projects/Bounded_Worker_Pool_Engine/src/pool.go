package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)


type WorkerPool struct {
	tasks chan Task
	results chan Result
	config PoolConfig
	active int32
	shutdown chan struct{}
	wg sync.WaitGroup
	quitMonitor chan struct{}
} 

func NewWorkerPool(cfg PoolConfig) *WorkerPool {
	return &WorkerPool{
		tasks: make(chan Task, cfg.QueueSize),
		results: make(chan Result, cfg.QueueSize),
		config: cfg,
		shutdown: make(chan struct{}),
		quitMonitor: make(chan struct{}),
	}
}

//method to start worker pool
func (p *WorkerPool) Start() {
	fmt.Printf("Aiit, worker pool starting with %d workers...\n", p.config.MinWorkers)

	//start minimum workers
	for i:= 0; i < p.config.MinWorkers; i++ {
		p.spawnWorker(i + 1)
	}

	//start dynamic scaler (monitor routie)
	go p.monitor()
}

//method to submit a task to a pool
func (p *WorkerPool) Submit (t Task) {
	p.tasks <- t
}

//method to stop pool gracefully
func (p *WorkerPool) Stop() {
	fmt.Println("Shutting down pool...")

	close(p.quitMonitor) 
	close(p.tasks) //even though it's closed, the tasks left in it will still be processed
	p.wg.Wait() //wait for all workers to finish
	close(p.results)

	fmt.Println("Pool stopped, arriverderci")
}


func (p *WorkerPool) monitor() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {

		select {
			case <-p.quitMonitor:
				return

			case <-ticker.C:
				queueLen := len(p.tasks)
				currentWorkers := atomic.LoadInt32(&p.active)

				// SCALE UP logic
				if queueLen > p.config.QueueSize/2 && int(currentWorkers) < p.config.MaxWorkers {
					fmt.Printf("[Scaler] Queue loaded (%d/%d). Spawning worker.\n", queueLen, p.config.QueueSize)
					p.spawnWorker(int(currentWorkers) + 1)
				}
		}
	}
}