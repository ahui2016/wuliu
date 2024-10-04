import sys
import json
import argparse

from pathlib import Path
from tinydb import TinyDB

from wuliu.const import *
from wuliu.common import (
    print_err_exit,
    read_project_info,
    check_not_in_backup,
    yaml_dump,
)
from wuliu.db import open_db, db_all_ids, db_dup_id, db_dup_checksum


def load_all_metadatas(db: TinyDB):
    files = Path(METADATA).glob("*.json")
    meta_list = []
    for f in files:
        print(".", end="")
        data = f.read_text(encoding="utf-8")
        meta = json.loads(data)
        meta_list.append(meta)
    db.insert_multiple(meta_list)

    all_files = db.all()
    duplicated = db_dup_id(all_files, db)
    if duplicated:
        print("【發現重複ID】")
        msg = yaml_dump(duplicated)
        print(msg)
        print()

    duplicated = db_dup_checksum(all_files, db)
    if duplicated:
        print("【發現重複檔案】")
        msg = yaml_dump(duplicated)
        print(msg)


def create_database(db_path: str):
    """
    必须确保数据库不存在，创建一个新的数据库。
    """
    if Path(db_path).exists():
        print_err_exit(f"file exists: {db_path}")
    print(f"Create {db_path}")
    with open_db(db_path) as db:
        load_all_metadatas(db)


def update_database_add_only(db_path: str):
    """
    發現 metadata 資料夾中的新檔案並添加到數據庫中。
    注意: 本函數只處理新檔案, 忽略檔案被刪除或修改的情況,
    如須確保數據庫完全更新, 請刪除數據庫後重新建立。
    """
    if not Path(db_path).exists():
        print_err_exit("未找到數據庫, 請使用 `w-db --create` 創建數據庫")
    print("尋找新檔案...")
    metafiles = Path(METADATA).glob("*.json")
    new_files = []
    with open_db(db_path) as db:
        ids = db_all_ids(db)
        for f in metafiles:
            data = f.read_text(encoding="utf-8")
            meta = json.loads(data)
            if meta[ID] in ids:
                continue
            new_files.append(meta)

        if not new_files:
            print("未發現新檔案")
            return

        print(f"發現 {len(new_files)} 個新檔案, 正在插入數據庫...")
        db.insert_multiple(new_files)
    print("OK")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()

    parser.add_argument("--create", action="store_true", help="create the database")

    parser.add_argument("--update", type=str, help="--update=add")

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    if args.create:
        create_database(Project_PY_DB)
        sys.exit()

    if args.update:
        if args.update != "add":
            print_err_exit("參數 '--update' 目前只接受 'add', 例: w-db --update=add")
        update_database_add_only(Project_PY_DB)
        sys.exit()

    parser.print_help()

"""
    with open_db() as db:
        File = Query()
        result = db.search(File.filename.matches("香港*"))
        print(result)
"""
