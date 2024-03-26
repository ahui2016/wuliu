import sys
import json
import shutil
import msgpack
import argparse
from pathlib import Path
from typing import Tuple
from operator import itemgetter
from wuliu.const import *
from wuliu.common import check_filename, print_err_exit


def create_new_album_info(filename: str):
    target_path = Path(filename)
    if target_path.exists():
        print(f'Error! File Exists: {filename}')
        return

    print(f'Create => {filename}')
    blob = json.dumps(New_Album_Info, indent=4)
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
    return info, err


def get_pics_metadata() -> list:
    """獲取全部圖片的屬性
    
    msgp_path 是由 "wuliu-db -dump" 導出的數據，檔名固定為 pics.msgp
    :return: 返回 File 列表 (參考 util/model.go 裏的 File)
    """
    with Path(Pics_msgp).open(mode='rb') as f:
        return msgpack.load(f)


def get_orderby(x: str) -> str:
    """
    :return: 默認返回 utime
    """
    match x:
        case 'ctime':
            return 'CTime'
        case 'filename':
            return 'Filename'
        case 'like':
            return 'Like'
        case _:
            return 'UTime'


def get_pics(info: dict) -> list:
    pics = get_pics_metadata()
    orderby = get_orderby(info['orderby'])
    reverse = not info['ascending']
    pics.sort(key=itemgetter(orderby), reverse=reverse)
    return pics


def create_album(pics: list, album_path: Path):
    pics_path = album_path.joinpath("pics")  # 原圖資料夾
    print(f'Create => {pics_path}')
    pics_path.mkdir(parents=True)
    print(f'Copy pictures from files to {pics_path}')
    for file in pics:
        src = Path("files").joinpath(file['Filename'])
        dst = pics_path.joinpath(file['ID']+src.suffix)
        print('.', end='')
        shutil.copyfile(src, dst)
    print('\nOK')


def make_album(pics: list, info: dict):
    """新建或更新相簿。
    """
    album_path = Path(Webpages).joinpath(info['name'])
    if album_path.exists():
        print("update album")
    else:
        create_album(pics, album_path)


# ↓↓↓ main ↓↓↓ 

parser = argparse.ArgumentParser()

parser.add_argument('-json', type=str, help='read album info')

parser.add_argument('--new-json', type=str,
        help='a filename of the json file of album info')

args = parser.parse_args()

if args.new_json:
    create_new_album_info(args.new_json)
    sys.exit()

if args.json:
    info, err = read_album_info(args.json)
    print_err_exit(err, front_msg=f'{args.json}, name: {info['name']}')
    pics = get_pics(info)
    make_album(pics, info)
