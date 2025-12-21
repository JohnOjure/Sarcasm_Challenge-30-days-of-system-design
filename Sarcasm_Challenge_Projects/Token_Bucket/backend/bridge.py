import time
import logging
from fastapi import FastAPI, Request, HTTPException
from fastapi import FastAPI, Request, HTTPException, status
from fastapi.responses import JSONResponse
import uvicorn

from .classes import RateLimitStore
from .schemas import RuleUpdate

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger("RateLimiterStore")

app = FastAPI()
rate_limit_store = RateLimitStore()

@app.middleware("http")
async def rate_limit_middleware(request: Request, call_next):
    client_id = request.client.host #the ip address of the client that made the request

    #exclude admin route from rate limiting
    if request.url.path == "/admin/update_rule":
        return await call_next(request)

    #check if the request is allowed
    if rate_limit_store.allow_request(client_id):
        logger.info(f"Allowed request from {client_id}")

        response = await call_next(request)
        return response
    else:
        logger.warning(f"Rejected request from {client_id}, rate limit exceeded")

        return JSONResponse(
            status_code=status.HTTP_429_TOO_MANY_REQUESTS,
            content={"detail": "Oh no..., rate limit exceeded. Pele, try again later."}
        )
    

@app.get("/")
async def root():
    return {"message": "All services activeee like Chivita. Request accepted", "timestamp": time.time()}

@app.post("/admin/update_rule")
async def update_rule(rule_update: RuleUpdate):
    rate_limit_store.update_rule(
        client_id=rule_update.client_id,
        capacity=rule_update.capacity,
        refill_rate=rule_update.refill_rate
    )

    return JSONResponse(
        status_code=status.HTTP_200_OK,
        content={"detail": f"Updated rate limit for {rule_update.client_id}"}
    )


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)