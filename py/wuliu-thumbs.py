import sys
import argparse
import msgpack
from typing import List, Dict, Set
from pathlib import Path
from wuliu.const import *
from wuliu.common import print_err, read_project_info, read_thumbs_msgp, create_thumb


def get_pics_metadata(msgp_path:Path) -> List:
    """獲取全部圖片的屬性
    
    msgp_path 是由 "wuliu-db -dump" 導出的數據。
    :return: 返回 File 列表 (參考 util/model.go 裏的 File)
    """
    with msgp_path.open(mode='rb') as f:
        return msgpack.load(f)


def pic_in_thumbs(pic:dict, thumbs:dict) -> bool:
    pic_id = pic[ID]
    old_checksum = thumbs.get(pic_id)
    return old_checksum == pic[Checksum]


def get_thumbs_orphans() -> Set[str]:
    """返回孤立縮略圖的ID"""
    current_thumbs = read_thumbs_msgp()
    current_ids = set(current_thumbs.keys())
    thumbs_paths = Path(Thumbs).glob('*.*')
    thumbs_ids = {filename.stem for filename in thumbs_paths}
    return thumbs_ids.difference(current_ids)


def updated_pics(pics:List, thumbs:Dict) -> List:
    """
    :return: pics need to create or re-create thumbnails
    """
    newpics = []
    for pic in pics:
        if not pic_in_thumbs(pic, thumbs):
            newpics.append(pic)
    return newpics


def deleted_pics_ids(pics:List, thumbs:Dict) -> Set[str]:
    """返回需要刪除的縮略圖的ID"""
    ids: Set[str] = set()
    pics_ids = {pic[ID] for pic in pics}
    thumbs_ids = {pic_id for pic_id in thumbs.keys()}
    deleted_ids = thumbs_ids.difference(pics_ids)
    return deleted_ids


def delete_thumbs(deleted_ids:Set[str], thumbs:Dict) -> Dict:
    """返回更新後的 thumbs"""
    for pic_id in deleted_ids:
        thumb_path = Path(Thumbs).joinpath(f'{pic_id}.jpg')
        print(f'Delete -> {thumb_path}')
        thumb_path.unlink(missing_ok=True)
        if pic_id in thumbs:
            del thumbs[pic_id]
    return thumbs


def delete_orphans(orphans_ids:Set[str]):
    """刪除多餘的縮略圖"""
    if len(orphans_ids) == 0:
        print('未發現多餘的縮略圖')
        return
    for pic_id in orphans_ids:
        thumb_path = Path(Thumbs).joinpath(f'{pic_id}.jpg')
        print(f'Delete -> {thumb_path}')
        thumb_path.unlink(missing_ok=True)


def create_thumbs(thumb_size, newpics:List, thumbs:Dict) -> Dict:
    """返回更新後的 thumbs"""
    for pic in newpics:
        pic_id = pic[ID]
        pic_path = Path(Files).joinpath(pic[Filename])
        thumb_path = Path(Thumbs).joinpath(f'{pic_id}.jpg')
        print(f'Create -> {thumb_path}')
        err = create_thumb(pic_path, thumb_path, thumb_size)
        print_err(err)
        if err is None:
            thumbs[pic_id] = pic[Checksum]
    return thumbs


def update_thumbs(pics:List, thumb_size):
    """生成縮略圖及刪除縮略圖"""
    thumbs = read_thumbs_msgp()
    newpics = updated_pics(pics, thumbs)
    deleted_ids = deleted_pics_ids(pics, thumbs)
    if len(newpics) + len(deleted_ids) == 0:
        print("未發現圖片變化 (無新增、無更新、無刪除)")
        return

    thumbs = delete_thumbs(deleted_ids, thumbs)
    thumbs = create_thumbs(thumb_size, newpics, thumbs)

    print(f'Update -> {Thumbs_msgp}')
    blob = msgpack.packb(thumbs)
    Path(Thumbs_msgp).write_bytes(blob)


# ↓↓↓ main ↓↓↓ 

parser = argparse.ArgumentParser()

parser.add_argument('-msgp', type=str, default='',
        help='the msgp file created by "wuliu-db -dump pics"')

parser.add_argument('--delete-orphans', action='store_true',
        help='delete orphans of the thumbnails')

args = parser.parse_args()

if args.msgp == '' and not args.delete_orphans:
    parser.print_help()
    sys.exit(1)

if args.msgp != '':
    proj_info = read_project_info()
    pics = get_pics_metadata(Path(args.msgp))
    update_thumbs(pics, proj_info[Thumb_Size])
    sys.exit(0)

if args.delete_orphans:
    orphans = get_thumbs_orphans()
    delete_orphans(orphans)
