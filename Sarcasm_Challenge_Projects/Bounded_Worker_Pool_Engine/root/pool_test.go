package main

import (
	"testing"
	"time"
)

//TestPoolExecution verifies that all submitted tasks are processed
func TestPoolExecution(t *testing.T) {
	cfg := PoolConfig{
		MinWorkers:  2,
		MaxWorkers:  4,
		QueueSize:   10,
		IdleTimeout: 1 * time.Second,
	}
	pool := NewWorkerPool(cfg)
	pool.Start()

	taskCount := 5
	for i := 0; i < taskCount; i++ {
		pool.Submit(Task{ID: i, Payload: "Test Payload", CreatedAt: time.Now()})
	}

	
	receivedCount := 0
	timeout := time.After(3 * time.Second) // Safety valve so test doesn't hang forever

	for i := 0; i < taskCount; i++ {
		select {
		case res := <-pool.results:
			if res.Error != nil {
				t.Errorf("Task %d failed: %v", res.TaskID, res.Error)
			}
			receivedCount++
		case <-timeout:
			t.Fatalf("Test timed out! Only received %d/%d results", receivedCount, taskCount)
		}
	}

	pool.Stop()

	if receivedCount != taskCount {
		t.Errorf("Expected %d results, got %d", taskCount, receivedCount)
	}
}


func TestPoolConfig(t *testing.T) {
	tests := []struct { //table of test cAses
		name        string
		inputConfig PoolConfig
		expectPanic bool
	}{
		{
			name: "Valid Config",
			inputConfig: PoolConfig{
				MinWorkers: 1, MaxWorkers: 5, QueueSize: 10, IdleTimeout: time.Second,
			},
			expectPanic: false,
		},
	}

	//loop over the table
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := NewWorkerPool(tc.inputConfig)
			if p == nil {
				t.Error("WorkerPool should not be nil")
			}
		})
	}
}