from typing import List, Optional
import yaml
from pydantic import BaseModel


class Config(BaseModel):
    version: str
    mods: List[str]


    def __init__(self):
        if