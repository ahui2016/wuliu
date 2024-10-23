import re
import argparse
from sqlite3 import Connection as Conn
from pathlib import Path

from wuliu.const import *
from wuliu.common import (
    read_project_info,
    check_not_in_backup,
    json_load,
    type_by_filename,
    name_to_id,
    path_write_json,
    print_err,
)
from wuliu.db import open_db, db_insert_file, db_delete_file, db_select_by_id


files_folder = Path(FILES)
meta_folder = Path(METADATA)


def check_exist(src: Path, dst: Path):
    if not src.exists():
        raise WuliuError(f"not found: {src}")
    if not src.is_file():
        raise WuliuError(f"not a file: {src}")
    if dst.exists():
        raise WuliuError(f"file exists: {dst}")


def check_filename(old_name: str, new_name: str):
    if old_name == new_name:
        raise WuliuError("新檔案名與舊檔案名相同")
    # if re.search(r'[\\\/\:\*\?\"\<\>\|]', new_name):
    if re.search(r'[\\/:*?"<>|]', new_name):
        raise WuliuError(r'檔案名稱不允許包含這些字符 \/:*?"<>|')


def rename_file(old_name: str, new_name: str):
    src = files_folder.joinpath(old_name)
    dst = files_folder.joinpath(new_name)
    print(f"Rename {src} => {dst}")
    check_exist(src, dst)
    src.rename(dst)


def rename_meta(old_name: str, new_name: str) -> dict:
    src = meta_folder.joinpath(old_name + ".json")
    dst = meta_folder.joinpath(new_name + ".json")
    print(f"Rename {src} => {dst}")
    check_exist(src, dst)
    meta = json_load(src)
    meta[ID] = name_to_id(new_name)
    meta[FILENAME] = new_name
    meta[TYPE] = type_by_filename(new_name)
    path_write_json(dst, meta)
    src.unlink()
    return meta


def rename_in_db(old_id: str, new_meta: dict, db: Conn):
    print("Update database...")
    db_insert_file(new_meta, db)
    db_delete_file(old_id, db)


def rename(old_id: str, new_name: str, db: Conn):
    try:
        old_meta = db_select_by_id(old_id, db)
        old_name = old_meta[FILENAME]
        check_filename(old_name, new_name)
    except WuliuError as err:
        print_err(str(err))
        return
    rename_file(old_name, new_name)
    new_meta = rename_meta(old_name, new_name)
    rename_in_db(old_id, new_meta, db)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-id", type=str, help="specify a file ID that exists")
    parser.add_argument("-name", type=str, help="set a new filename")

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    db = open_db(Project_PY_DB)
    if args.id and (not args.name):
        print_err("Required '-name'")
    elif args.name and (not args.id):
        print_err("Required '-id'")
    elif args.id and args.name:
        rename(args.id, args.name, db)
    else:
        parser.print_help()
    db.close()
