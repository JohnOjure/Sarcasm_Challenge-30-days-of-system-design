package main

import (
	"fmt"
	"time"
)


func main() {
	config := PoolConfig{
			MinWorkers:  2,
			MaxWorkers:  5,
			QueueSize:   10,
			IdleTimeout: 2 * time.Second,
	}

	pool := NewWorkerPool(config)
	pool.Start()

	//simulate a process consuming the results of submitted tasks
	go func() {
		for res := range pool.results {
			fmt.Printf("Result: [Task %d] %s\n", res.TaskID, res.Output)
		}
	}()

	//simulate processes submitting tsaks to the pool
	fmt.Println("About to submit a huge amount of tasks")
	for i := 1; i <= 20; i++ {
		pool.Submit(Task{ID: i, Payload: "Some type stuff", CreatedAt: time.Now()})
		time.Sleep(50 * time.Millisecond)
	} 

	time.Sleep(5 * time.Second) //scaling should occur

	pool.Stop()
}