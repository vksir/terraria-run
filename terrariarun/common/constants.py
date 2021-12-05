import os
import platform


USER_HOME = os.environ['HOME'] if platform.system() == 'Linux' else os.environ['USERPROFILE']
HOME = f'{USER_HOME}/.terraria_run'
CFG_PATH = f'{HOME}/cfg.yaml'
LOG_PATH = f'{HOME}//terraria_run.log'


SERVER_DIR = f'{USER_HOME}/terraria'
SERVER_PATH = f'{SERVER_DIR}/tModLoaderServer'
SERVER_LOG_PATH = f'{HOME}/terraria.log'

WORLD_DIR = f'{USER_HOME}/.local/share/Terraria/ModLoader/Worlds'
WORLD_BACKUP_DIR = f'{USER_HOME}/.local/share/Terraria/ModLoader/Worlds/backup'
WORLD_NAME = 'Aurora'
WORLD_FILE_PATHS = [f'{WORLD_DIR}/{WORLD_NAME}.{file_type}'
                    for file_type in ('wld', 'wld.bak', 'twld', 'twld.bak')]
MOD_DIR = f'{USER_HOME}/.local/share/Terraria/ModLoader/Mods'
MOD_CFG_PATH = f'{MOD_DIR}/enabled.json'
