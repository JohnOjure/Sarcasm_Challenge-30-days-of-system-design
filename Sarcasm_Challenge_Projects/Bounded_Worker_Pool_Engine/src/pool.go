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

	//start dynamic scaler
	go p.monitor()
}

//method to submit a task to a pool
func (p *WorkerPool) Submit (t Task) {
	p.tasks <- t
}

