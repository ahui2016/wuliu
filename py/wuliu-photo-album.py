import sys
import json
import shutil
import msgpack
import argparse
from pathlib import Path
from typing import Tuple, Set
from operator import itemgetter
from wuliu.const import *
from wuliu.common import check_filename, print_err, print_err_exit, create_thumb, read_project_info


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


def read_pics_msgp(album_path: Path):
    pics_msgp_path = album_path.joinpath('pics.msgp')
    data = pics_msgp_path.read_bytes()
    return msgpack.unpackb(data)


def write_pics_msgp(pics: dict, album_path: Path):
    pics_msgp_path = album_path.joinpath('pics.msgp')
    blob = msgpack.packb(pics)
    print(f'Write => {pics_msgp_path}')
    pics_msgp_path.write_bytes(blob)


def create_album(pics: list, album_path: Path, thumb_size):
    pics_path = album_path.joinpath('pics')  # 原圖資料夾
    print(f'mkdir => {pics_path}')
    pics_path.mkdir(parents=True)
    thumbs_path = album_path.joinpath('thumbs')  # 縮略圖資料夾
    print(f'mkdir = {thumbs_path}')
    thumbs_path.mkdir(parents=True)
    print(f'Copy pictures from files to {pics_path}')
    album_pics = dict()

    for file in pics:
        file_id = file[ID]
        src = Path(Files).joinpath(file[Filename])
        dst = pics_path.joinpath(file_id+src.suffix)
        print('.', end='')
        shutil.copyfile(src, dst)
        thumb = thumbs_path.joinpath(file_id+'.jpg')
        err = create_thumb(src, thumb, thumb_size)
        print_err_exit(err)
        album_pics[file_id] = file

    print()
    write_pics_msgp(album_pics, album_path)
    print('OK')


def get_deleted_pics(pics:list, old_pics:dict) -> dict:
    """返回需要刪除的縮略圖的ID 和 filename"""
    ids: Set[str] = set()
    pics_ids = {pic[ID] for pic in pics}
    old_ids = {pic_id for pic_id in old_pics.keys()}
    deleted_ids = old_ids.difference(pics_ids)
    
    deleted_pics = dict()
    for pic_id in deleted_ids:
        deleted_pics[pic_id] = old_pics[pic_id][Filename]

    return deleted_pics


def delete_album_pics(deleted_pics:dict, album_pics:dict, album_path:Path) -> dict:
    """返回更新後的 album_pics"""
    if len(deleted_pics) == 0:
        return album_pics

    pics_dir = album_path.joinpath('pics')  # 原圖資料夾
    thumbs_dir = album_path.joinpath('thumbs')  # 縮略圖資料夾

    for pic_id, filename in deleted_pics.items():
        suffix = Path(filename).suffix
        pic_path = pics_dir.joinpath(f'{pic_id}{suffix}')
        thumb_path = thumbs_dir.joinpath(f'{pic_id}.jpg')
        print(f'Delete => [{pic_id}] {filename}')
        pic_path.unlink(missing_ok=True)
        thumb_path.unlink(missing_ok=True)
        album_pics.pop(pic_id, None)

    print()
    return album_pics


def pic_exists(pic:dict, old_pics:dict) -> bool:
    pic_id = pic[ID]
    if pic_id not in old_pics:
        return False

    old_pic = old_pics[pic_id]
    old_checksum = old_pic[Checksum]
    return old_checksum == pic[Checksum]


def get_updated_pics(pics:list, album_pics:dict) -> list:
    """
    :return: pics need to copy or overwrite
    """
    newpics = []
    for pic in pics:
        if not pic_exists(pic, album_pics):
            newpics.append(pic)
    return newpics


def update_album_pics(newpics:list, album_pics:dict, album_path: Path, thumb_size) -> dict:
    """返回更新後的 album_pics"""
    if len(newpics) == 0:
        return album_pics

    pics_dir = album_path.joinpath('pics')  # 原圖資料夾
    thumbs_dir = album_path.joinpath('thumbs')  # 縮略圖資料夾

    for pic in newpics:
        pic_id = pic[ID]
        src = Path(Files).joinpath(pic[Filename])
        dst = pics_dir.joinpath(pic_id+src.suffix)
        thumb = thumbs_dir.joinpath(pic_id+'.jpg')
        print(f'Add or update: [{pic_id}] {pic[Filename]}')
        shutil.copyfile(src, dst)
        err = create_thumb(src, thumb, thumb_size)
        print_err(err)
        if err is None:
            album_pics[pic_id] = pic

    return album_pics


def update_album_pics_msgp(pics:list, album_pics:dict, album_path:Path):
    for pic in pics:
        pic_id = pic[ID]
        if pic_id in album_pics:
            album_pics[pic_id] = pic

    # TODO: 排序在前端 js 做。
    write_pics_msgp(album_pics, album_path)


def update_album(pics:list, album_path:Path, thumb_size):
    old_pics = read_pics_msgp(album_path)
    deleted_pics = get_deleted_pics(pics, old_pics)
    album_pics = delete_album_pics(deleted_pics, old_pics, album_path)
    
    updated_pics = get_updated_pics(pics, album_pics)
    album_pics = update_album_pics(updated_pics, album_pics, album_path, thumb_size)

    if len(deleted_pics)+len(updated_pics) == 0:
        print('圖片無變化 (圖片無新增、更改或刪除)')
        return

    update_album_pics_msgp(pics, album_pics, album_path)


def make_album(pics: list, album_info: dict, proj_info: dict):
    """新建或更新相簿。
    """
    album_path = Path(Webpages).joinpath(album_info['name'])
    thumb_size = proj_info[Thumb_Size]

    if album_path.exists():
        update_album(pics, album_path, thumb_size)
    else:
        create_album(pics, album_path, thumb_size)


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
    proj_info = read_project_info()
    album_info, err = read_album_info(args.json)
    print_err_exit(err, front_msg=f'{args.json}, name: {album_info['name']}')
    pics = get_pics(album_info)
    make_album(pics, album_info, proj_info)
