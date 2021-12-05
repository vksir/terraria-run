import os
import shlex
import subprocess
import time

from terrariarun import log
from terrariarun.common import constants


def run_cmd(cmd: str, cwd=None, sudo=False) -> str:
    if sudo:
        cmd += 'sudo '

    p = subprocess.run(cmd, cwd=cwd, shell=True,
                       encoding='utf-8',
                       stderr=subprocess.STDOUT,
                       stdout=subprocess.PIPE)
    if p.returncode:
        log.error(f'run_cmd failed: cmd={cmd}, out={p.stdout}')
    return p.stdout


class ProcessHandler:
    def __init__(self, cmd: str, log_path: str, cwd=None):
        self._log_writer = open(log_path, 'w', encoding='utf-8')
        self._log_reader = open(log_path, 'r', encoding='utf-8')
        self._cmd_reader = open(log_path, 'r', encoding='utf-8')

        self._cmd = shlex.split(cmd)
        log.info(f'start process: cmd={self._cmd}')
        self._proc = subprocess.Popen(self._cmd, cwd=cwd, shell=False,
                                      encoding='utf-8', bufsize=1,
                                      stderr=subprocess.STDOUT,
                                      stdout=self._log_writer,
                                      stdin=subprocess.PIPE)

    def __del__(self):
        self._log_writer.close()
        self._log_reader.close()
        self._cmd_reader.close()

    def run_cmd(self, cmd: str) -> str:
        if not cmd.endswith('\n'):
            cmd = cmd + '\n'
        self._cmd_reader.read()
        self._proc.stdin.write(cmd)
        out = self._cmd_reader.read()
        return out

    def read(self) -> str:
        return self._log_reader.read()

    def exit(self, timeout=15):
        self._proc.send_signal(2)
        cost_time = 0
        while self._proc.poll() is None and cost_time < timeout:
            time.sleep(1)
            cost_time += 1
        if cost_time >= timeout:
            log.error('process exit timeout, begin to kill')
            self._proc.kill()
        out = run_cmd(f'ps -ef | grep -v grep | grep {self._cmd[0]}')
        if out:
            log.error(f'process exit failed: cmd={self._cmd}, out={out}')
            return
        log.info(f'process exit success: cmd={self._cmd}')


def init_path():
    for path_name in dir(constants):
        path = getattr(constants, path_name)
        if (not path_name.endswith('HOME') and not path_name.endswith('DIR')) \
                or os.path.exists(path):
            continue
        run_cmd(f'mkdir -p {path}')
