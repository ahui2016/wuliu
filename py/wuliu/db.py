import json
import sqlite3
import sqlite3.Connection as Conn
from operator import itemgetter
from .const import *


Create_Table_File = "CREATE TABLE file(id TEXT PRIMARY KEY, doc TEXT)"
Insert_File = "INSERT INTO file(id, doc) VALUES(?, ?)"
Select_File = "SELECT id, doc FROM file"


def open_db(db_path) -> Conn:
    return sqlite3.connect(db_path)


def db_init(db: Conn):
    with db:
        db.execute(Create_Table_File)


def db_cache(db: Conn) -> dict:
    with db:
        data = db.execute(Select_File).fetchall()
        return {k:v for (k,v) in data}


def db_insert_many(data: list, db: Conn):
    with db:
        db.executemany(Insert_File, data)


def db_insert_files(files: list, db: Conn):
    data = files_to_pairs(files)
    db_insert_many(data, db)


def file_to_pair(file: dict) -> tuple:
    """
    convert a dict to a key-value pair,
    where the key is the id, and the value is a JSON.
    """
    return (file[ID], json.dumps(file))


def files_to_pairs(files: list) -> list:
    return [file_to_pair(file) for file in files]


def db_insert_file(file:dict, db: Conn):
    with db:
        db.execute(Insert_File, file_to_pair(file))


def db_files_exist(files: list, cache: dict) -> list:
    """
    :return: exist_files 名稱或內容相同的檔案
    """
    checksums = [f[CHECKSUM] for f in cache.values()]
    exist_files = []
    for f in files:
        if (f[ID] in cache) or (f[CHECKSUM] in checksums):
            exist_files.append(f)
    return exist_files


def db_get_files(cache: dict, n: int | None, orderby: str | None) -> list:
    """
    :orderby: size/like/ctime/utime (default "utime")
    :n: n < 0 表示全部, n 等於 None 或 0 表示默認值(5)
    """
    if orderby not in ["size", "like", "ctime"]:
        orderby = "utime"

    if orderby == "like":
        files = [file for file in cache.values() if file[LIKE] > 0]
    else:
        files = cache.values()

    files.sort(key=itemgetter(orderby), reverse=True)

    if n is None or n == 0:
        n = 5

    files_len = len(files)
    if n < 0 or n >= files_len:
        return files

    return files[:n]


def db_all_ids(db: TinyDB) -> set:
    """
    導出全部檔案ID
    """
    files = db.all()
    return {f[ID] for f in files}


def db_dup_id(files: list, db: TinyDB) -> dict:
    """尋找重複的ID"""
    # files = db.all(), 另一個函數也要使用 files
    # 檔案名受檔案系統的限制, 不會重複, 因此只需要檢查ID
    id_count: dict[str, int] = {}
    for f in files:
        n = id_count.get(f[ID], 0)
        id_count[f[ID]] = n + 1
    duplicated = {}
    for k, v in id_count.items():
        if v > 1:
            docs = db.search(File.id == k)
            duplicated[k] = [dict(doc) for doc in docs]
    return duplicated


def db_dup_checksum(files: list, db: TinyDB) -> dict:
    """尋找重複的 checksum (意味着檔案內容完全相同)"""
    # files = db.all(), 另一個函數也要使用 files
    count: dict[str, int] = {}
    for f in files:
        n = count.get(f[CHECKSUM], 0)
        count[f[ID]] = n + 1
    duplicated = {}
    for k, v in count.items():
        if v > 1:
            docs = db.search(File.id == k)
            duplicated[k] = [dict(doc) for doc in docs]
    return duplicated
