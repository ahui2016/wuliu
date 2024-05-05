from pathlib import Path

from wuliu.albums import *
from wuliu.const import *

def get_docs_metadata() -> list:
    """獲取全部可預覽檔案的屬性
    
    Docs_msgp 是由 "wuliu-db -dump docs" 導出的數據，檔名固定為 docs.msgp
    :return: 返回 File 列表 (參考 util/model.go 裏的 File)
    """
    return get_files_metadata(Docs_msgp)


def read_docs_msgp(album_path: Path) -> dict:
    return read_album_msgp(album_path, Docs_msgp)

