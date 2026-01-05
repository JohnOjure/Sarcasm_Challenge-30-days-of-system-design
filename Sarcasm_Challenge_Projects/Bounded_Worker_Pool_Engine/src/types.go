package main

import "time"

//unit of work 
type Task struct {
	ID int
	Payload string
	CreatedAt time.Time
}

//result of a task
type Result struct {
	TaskID int
	Output string
	ProcessedAt time.Time
	Error error
}

//configuration for the worker pool
type PoolConfig struct {
	MinWorkers int
	MaxWorkers int
	QueueSize int
	IdleTimeout time.Duration
}
