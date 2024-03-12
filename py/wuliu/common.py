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


def print_err(err:str):
    """如果有错误就打印, 没错误就忽略."""
    if err:
        print(f"Error! {err}", file=sys.stderr)


def print_err_exist(err:str):
    """若有错误则打印并结束程序, 无错误则忽略."""
    if err:
        sys.exit(f"Error! {err}")
