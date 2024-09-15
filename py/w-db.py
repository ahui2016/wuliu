import sys
import json
import argparse

from pathlib import Path
from tinydb import TinyDB, Query
from tinydb.storages import JSONStorage
from tinydb.middlewares import CachingMiddleware

from wuliu.const import *
from wuliu.common import print_err, print_err_exit, read_project_info, check_not_in_backup


def open_db() -> TinyDB:
    return TinyDB(Project_PY_DB, storage=CachingMiddleware(JSONStorage))


def load_all_metadatas(db: TinyDB):
    files = Path(METADATA).glob("*.json")
    meta_list = []
    for f in files:
        print(".", end="")
        data = f.read_text(encoding="utf-8")
        meta = json.loads(data)
        meta_list.append(meta)
    db.insert_multiple(meta_list)


def create_database():
    """
    必须确保数据库不存在，创建一个新的数据库。
    """
    if Path(Project_PY_DB).exists():
        print_err_exit(f"file exists: {Project_PY_DB}")
    print(f"Create {Project_PY_DB}")
    with open_db() as db:
        load_all_metadatas(db)


if __name__ == "__main__":

    parser = argparse.ArgumentParser()

    parser.add_argument("--create", action="store_true",
        help="create the database")

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    if args.create:
        create_database()
        sys.exit()

    parser.print_help()

"""
    with open_db() as db:
        File = Query()
        result = db.search(File.filename.matches("香港*"))
        print(result)
"""
