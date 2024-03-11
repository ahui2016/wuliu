import msgpack
from .const import *


def files_to_pics(files):
    return {file[ID]: file for file in files if file[Type].startswith('image')}


def dump_pics(pics, path):
    """把 pics 寫入 path, 方便後續使用。
    """
    print(f"Write to {path}")
    msgpack.dump(pics, path)
