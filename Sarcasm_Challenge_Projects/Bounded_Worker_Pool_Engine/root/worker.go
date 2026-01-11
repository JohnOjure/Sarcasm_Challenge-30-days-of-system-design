package main

import (

	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

func (p *WorkerPool) spawnWorker(id int) {
	atomic.AddInt32(&p.active, 1)
	p.wg.Add(1)

	go func(workerID int) {
		defer p.wg.Done()
		defer atomic.AddInt32(&p.active, -1)

		fmt.Printf("Worker %d started", workerID)
		timer := time.NewTimer(p.config.IdleTimeout)
		defer timer.Stop()

		for {
			select{
			case task, ok := <-p.tasks:
				if !ok { //if the channel has closed
					return
				}

				if !timer.Stop() { //reset idle timer
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(p.config.IdleTimeout)

				res := p.execute(workerID, task) //process the task
				p.results <- res
			case <-timer.C:
				//timeout, this worker hasnt received 
				//any tasks so spin it down if there are temporary workers

				current := atomic.LoadInt32(&p.active) //we use atomic warrever warrever when we are dealing with metric 
				// values that multiple goroutines are interacting with concurrently

				if int(current) > p.config.MinWorkers {
					fmt.Printf("So like, worker %d has been idle for too long so it's shutting down", workerID)
					return
				}
				timer.Reset(p.config.IdleTimeout)
			}
		}
	} (id) //an anonympus function

}

func (p *WorkerPool) execute(workerID int, t Task) Result {
	sleepTime := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(sleepTime)

	return Result{
		TaskID: t.ID,
		Output: fmt.Sprintf("Processed by Worker %d in %v", workerID, sleepTime),
	}
}
