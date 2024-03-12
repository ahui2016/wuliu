from pathlib import Path
import sys
import argparse
import msgpack
from wuliu.const import *
from wuliu.common import print_err, read_project_info, create_thumb


def get_pics_metadata(msgp_path:Path):
    """獲取全部圖片的屬性
    
    msgp_path 是由 "wuliu-db -dump" 導出的數據。
    :return: 返回 File 列表 (參考 util/model.go 裏的 File)
    """
    with msgp_path.open(mode='rb') as f:
        return msgpack.load(f)


def create_thumbs(pics, thumb_size):
    for pic in pics:
        pic_path = Path(Files).joinpath(pic[Filename])
        thumb_path = Path(Thumbs).joinpath(f'{pic[ID]}.jpg')
        print(f'Create -> {thumb_path}')
        err = create_thumb(pic_path, thumb_path, thumb_size)
        print_err(err)


# ↓↓↓ main ↓↓↓ 

parser = argparse.ArgumentParser()

parser.add_argument('-msgp', type=str, default='',
        help='the msgp file created by "wuliu-db -dump"')

args = parser.parse_args()
print(args)

if args.msgp == '':
    print('wuliu-pics.py: error: required argument -msgp')
    sys.exit(1)

proj_info = read_project_info()

pics = get_pics_metadata(Path(args.msgp))
create_thumbs(pics, proj_info[Thumb_Size])
