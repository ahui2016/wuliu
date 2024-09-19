from tinydb import TinyDB, Query
from tinydb.storages import JSONStorage
from tinydb.middlewares import CachingMiddleware
from .const import *


def open_db(db_path) -> TinyDB:
    return TinyDB(db_path, storage=CachingMiddleware(JSONStorage))


def files_in_db(files: list, db: TinyDB) -> list:
    """
    :return: exist_files
    """
    exist_files = []
    File = Query()
    for f in files:
        ef = db.get(File.id == f[ID])
        if ef is not None:
            exist_files.append(ef)
    return exist_files
