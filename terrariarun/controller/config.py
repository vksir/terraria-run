import os
from enum import Enum
from typing import List, Optional
import yaml
from pydantic import BaseModel
from terrariarun.common.constants import *


class ServerType(Enum):
    TERRARIA = 'Terraria'
    T_MOD_LOADER = 'tModLoader'


class Config(BaseModel):
    server: ServerType = ServerType.TERRARIA
    version: str = ''
    port: str = '7777'
    room_name: str = 'Terraria Run'
    room_passwd: str = '6666'
    mods: List[str] = []

    def __init__(self):
        if not os.path.exists(CFG_PATH):
            super().__init__(**{})
            self.save()
        self.read()

    def read(self):
        with open(CFG_PATH, 'r', encoding='utf-8') as f:
            cfg = yaml.load(f, Loader=yaml.Loader)
            super().__init__(**cfg)

    def save(self):
        with open(CFG_PATH, 'w', encoding='utf-8') as f:
            yaml.dump(self.dict(), f)

    def dict(self, *args, **kwargs):
        res = super().dict(*args, **kwargs)
        for key, value in res.items():
            if isinstance(value, Enum):
                res[key] = value.value
        return res
