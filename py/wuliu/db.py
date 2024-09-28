from tinydb import TinyDB, Query
from tinydb.storages import JSONStorage
from tinydb.middlewares import CachingMiddleware
from operator import itemgetter
from .const import *


def open_db(db_path) -> TinyDB:
    return TinyDB(db_path, storage=CachingMiddleware(JSONStorage))


def files_in_db(files: list, db: TinyDB) -> list:
    """
    :return: exist_files 名稱或內容相同的檔案
    """
    exist_files = []
    File = Query()
    for f in files:
        ef = db.get((File.id == f[ID]) | (File.checksum == f[CHECKSUM]))
        if ef is not None:
            exist_files.append(ef)
    return exist_files


def db_new_files(db: TinyDB, n: int | None) -> list:
    """
    返回 n 个最新添加到数据库中的档案。 n 不可大于 100
    """
    if not n:
        n = 10

    if n > 100:
        n = 100

    db_size = len(db)
    if n >= db_size:
        files = db.all()
        return [dict(f) for f in reversed(files)]

    skip = db_size - n
    files = []
    for f in db:
        if skip > 0:
            skip -= 1
            continue
        files.append(dict(f))
    files.reverse()
    return files


def db_all_files(db: TinyDB, orderby: str | None) -> list:
    """
    :orderby: size/like/utime (default "utime")
    """
    if orderby != "size" and orderby != "like":
        orderby = "utime"

    if orderby == "like":
        File = Query()
        files = db.search(File.like > 0)
    else:
        files = db.all()

    files = [dict(f) for f in files]
    files.sort(key=itemgetter(orderby), reverse=True)
    return files
