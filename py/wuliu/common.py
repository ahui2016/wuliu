import re
import sys
import json
import msgpack
from typing import Union
from pathlib import Path
from PIL import Image, ImageOps
from .const import *


def read_project_info():
    data = Path(Project_JSON).read_text()
    return json.loads(data)


def read_thumbs_msgp() -> dict:
    file = Path(Thumbs_msgp)
    if not file.exists():
        return dict()
    data = file.read_bytes()
    return msgpack.unpackb(data)


def open_image(file: str|Path) -> Union[Image, None]:
    try:
        img = Image.open(file)
    except OSError:
        img = None
    return img


def create_thumb_img(img:Image, thumb_path:Path, thumb_size):
    """請使用 create_thumb
    """
    img = ImageOps.exif_transpose(img)
    img = img.convert("RGB")
    img = ImageOps.fit(img, thumb_size)
    img.save(thumb_path)


def create_thumb(pic_path, thumb_path, thumb_size) -> str | None:
    img = open_image(pic_path)
    if img is None:
        return f"Not Image: {pic_path}"
    create_thumb_img(img, thumb_path, thumb_size)
    return None


def files_to_pics(files):
    return {file[ID]: file for file in files if file[Type].startswith('image')}


def check_filename(name: str) -> str:
    """
    :return: 有错返回 err: str, 无错返回空字符串。
    """
    if Filename_Forbid_Pattern.search(name) is None:
        return ''
    else:
        return '只能使用 0-9, a-z, A-Z, _(下劃線), -(連字號), .(點)' \
               '\n注意：不能使用空格，請用下劃線或連字號替代空格。'


def print_err(err:str|None):
    """如果有错误就打印, 没错误就忽略."""
    if err:
        print(f"Error! {err}", file=sys.stderr)


def print_err_exit(err:str|None, front_msg:str=''):
    """若有错误则打印并结束程序, 无错误则忽略.
    如果提供了 front_msg, 则在 err 之前显示。
    """
    if err:
        if front_msg:
            print(f'Error! {front_msg}', file=sys.stderr)
            sys.exit(err)
        else:
            sys.exit(f'Error! {err}')
