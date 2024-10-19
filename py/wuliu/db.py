import json
import sqlite3
import sqlite3.Connection as Conn
from operator import itemgetter
from .const import *


"""
数据库用 sqlite, 但把 sqlite 当作 key-value 数据库来用。
因为 key-value 数据库就够用了，但多数 key-value 数据库都是语言限定的，
不同用，因此用这个办法，让 sqlite 变成一个各种语言通用的 key-value 数据库。

具体来说，平时通常把一个表的全部条目读出来，保存到一个 `dict` 中，成为 cache,
如果只涉及「读」操作，一般就只使用 cache, 如果涉及「写」操作则需要同时更新 cache 和数据库。
"""


Create_Table_File = "CREATE TABLE file(id TEXT PRIMARY KEY, doc TEXT)"
Insert_File = "INSERT INTO file(id, doc) VALUES(?, ?)"
Select_File = "SELECT id, doc FROM file"


def open_db(db_path) -> Conn:
    return sqlite3.connect(db_path)


def db_create_tables(db: Conn):
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


def db_dup_checksum(db: Conn) -> dict:
    """
    尋找数据库中重複的 checksum
    返回 dict, key 是 checksum, value 是 files
    """
    cache = db_cache(db)
    count: dict[str, dict] = {}
    for f in cache.values():
        items = count.get(f[CHECKSUM], list())
        count[f[CHECKSUM]] = items.append(f)
    return {k:v for (k,v) in count.items() if len(v)>1}
