# Bounded Worker Pool Engine (Go)

A high-performance, concurrent worker pool implementation in Go. This project demonstrates advanced concurrency patterns including **channel-based backpressure**, **dynamic worker auto-scaling**, and **graceful shutdown** mechanisms.

## 🚀 Features

* **Bounded Task Queue:** Prevents memory overload by blocking task submission when the queue is full (Backpressure).
* **Dynamic Auto-Scaling:** Automatically spawns new workers (up to a limit) when the queue load exceeds 50%.
* **Scale Down:** Idle workers automatically shut themselves down to save resources.
* **Graceful Shutdown:** Ensures all active tasks are completed before the application exits.
* **Race Condition Free:** Thread-safe operations using `sync/atomic` and `sync.WaitGroup`.

## 📂 Project Structure

```text
├── main.go        # Entry point: Configuration and simulation of traffic
├── pool.go        # The Engine: Manages the queue, scaler monitor, and lifecycle
├── worker.go      # The Labor: Individual worker logic and self-destruction
├── types.go       # Data Models: Task, Result, and PoolConfig definitions
├── pool_test.go   # Unit Tests: Verifies execution and config validation
├── Dockerfile     # Multi-stage build for containerization
└── go.mod         # Go module definition
```

## ⚙️ Architecture

The system uses a Fan-Out pattern with a dynamic scaler monitor.
```
graph TD
    P[Producer/Main] -->|Submits Task| Q(Buffered Channel)
    
    subgraph Worker Pool
        Q --> W1[Worker 1]
        Q --> W2[Worker 2]
        Q -.-> W3[Worker 3 (Dynamic)]
        
        M[Scaler Monitor] -.->|Checks Depth| Q
        M -.->|Spawns| W3
    end
    
    W1 --> R(Results Channel)
    W2 --> R
    W3 --> R
    
    R --> C[Consumer]
```

## 🛠️ Installation & Usage

### Prerequisites
* Go 1.21+
* Docker (Optional)

### 1.Run Locally
Initialize module (first time only):

```
go mod init worker-pool
go mod tidy
```

Run the application:

```
go run .
```

### 2. Run with Docker
Build the lightweight Alpine image:

```
docker build -t worker-pool-app .
```

Run the container:

```
docker run --rm worker-pool-app
```

### 3. Running Tests
Run unit tests with race condition detection:

```
go test -v -race
```

## 🎛️ Configuration
You can tune the pool `behavior` in main.go via the `PoolConfig` struct



