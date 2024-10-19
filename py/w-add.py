import sys
import shutil
import argparse
import humanize
import sqlite3
import sqlite3.Connection as Conn
from pathlib import Path

from wuliu.const import *
from wuliu.common import (
    print_err,
    print_err_exit,
    read_project_info,
    check_not_in_backup,
    get_filenames,
    time_now,
    name_to_id,
    file_sum512,
    type_by_filename,
    json_dumps,
    yaml_dump,
    yaml_load_file,
    check_keywords,
)
from wuliu.db import open_db, db_cache, db_files_exist


input_folder = Path(INPUT)
files_folder = Path(FILES)
meta_folder = Path(METADATA)


def new_file(name: str) -> dict:
    f = New_File()
    f[ID] = name_to_id(name)
    f[FILENAME] = name
    now = time_now()
    f[CTIME] = now
    f[UTIME] = now
    return f


def new_files_from(names: list[str], folder: str) -> list:
    """
    把档案名 (names) 转换为 files, 此时假设档案在 folder 资料夹内。
    """
    files = []
    folder_path = Path(folder)
    for name in names:
        file_path = folder_path.joinpath(name)
        file_stat = file_path.lstat()
        checksum = file_sum512(file_path)
        f = new_file(name)
        f[CHECKSUM] = checksum
        f[SIZE] = file_stat.st_size
        f[TYPE] = type_by_filename(name)
        files.append(f)
    return files


def read_config(file_path: Path) -> dict:
    """返回 None 表示有错误"""
    cfg = yaml_load_file(file_path)
    if len(cfg[IDS]) > 0:
        print_err_exit("添加新檔案時不可通過 ID 指定檔案")

    keywords = cfg[KEYWORDS] + cfg[COLLECTIONS] + cfg[ALBUMS]
    err = check_keywords(keywords)
    print_err_exit(err)
    return cfg


def config_add_files(ids: list, filenames: list) -> dict:
    cfg = Edit_Files_Config()
    cfg[IDS] = ids
    cfg[FILENAMES] = filenames
    return cfg


def create_config_yaml(filename: str, allow_danger: bool):
    file_path = Path(filename)
    if file_path.exists() and not allow_danger:
        print_err(f"file exists: {filename}")
        return
    names_in_input = get_filenames(input_folder)
    text = yaml_dump(config_add_files([], names_in_input))
    print(f"Create => {filename}")
    file_path.write_text(text, encoding="utf-8")


def find_input_files(cfg_path: str):
    """
    寻找 input 资料夹里的全部档案。
    :return: (files, cfg)
    """
    names_in_input = get_filenames(input_folder)
    if not cfg_path:
        return new_files_from(names_in_input, INPUT), None

    cfg = read_config(Path(cfg_path))
    if len(cfg[FILENAMES]) == 0:
        cfg[FILENAMES] = names_in_input

    filenames = []
    for name in cfg[FILENAMES]:
        if name in names_in_input:
            filenames.append(name)
        else:
            print(f"Not Found: {name}")

    files = new_files_from(filenames, INPUT)
    for f in files:
        f[LIKE] = cfg[LIKE]
        f[LABEL] = cfg[LABEL]
        f[NOTES] = cfg[NOTES]
        f[KEYWORDS] = cfg[KEYWORDS]
        f[COLLECTIONS] = cfg[COLLECTIONS]
        f[ALBUMS] = cfg[ALBUMS]
    return files, cfg


def print_files(files: list, cfg: dict | None):
    if len(files) == 0:
        print("在input資料夾中未發現新檔案")

    print("【待執行操作如下所示(未正式執行)】")
    for f in files:
        size = humanize.naturalsize(f[SIZE])
        size = f"({size})"
        print(f"{size.ljust(11, ' ')} {f[FILENAME]}")

    if not cfg:
        return

    print(f"like: {cfg[LIKE]}")
    print(f"label: {cfg[LABEL]}")
    print(f"notes: {cfg[NOTES]}")
    print(f"keywords: {', '.join(cfg[KEYWORDS])}")
    print(f"collections: {', '.join(cfg[COLLECTIONS])}")
    print(f"albums: {', '.join(cfg[ALBUMS])}")


def print_id_name(files: list):
    for f in files:
        print(f"{f[ID]}: {f[FILENAME]}")


def check_exist(files: list, cache: dict) -> bool:
    """
    :return: has_exist_files
    """
    exist_files = db_files_exist(files, cache)
    if len(exist_files) > 0:
        print("【注意！】數據庫中有名稱或内容相同的檔案：")
        print_id_name(exist_files)
        return True

    dst_files = []
    for f in files:
        filename = f[FILENAME]
        dst = files_folder.joinpath(filename)
        meta = meta_folder.joinpath(filename + ".json")
        dst_files.extend([dst, meta])

    exist_files = []
    for f in dst_files:
        if f.exists():
            exist_files.append(f)

    if len(exist_files) > 0:
        print("【注意！】同名檔案已存在：")
        for f in exist_files:
            print(f)
        return True

    return False


# 添加档案时不需要检查磁盘空间，
# 因为 input 资料夹与 files 资料夹在同一个磁盘分区内。
def add_files(files: list, db: Conn):
    if len(files) == 0:
        print("warning: No file to add.")
        return

    for f in files:
        filename = f[FILENAME]
        src = input_folder.joinpath(filename)
        dst = files_folder.joinpath(filename)
        print(f"Add => {dst}")
        shutil.move(src, dst)

        meta_path = meta_folder.joinpath(filename + ".json")
        print(f"Create => {meta_path}")
        text = json_dumps(f)
        meta_path.write_text(text, encoding="utf8")

        db.insert(f)

    print("Done.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()

    parser.add_argument(
        "-danger", action="store_true", help="allow dangerous operations"
    )

    parser.add_argument(
        "--new-yaml", type=str, help="create a YAML file for adding files"
    )

    parser.add_argument("-yaml", type=str, help="use a YAML file to add files")

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    if args.new_yaml:
        create_config_yaml(args.new_yaml, args.danger)
        sys.exit()

    files, cfg = find_input_files(args.yaml)

    db = open_db(Project_PY_DB)
    cache = db_cache()

    if check_exist(files, cache):
        db.close()
        sys.exit()

    if args.danger:
        add_files(files, db)
        db.close()
        sys.exit()

    db.close()
    print_files(files, cfg)
