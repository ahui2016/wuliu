import shutil
import argparse
import humanize
from sqlite3 import Connection as Conn
from pathlib import Path

from wuliu.const import *
from wuliu.common import read_project_info
from wuliu.db import open_db, db_select_by_id, db_cache


files_folder = Path(FILES)
buffer_folder = Path(BUFFER)
meta_folder = Path(METADATA)


def checkSizeLimit(size: int, info: dict):
    limit = info["ExportSizeLimit"] * MB
    if size > limit:
        size_str = humanize.naturalsize(size)
        raise WuliuError(f"檔案體積({size_str}) 超過上限({info["ExportSizeLimit"]} MB)")


def export_file(file: dict, info: dict):
    checkSizeLimit(file[SIZE], info)
    src = files_folder.joinpath(file[FILENAME])
    dst = buffer_folder.joinpath(file[FILENAME])
    if dst.exists():
        print(f"[warning] file exists: {dst}")
        return
    print(f"Export => {dst}")
    shutil.copyfile(src, dst)


def export_file_by_id(file_id: str, db: Conn, info: dict):
    file = db_select_by_id(file_id, db)
    export_file(file, info)


def export_meta(file_id: str, db: Conn):
    file = db_select_by_id(file_id, db)
    src = meta_folder.joinpath(file[FILENAME] + ".json")
    dst = buffer_folder.joinpath(file[FILENAME] + ".json")
    if dst.exists():
        print(f"[warning] file exists: {dst}")
        return
    print(f"Export => {dst}")
    shutil.copyfile(src, dst)


def export_file_by_collection(coll: str, db: Conn, info: dict, danger: bool | None):
    cache = db_cache(db)
    files = [f for f in cache.values() if coll in f[COLLECTIONS]]
    if not files:
        print(f"找不到與 collection:{coll} 匹配的檔案")
        return
    if len(files) > 10 and not danger:
        print(f"找到 {len(files)} 個檔案, 請使用 '-danger' 參數執行導出。")
        return
    for f in files:
        export_file(f, info)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()

    parser.add_argument(
        "-danger", action="store_true", help="allow dangerous operations"
    )
    parser.add_argument("-file", type=str, help="specify a file ID and export the file")
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
    db = open_db(Project_PY_DB)

    if args.file:
        export_file_by_id(args.file, db, info)
    elif args.meta:
        export_meta(args.meta, db)
    elif args.id:
        export_file_by_id(args.id, db, info)
        export_meta(args.id, db)
    elif args.coll:
        export_file_by_collection(args.coll, db, info, args.danger)
    else:
        parser.print_help()

    db.close()
