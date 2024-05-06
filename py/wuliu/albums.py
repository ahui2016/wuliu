import sys
import json
import shutil
import msgpack
from pathlib import Path
from typing import Tuple, Set
from .const import *
from .common import check_filename, print_err_exit


"""
這是 wuliu-photo-album.py 與 wuliu-docs-preview.py 的共用函數。
"""


def create_new_album_info(filename: str):
    target_path = Path(filename)
    if target_path.exists():
        print(f'Error! File Exists: {filename}')
        return

    print(f'Create => {filename}')
    blob = json.dumps(New_Album_Info, ensure_ascii=False, indent=4)
    target_path.write_text(blob, encoding='utf8')


def read_album_info(filename: str):
    """
    :return: Tuple[dict|None, str]
    """
    data = Path(filename).read_text(encoding='utf8')
    info = json.loads(data)

    if info['name'] == '':
        return None, f'請在 {filename} 中填寫 name, 不可留空。'
    err = check_filename(info['name'])

    if not info['union']:
        return None, 'union 請設為 true, 取交集的功能暫未實現。'

    return info, err


def get_files_metadata(msgp_name: str) -> list:
    """獲取指定 msgp 中的全部檔案屬性。
    
    msgp_name 是由 "wuliu-db -dump" 導出的數據檔案名稱。
    :return: 返回 File 列表 (參考 util/model.go 裏的 File)
    """
    if not Path(msgp_name).exists():
        print(f'ERROR: 找不到 {msgp_name}')
        print('請執行 "wuliu-db -dump" 導出 msgp 檔案 (參考README.md)')
        sys.exit(1)
    with Path(msgp_name).open(mode='rb') as f:
        return msgpack.load(f)


def read_album_msgp(album_path: Path, msgp_name: str) -> dict:
    """讀取 album_path 中的 msgp 檔案 (舊的檔案屬性)。"""
    msgp_path = album_path.joinpath(msgp_name)
    if not msgp_path.exists():
        print_err_exit(f'ERROR: No such file: {msgp_path}')
    data = msgp_path.read_bytes()
    return msgpack.unpackb(data)


def write_album_msgp(
        files: dict, album_info: dict, album_path: Path, msgp_name: str, tmpl_name: str):
    if msgp_name == Pics_msgp:
        msgp_path = album_path.joinpath(msgp_name)
        blob = msgpack.packb(files)
        print(f'Write => {msgp_path}')
        msgp_path.write_bytes(blob)

    album_info['files'] = list(files.values())
    blob = json.dumps(album_info, ensure_ascii=False, indent=4)
    blob = 'const albumData = ' + blob;
    js_path = album_path.joinpath('files.js')
    print(f'Write => {js_path}')
    js_path.write_text(blob, encoding='utf8')
    
    src = Path(Webpages).joinpath('templates', tmpl_name)
    dst = album_path.joinpath('index.html')
    if not dst.exists():
        print(f'Write => {dst}')
        shutil.copyfile(src, dst)


def keywords_union(files: list, album_info: dict) -> set:
    result: Set[str] = set()
    for x in album_info['keywords']:
        good = {f[ID] for f in files if x in f[Keywords]}
        result = result.union(good)
    return result


def collections_union(files: list, album_info: dict) -> set:
    result: Set[str] = set()
    for x in album_info['collections']:
        good = {f[ID] for f in files if x in f[Collections]}
        result = result.union(good)
    return result


def albums_union(files: list, album_info: dict) -> set:
    result: Set[str] = set()
    for x in album_info['albums']:
        good = {f[ID] for f in files if x in f[Albums]}
        result = result.union(good)
    return result


def filter_files(files: list, album_info: dict) -> list:
    ids: Set[str] = set()
    
    if album_info['label'] + album_info['notes'] == '' and \
            len(album_info['keywords'])+len(album_info['collections'])+len(album_info['albums']) == 0:
        return files
    
    if album_info['label'] != '':
        by_label = {f[ID] for f in files if album_info['label'] == f[Label]}
        ids = ids.union(by_label)

    if album_info['notes'] != '':
        by_notes = {f[ID] for f in files if album_info['notes'] == f[Notes]}
        ids = ids.union(by_notes)
    
    by_keywords = keywords_union(files, album_info)
    ids = ids.union(by_keywords)

    by_collections = collections_union(files, album_info)
    ids = ids.union(by_collections)

    by_albums = albums_union(files, album_info)
    ids = ids.union(by_albums)
    
    result = []
    for f in files:
        if f[ID] in ids:
            result.append(f)

    return result
