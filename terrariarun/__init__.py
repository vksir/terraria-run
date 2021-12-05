import logging
from importlib import import_module
from fastapi import FastAPI
from terrariarun.common.log import log, init_log
from terrariarun.common.utils import init_path


init_path()
init_log()
log.setLevel(logging.DEBUG)
app = FastAPI(title='Terraria Run', version='0.1.0')
import_module('terrariarun.routes.routes')

