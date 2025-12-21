from pydantic import BaseModel

class RuleUpdate(BaseModel):
    client_id: str
    capacity: int
    refill_rate: float


    