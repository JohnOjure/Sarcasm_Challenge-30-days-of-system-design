import pytest
import time
from fastapi.testclient import TestClient
from unittest.mock import patch

from ..backend.bridge import app, rate_limit_store as limiter
from ..backend.classes import Bucket as TokenBucket

client = TestClient(app)

#unit tests
def test_token_bucket_initialization():
    #verify a bucket starts full
    bucket = TokenBucket(capacity=10, refill_rate=1.0)
    assert bucket.available_tokens == 10

def test_token_bucket_consumption():
    #verify consuming tokens reduces the count
    bucket = TokenBucket(capacity=10, refill_rate=1.0)
    success = bucket.consume(1)
    assert success is True
    assert bucket.available_tokens == 9

def test_token_bucket_exhaustion():
    #verify we cannot consume more than we have
    bucket = TokenBucket(capacity=1, refill_rate=1.0)
    bucket.consume(1) 
    success = bucket.consume(1)
    assert success is False

def test_token_bucket_refill():
    #verify the bucket refills over time

    with patch('time.monotonic') as mock_time:
        mock_time.side_effect = [0.0, 0.0, 0.5] #simuate time progression
        
        #time is 0.0
        bucket = TokenBucket(capacity=10, refill_rate=10.0)
        
        #time is 0.0 still
        bucket.consume(10)
        assert bucket.available_tokens == 0
    
        # time.sleep(0.5)
        #time is now 0.5
        bucket.consume(0) #trigger lazy refill check
    
        # #should have like about 5 tokens now
        # assert 4.0 < bucket.tokens < 6.0
        assert bucket.available_tokens == 5.0


#integration tests
def test_rate_limit_middleware_enforcement():
    #test full flow

    limiter._buckets.clear()
    
    client_ip = "test_client"
    
    #simulate 10 successful requests
    for _ in range(10):
        response = client.get("/")
        assert response.status_code == 200

    response = client.get("/")
    assert response.status_code == 429
    # assert response.json() == {"error": "Rate limit exceeded", "retry_after": "Wait for token refill"}

def test_dynamic_rule_update():
    #test that rules can be changed dynamically

    limiter._buckets.clear()
    
    for _ in range(10):
        client.get("/")
    
    assert client.get("/").status_code == 429
    
    #give this ip a bigger bucket
    update_payload = {
        "client_id": "testclient", 
        "capacity": 50,
        "refill_rate": 10.0
    }
    response = client.post("/admin/update_rule", json=update_payload)
    assert response.status_code == 200

    time.sleep(0.1)
    
    response = client.get("/")
    assert response.status_code == 200