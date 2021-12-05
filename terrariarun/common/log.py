import logging
from terrariarun.common.constants import LOG_PATH


log = logging.getLogger(__name__)
formatter = logging.Formatter('[%(asctime)s] [%(levelname)s] [%(filename)s] %(message)s')


def init_log():
    sh = logging.StreamHandler()
    sh.setFormatter(formatter)
    fh = logging.FileHandler(LOG_PATH, 'w', encoding='utf-8')
    fh.setFormatter(formatter)
    log.addHandler(sh)
    log.addHandler(fh)





