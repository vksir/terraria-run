import datetime
import json
import os
import zipfile
from os.path import expanduser

GAME_STEAM_ID = '1281930'
MOD_DIR = r'E:\Program Files (x86)\Steam\steamapps\workshop\content\1281930'

HOME = expanduser('~')
MOD_ENABLED_PATH = os.path.join(HOME, r'Documents\My Games\Terraria\tModLoader\Mods\enabled.json')
MOD_COLLECT_DIR = os.path.join(HOME, r'Documents\My Games\Terraria\tModLoader\Mods\MyMods')
MOD_COLLECT_ZIP_PATH = os.path.join(HOME, r'Documents\My Games\Terraria\tModLoader\Mods\MyMods.zip')


def get_enabled_mods():
    with open(MOD_ENABLED_PATH, 'r', encoding='utf-8') as f:
        return json.load(f)


def try_get_mod_file(mod_dir):
    filenames = os.listdir(mod_dir)
    for filename in filenames:
        file_path = os.path.join(mod_dir, filename)
        if os.path.isfile(file_path) and filename.endswith('.tmod'):
            return file_path
    return None


def get_mod_path(mod_dir):
    mod_path = try_get_mod_file(mod_dir)
    if mod_path is not None:
        return mod_path
    date_dirname_lst = [dirname for dirname in os.listdir(mod_dir) if os.path.isdir(os.path.join(mod_dir, dirname))]
    date_dirname_lst = sorted(date_dirname_lst, key=lambda date: datetime.datetime.strptime(date, '%Y.%m').timestamp())
    newest_dirname = date_dirname_lst[-1]
    newest_mod_dir = os.path.join(mod_dir, newest_dirname)
    mod_path = os.path.join(newest_mod_dir, os.listdir(newest_mod_dir)[0])
    print(f'date_dirname_lst: {date_dirname_lst}, mod_path={mod_path}')
    return mod_path


def get_exist_mods():
    mods = {}
    mod_dirname_lst = os.listdir(MOD_DIR)
    print(f'mod_dirname: {mod_dirname_lst}')
    for mod_dirname in mod_dirname_lst:
        mod_dir = os.path.join(MOD_DIR, mod_dirname)
        mod_path = get_mod_path(mod_dir)

        filename = os.path.split(mod_path)[1]
        mod_name = filename.split('.')[0]
        mods[mod_name] = mod_path
    return mods


def main():
    enabled_mods = get_enabled_mods()
    exist_mods = get_exist_mods()
    print(f'enabled mods: {enabled_mods}')
    mods = {mod: exist_mods[mod] for mod in enabled_mods}
    print(json.dumps(mods, indent=4))

    with zipfile.ZipFile(MOD_COLLECT_ZIP_PATH, 'w') as myzip:
        for mod_path in mods.values():
            filename = os.path.split(mod_path)[1]
            myzip.write(mod_path, arcname=filename)


if __name__ == '__main__':
    main()
