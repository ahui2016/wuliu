import sys
import json
import argparse

from sqlite3 import Connection as Conn
from pathlib import Path

from wuliu.const import *
from wuliu.common import (
    print_err_exit,
    read_project_info,
    check_not_in_backup,
    yaml_dump,
)
from wuliu.db import (
    open_db,
    db_create_tables,
    db_insert_files,
    db_dup_checksum,
    db_cache,
)


def load_all_metadatas(db: Conn):
    """
    向新数据库写入原有的 metadata
    """
    files = Path(METADATA).glob("*.json")
    meta_list = []
    for f in files:
        data = f.read_text(encoding="utf-8")
        meta = json.loads(data)
        meta_list.append(meta)
    db_insert_files(meta_list, db)

    cache = db_cache(db)
    duplicated = db_dup_checksum(cache)
    if duplicated:
        print("【發現重複檔案】")
        msg = yaml_dump(duplicated)
        print(msg)


def create_database(db_path: str):
    """
    创建一个新的数据库。
    """
    if Path(db_path).exists():
        print_err_exit(f"file exists: {db_path}")
    print(f"Create {db_path}")
    print("(如果檔案較多, 請耐心等待...)")

    db = open_db(db_path)
    db_create_tables(db)
    load_all_metadatas(db)
    db.close()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()

    parser.add_argument("--create", action="store_true", help="create the database")

    # parser.add_argument("--update", type=str, help="--update=add")

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    if args.create:
        create_database(Project_PY_DB)
        sys.exit()

    parser.print_help()
