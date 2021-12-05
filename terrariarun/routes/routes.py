from terrariarun import app
from terrariarun.common.utils import run_cmd


@app.on_event('startup')
def startup():
    pass


@app.on_event('shutdown')
def shutdown():
    pass


@app.post('/install')
def install():
    pass
