import time
import threading

class Bucket:
    def __init__ (self, capacity: int, refill_rate: float):
        self.capacity = capacity
        self.refill_rate = refill_rate #number of tokens added per second
        self.available_tokens = capacity 
        self.last_refill = time.monotonic()

    
    def consume(self, tokens: int = 1):
        now = time.monotonic() #get the current time

        #refill tokens based on the elapsed time first
        elapsed = now - self.last_refill
        new_tokens = elapsed * self.refill_rate

        #update the tokens with the new tokens
        if new_tokens > 0:
            self.available_tokens = min(self.capacity, self.available_tokens + new_tokens)
            self.last_refill = now

        #now try to consume the requested tokens
        if self.available_tokens >= tokens:
            self.available_tokens -= tokens
            return True #allowed
        return False #rejected
    

class RateLimitStore:
    def __init__(self):
        self._buckets = {}
        self._lock = threading.Lock()

    def get_bucket(self, client_id: str):
        #if this is a new user, create a new bucket for them
        if client_id not in self._buckets:
            with self._lock:
                if client_id not in self._buckets:
                    self._buckets[client_id] = Bucket(capacity=10, refill_rate=1.0)
        return self._buckets[client_id]

    def update_rule(self, client_id: str, capacity: int, refill_rate: float):
        if client_id in self._buckets:
            with self._lock:
                bucket = self._buckets[client_id]
                bucket.capacity = capacity
                bucket.refill_rate = refill_rate
            print(f"Updated rate limit for {client_id}: capacity={capacity}, refill_rate={refill_rate}")

    def allow_request(self, client_id: str):
        bucket = self.get_bucket(client_id)

        with self._lock:
            return bucket.consume(1)