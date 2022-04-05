import json
import re
import time
from typing import Union
from threading import Thread

import requests

from terrariarun import log
from terrariarun.common.constants import *
from terrariarun.common.utils import run_cmd, ProcessHandler
from terrariarun.controller.config import Config


LOG_HANDLE_INTERVAL = 3


class LogHandler(Thread):
    def __init__(self, proc_handler: ProcessHandler):
        super().__init__()
        self._proc_handler = proc_handler

    def run(self) -> None:
        while True:
            out = self._proc_handler.read()
            self._handle(out)
            time.sleep(LOG_HANDLE_INTERVAL)

    def _handle(self, out: str):
        print(out)


class Controller:
    __instance = None
    _proc_handler: Union[ProcessHandler, None] = None
    _log_handler: Union[LogHandler, None] = None
    _cfg: Config = None

    def __new__(cls):
        if cls.__instance is None:
            cls.__instance = super().__new__(cls)
            cls.__instance._init()
        return cls.__instance

    def _init(self):
        self._cfg = Config()

    def start(self):
        if self._proc_handler or self._log_handler:
            return
        self.update(force=False)
        self.create_world(force=False)
        self.update_mods(force=False)

        self._proc_handler = self._run()
        self._log_handler = LogHandler(self._proc_handler)
        self._log_handler.start()

    def stop(self):
        if not self._proc_handler or not self._log_handler:
            return
        self._proc_handler.exit()
        self._proc_handler, self._log_handler = None, None

    def restart(self):
        self.stop()
        self.start()

    def update(self, force=True):
        cur_version = self._get_latest_version()
        if not force and self._cfg.version == cur_version:
            return
        self._cfg.version = cur_version
        self._cfg.save()

        filename = f'tModLoader.Linux.{self._cfg.version}.tar.gz'
        url = f'{MOD_DOWNLOAD_URL}/{filename}'
        for cmd in [
            f'curl -OL {url} -x {PROXY_URL}',
            f'tar -xzvf {filename}',
            f'rm {filename}'
        ]:
            run_cmd(cmd, cwd=SERVER_DIR)

    @staticmethod
    def update_mods(force=True):
        if not os.path.exists(f'{MOD_CFG_PATH}'):
            with open(f'{MOD_CFG_PATH}', 'w', encoding='utf-8') as f:
                json.dump([], f)
        with open(f'{MOD_CFG_PATH}', 'r', encoding='utf-8') as f:
            mods = json.load(f)

        for mod in mods:
            mod_path = f'{MOD_DIR}/{mod}.tmod'
            if not force and os.path.exists(mod_path):
                continue
            log.info(f'begin download mod: mod={mod}')
            run_cmd(f'curl -OL {MOD_DOWNLOAD_URL}/{mod}.tmod', cwd=MOD_DIR)

    def create_world(self, force=True):
        """

        Choose World:

        1	Small
        2	Medium
        3	Large
        Choose size:

        1	Normal
        2	Expert
        Choose difficulty:
        Enter world name:
        """
        if not force and self._is_world_exists:
            return

        log.info('begin create_world')
        proc_handler = ProcessHandler(SERVER_PATH, SERVER_LOG_PATH)
        for cmd in [
            'n',
            '3',
            '2',
            WORLD_NAME,
        ]:
            proc_handler.run_cmd(cmd)

    def _run(self):
        """

        Choose World:
        Max players (press enter for 8):
        Server port (press enter for 7777):
        Automatically forward port? (y/n):
        Server password (press enter for none):
        """
        proc_handler = ProcessHandler(SERVER_PATH, SERVER_LOG_PATH)
        for cmd in [
            '1',
            '',
            '',
            '',
            self._cfg.room_passwd
        ]:
            proc_handler.run_cmd(cmd)
        return proc_handler

    @property
    def _is_world_exists(self):
        return all(os.path.exists(path) for path in WORLD_FILE_PATHS)

    @staticmethod
    def _backup_world():
        log.info('begin backup_world')
        file_name = time.strftime(f'{WORLD_NAME}_%Y-%m-%d_%H-%M-%S', time.localtime())
        run_cmd(f'tar -cvzf {WORLD_BACKUP_DIR}/{file_name}.tar.gz *', cwd=WORLD_DIR)
        for path in WORLD_FILE_PATHS:
            run_cmd(f'rm {path}')

    @staticmethod
    def _get_latest_version():
        proxies = {
            'http': PROXY_URL,
            'https': PROXY_URL
        }
        resp = requests.get(SERVER_DOWNLOAD_URL, proxies=proxies)
        version = resp.url.split('/')[-1]
        return version


if __name__ == '__main__':
    ctr = Controller()
    # agent.create_world()
    ctr.start()
