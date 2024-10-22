import re
import argparse
from sqlite3 import Connection as Conn
from pathlib import Path

from wuliu.const import *
from wuliu.common import read_project_info, json_load, type_by_filename, name_to_id, path_write_json
from wuliu.db import open_db, db_cache, db_insert_file, db_delete_file, db_select_by_id


files_folder = Path(FILES)
meta_folder = Path(METADATA)


def check_exist(src: Path, dst: Path):
    if not src.exists():
        raise WuliuError(f"not found: {src}")
    if not src.is_file():
        raise WuliuError(f"not a file: {src}")
    if dst.exists():
        raise WuliuError(f"file exists: {dst}")


def check_filename(name):
    if re.search(r'[\\\/\:\*\?\"\<\>\|]', name):
        raise ValueError(r'檔案名稱不允許包含這些字符 \/:*?"<>|')


def rename_file(old_name: str, new_name: str):
    src = files_folder.joinpath(old_name)
    dst = files_folder.joinpath(new_name)
    print(f"Rename {src} => {dst}")
    check_exist(src, dst)
    src.rename(dst)


def rename_meta(old_name: str, new_name: str):
    src = meta_folder.joinpath(old_name+".json")
    dst = meta_folder.joinpath(new_name+".json")
    print(f"Rename {src} => {dst}")
    check_exist(src, dst)
    meta = json_load(src)
    meta[ID] = name_to_id(new_name)
    meta[FILENAME] = new_name
    meta[TYPE] = type_by_filename(new_name)
    path_write_json(dst, meta)
    src.unlink()


def rename_in_db(old_id: str, newfile: dict, db: Conn):
    print(f"insert into db => {newfile[FILENAME]}")
    db_insert_file(newfile, db)
    print(f"delete from db => id:{old_id}")
    db_delete_file(old_id, db)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()

    parser.add_argument("-name", type=str, help="set a new filename")
    parser.add_argument(
        "-meta", type=str, help="specify a file ID and export the file's metadata(json)"
    )
    parser.add_argument(
        "-id", type=str, help="specify a file ID and export the file and its metadata"
    )
    parser.add_argument("-coll", type=str, help="export files by collection")

    args = parser.parse_args()
    info = read_project_info()
    # check_not_in_backup(info)

    old_file = db_select_by_id(old_id, db)

    check_filename(args.name)
