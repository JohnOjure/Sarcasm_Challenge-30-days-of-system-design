import time
import os
import fastapi
from fastapi import FastAPI, Request, HTTPException
from fastapi import FastAPI, Request, HTTPException, status
from fastapi.responses import JSONResponse

app = FastAPI()

@app.middleware("http")
async def rate_limit_middleware(request: Request)