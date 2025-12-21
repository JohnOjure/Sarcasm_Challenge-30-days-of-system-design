# API Gateway Rate Limiter (Token Bucket)

A high-performance, thread-safe microservice built with **FastAPI** that implements the **Token Bucket Algorithm** to rate-limit HTTP traffic.

This project demonstrates core backend engineering concepts: middleware architecture, concurrency control (thread locking), and precise time accounting for distributed systems.

## 🚀 Features

- **Token Bucket Algorithm:** Allows for bursty traffic while enforcing a long-term average rate.
- **Thread-Safe State:** Uses `threading.Lock` to handle concurrent requests without race conditions.
- **FastAPI Middleware:** Intercepts requests at the gateway level, rejecting traffic before it reaches business logic.
- **Dynamic Configuration:** Supports runtime updates to rate limits (Capacity/Refill Rate) via an Admin API without restarting the server.
- **Precise Timing:** Uses `time.monotonic()` to prevent drift from system clock updates (NTP).

## 🛠️ Tech Stack

- **Python 3.9+**
- **FastAPI:** Modern web framework for building APIs.
- **Uvicorn:** ASGI server for production.
- **Pytest:** Testing framework (includes `unittest.mock` for time simulation).

## 📦 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/rate-limiter.git
   cd rate-limiter
   ```

2. **Install dependencies**
   ```bash
   pip install fastapi uvicorn pytest httpx
   ```

## ⚡ Usage

1. **Run the Server**  
   The application runs on port 8000.
   ```bash
   python main.py
   ```
   > Note: Since state is in-memory, run with `workers=1` to ensure shared state (e.g., `uvicorn main:app --workers 1`).

2. **Test Rate Limiting**  
   Send requests to the root endpoint. By default: burst capacity of 10 tokens and refill rate of 1 token/sec.
   ```bash
   # Allowed (200 OK)
   curl -v http://localhost:8000/

   # After exceeding the limit quickly:
   # HTTP/1.1 429 Too Many Requests
   # {"error": "Rate limit exceeded", ...}
   ```

3. **Dynamic Rule Update**  
   Update rate limit rules for a specific client IP on the fly.
   ```bash
   curl -X POST http://localhost:8000/admin/update-rule \
     -H "Content-Type: application/json" \
     -d '{
       "client_id": "127.0.0.1",
       "capacity": 50,
       "refill_rate": 10.0
     }'
   ```

## 🧪 Testing

Comprehensive test suite using `pytest`:

- Unit tests: Validate token bucket logic with mocked time.
- Integration tests: Verify middleware correctly enforces limits.

Run tests:
```bash
pytest test_main.py -v
```

## 📐 Architecture

### Token Bucket Logic

Each client (identified by IP) has a bucket with:
- `capacity`: maximum tokens (burst size)
- `refill_rate`: tokens added per second

**Lazy Refill Strategy**:  
No background timer. On each request:
- Calculate new tokens: `(current_time - last_visit) * refill_rate`
- Add them to the bucket (capped at capacity)
- Update `last_visit` timestamp

**Consumption**:
- If current tokens ≥ 1 → consume 1 token, allow request
- Else → reject with 429 Too Many Requests

### Thread Safety

FastAPI/Uvicorn uses multiple workers/threads. To prevent race conditions when updating bucket state, all operations are protected by a per-client `threading.Lock`.

The critical section ensures atomicity for:
```
Read → Refill → Check → Consume → Update
```

This guarantees accurate rate limiting even under high concurrency.
