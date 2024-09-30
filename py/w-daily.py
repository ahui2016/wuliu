import sys
import json
import shutil
import argparse
from pathlib import Path
from tinydb import TinyDB, Query
from wuliu.const import *
from wuliu.common import (
    time_now,
    print_err,
    file_sum512,
    type_by_filename,
    read_project_info,
    check_not_in_backup,
)
from wuliu.db import open_db


File = Query()
input_folder = Path(INPUT)
files_folder = Path(FILES)
buffer_folder = Path(BUFFER)
meta_folder = Path(METADATA)


def today() -> str:
    return time_now()[:10]


def file_exists(file_id: str, db: TinyDB):
    """返回空字符串或 None 表示沒有錯誤"""
    result = db.get(File.id == file_id)
    file_path = files_folder.joinpath(file_id + ".html")
    return file_path.exists() or result


def export_file(file_id: str) -> str | None:
    """返回空字符串或 None 表示沒有錯誤"""
    filename = file_id + ".html"
    src = files_folder.joinpath(filename)
    dst = buffer_folder.joinpath(filename)
    if dst.exists():
        return f"file exists: {dst}"
    print(f"Export => {dst}")
    shutil.copyfile(src, dst)

    meta_file = filename + ".json"
    src = meta_folder.joinpath(meta_file)
    dst = buffer_folder.joinpath(meta_file)
    if dst.exists():
        return f"file exists: {dst}"
    print(f"Export => {dst}")
    shutil.copyfile(src, dst)

    return None


def new_file(file_id: str) -> dict:
    filename = file_id + ".html"
    file_path = files_folder.joinpath(filename)
    f = New_File()
    f[ID] = file_id
    f[FILENAME] = filename
    now = time_now()
    f[CTIME] = now
    f[UTIME] = now
    f[CHECKSUM] = file_sum512(file_path)
    file_stat = file_path.lstat()
    f[SIZE] = file_stat.st_size
    f[TYPE] = type_by_filename(filename)
    f[LABEL] = file_id
    f[COLLECTIONS] = [MY_DAILY]
    return f


def copyfile_or_not(src, dst):
    print(f"Copy => {dst}")
    if dst.exists():
        print(f"未執行複製! file exists: {dst}")
    else:
        shutil.copyfile(src, dst)


def create_daily(file_id: str, db: TinyDB):
    # 剛剛檢查過檔案不存在, 因此這裡不再檢查。
    filename = file_id + ".html"
    text = "<!DOCTYPE html>"
    src = files_folder.joinpath(filename)
    print(f"Create => {src}")
    src.write_text(text, encoding="utf-8")
    dst = buffer_folder.joinpath(filename)
    copyfile_or_not(src, dst)

    meta_filename = filename + ".json"
    file_meta = new_file(file_id)
    meta = json.dumps(file_meta, ensure_ascii=False, indent=4)
    src = meta_folder.joinpath(meta_filename)
    print(f"Create => {src}")
    src.write_text(meta, encoding="utf-8")
    dst = buffer_folder.joinpath(meta_filename)
    copyfile_or_not(src, dst)

    print("Insert into the database...")
    db.insert(file_meta)
    print("OK")


def create_export(day: str, db: TinyDB):
    file_id = DAILY_PREFIX + day
    if file_exists(file_id, db):
        err = export_file(file_id)
        print_err(err)
    else:
        create_daily(file_id, db)


if __name__ == "__main__":
    # 在 Windows 中使用 `>` 重定向打印到文件时可能会遇到编码问题，因此需要这行设置。
    # sys.stdout.reconfigure(encoding='utf-8')  # type: ignore

    parser = argparse.ArgumentParser()

    parser.add_argument("-edit", type=str, help="'-edit today' or '-edit 1970-12-31'")

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    if args.edit == "today":
        args.edit = today()

    if args.edit:
        with open_db(Project_PY_DB) as db:
            create_export(args.edit, db)
        sys.exit()

    parser.print_help()
