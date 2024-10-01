import json
import shutil
from pathlib import Path
from tinydb import TinyDB, Query
from wuliu.const import *
from wuliu.common import (
    print_err_exit,
    get_filenames,
    type_by_filename,
    yaml_dump,
    yaml_load_file,
    file_sum512,
    time_now,
    path_write_json,
)


File = Query()


def files_or_meta(filetype: str) -> str:
    if filetype == "text/json":
        return "metadata"
    return "files"


def buffer_files() -> dict:
    names = get_filenames(BUFFER)
    cfg = {}
    for name in names:
        file_type = type_by_filename(name)
        target = files_or_meta(file_type)
        cfg[name] = target
    if not cfg:
        print_err_exit("在buffer資料夾中未發現檔案")
    return cfg


def new_cfg_yaml(cfg_path: Path):
    text = yaml_dump(buffer_files())
    print(f"Create => {cfg_path}")
    cfg_path.write_text(text, encoding="utf-8")


def read_config(cfg_path: Path) -> dict:
    cfg = yaml_load_file(cfg_path)
    if not cfg:
        print_err_exit("指定的 yaml 沒有內容")
    return cfg


def check_dst(dst: Path):
    if not dst.exists():
        print_err_exit(f"無法覆蓋不存在的檔案: {dst}")


def print_preview(cfg: dict):
    for name, target in cfg.items():
        print(f"{target} <== buffer/{name}")
        dst = Path(target).joinpath(name)
        check_dst(dst)


def overwrite_into_files(name: str, dst: Path, db: TinyDB):
    src = Path(BUFFER).joinpath(name)
    meta_path = Path(METADATA).joinpath(name+".json")
    meta = json.load(meta_path)

    sum = file_sum512(src)
    if sum == meta[CHECKSUM]:
        print(f"檔案內容沒有變化: {name}")
        return

    meta[CHECKSUM] = sum
    meta[UTIME] = time_now()
    file_stat = src.lstat()
    meta[SIZE] = file_stat.st_size

    shutil.move(src, dst)
    path_write_json(meta_path, meta)

    updated = {
        CHECKSUM: meta[CHECKSUM],
        SIZE: meta[SIZE],
        UTIME: meta[UTIME],
    }
    db.update(updated, File.id == meta[ID])


def overwrite_into_meta(name: str, dst: Path, db: TinyDB):
    src = Path(BUFFER).joinpath(name)
    meta = json.load(src)
    shutil.move(src, dst)
    db.update(meta, File.id == meta[ID])


def overwrite_files(cfg:dict, db: TinyDB):
    for name, target in cfg.items():
        print(f"{target} <== buffer/{name}")
        dst = Path(target).joinpath(name)
        check_dst(dst)
        if target == FILES:
            overwrite_into_files(name, dst, db)
        elif target == METADATA:
            overwrite_into_meta(name, dst, db)
        else:
            print_err_exit(f"不認識: {target}")


