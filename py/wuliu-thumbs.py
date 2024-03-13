from pathlib import Path
import sys
import argparse
import msgpack
from wuliu.const import *
from wuliu.common import print_err, read_project_info, read_thumbs_msgp, create_thumb


def get_pics_metadata(msgp_path:Path):
    """獲取全部圖片的屬性
    
    msgp_path 是由 "wuliu-db -dump" 導出的數據。
    :return: 返回 File 列表 (參考 util/model.go 裏的 File)
    """
    with msgp_path.open(mode='rb') as f:
        return msgpack.load(f)


def pic_in_thumbs(pic, thumbs) -> bool:
    pic_id = pic[ID]
    old_checksum = thumbs.get(pic_id)
    return old_checksum == pic[Checksum]


def updated_pics(pics, thumbs):
    """
    :return: pics need to create or re-create thumbnails
    """
    newpics = dict()
    for pic in pics:
        if not pic_in_thumbs(pic, thumbs):
            newpics[pic[ID]] = pic
    return newpics


def create_thumbs(pics, thumb_size):
    thumbs = read_thumbs_msgp()
    newpics = updated_pics(pics, thumbs)
    if len(newpics) == 0:
        return

    for pic in newpics:
        pic_id = pic[ID]
        pic_path = Path(Files).joinpath(pic[Filename])
        thumb_path = Path(Thumbs).joinpath(f'{pic_id}.jpg')
        print(f'Create -> {thumb_path}')
        err = create_thumb(pic_path, thumb_path, thumb_size)
        print_err(err)
        if err is None:
            thumbs[pic_id] = pic[Checksum]

    print(f'Update -> {Thumbs_msgp}')
    blob = msgpack.packb(thumbs)
    Path(Thumbs_msgp).write_bytes(blob)


# ↓↓↓ main ↓↓↓ 

parser = argparse.ArgumentParser()

parser.add_argument('-msgp', type=str, default='',
        help='the msgp file created by "wuliu-db -dump"')

args = parser.parse_args()

if args.msgp == '':
    print('wuliu-pics.py: error: required argument -msgp')
    sys.exit(1)

proj_info = read_project_info()

pics = get_pics_metadata(Path(args.msgp))
create_thumbs(pics, proj_info[Thumb_Size])
