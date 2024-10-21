import re


Repo_Name = "Wuliu File Manager"
Repo_URL = "https://github.com/ahui2016/wuliu"

Project_JSON = "project.json"
Project_PY_DB = "project.sqlite.db"
Thumb_Size = "ThumbSize"

ID = "id"
FILENAME = "filename"
CHECKSUM = "checksum"
SIZE = "size"
TYPE = "type"
LIKE = "like"
LABEL = "label"
NOTES = "notes"
KEYWORDS = "keywords"
COLLECTIONS = "collections"
ALBUMS = "albums"
CTIME = "ctime"
UTIME = "utime"

IDS = "ids"
FILENAMES = "filenames"

FILES = "files"
BUFFER = "buffer"
METADATA = "metadata"
INPUT = "input"
WEBPAGES = "webpages"
TEMPLATES = "templates"
THUMBS = "webpages/thumbs"
THUMBS_MSGP = "thumbs.msgp"
PICS_MSGP = "pics.msgp"
DOCS_MSGP = "docs.msgp"

DAILY_PREFIX = "daily-"
MY_DAILY = "my-daily"
DAILY_NEW_HTML = "daily_new.html"
DAILY_INDEX_HTML = "daily_index.html"
DAILY_JS = "daily.js"

MIME_OCTET_STREAM = "application/octet-stream"

Filename_Forbid_Pattern = re.compile(r"[^._0-9a-zA-Z\-]")
"""文件名只能使用 0-9, a-z, A-Z, _(下划线), -(短横线), .(点)。"""

Keywords_Forbid_Pattern = re.compile(r"[,\s]")
"""keywords, collections, albums 等不允许包含半角逗号和空格"""

MB = 1024 * 1024


class WuliuError(Exception):
    pass


def New_Album_Info() -> dict:
    return dict(  # 用於創建新相簿
        name="",  # 相簿名稱，必填，只允許使用 0-9, a-z, A-Z, _, -, .(点)
        ids=[],  # 通過 ID 指定圖片 (一旦指定ID, 其他條件無效)
        label="",  # label, notes, keyword, collections, albums
        notes="",  # 這五項可取併集(默認), 也可取交集, 通過下面的 union 控制。
        keywords=[],  # 這五項其中留空的，則被忽略。
        collections=[],  # 如果這五項及 ids 都留空，則表示 "全部圖片"。
        albums=[],
        union=True,  # True: 併集(聯集), False: 交集
        orderby="utime",  # 排序依據: utime/ctime/filename/like
        ascending=False,  # False: 降序, True: 昇序, 如果指定 ids, 則以 ids 為準
    )


def New_File() -> dict:
    return dict(
        id="",  # 档案名称的 CRC32
        filename="",  # 档案名称
        checksum="",  # BLAKE2b
        size=0,  # length in bytes for regular files
        type="",  # 檔案類型, 例: text/js, office/docx
        like=0,  # 點贊
        label="",  # 标签，便於搜尋
        notes="",  # 備註，便於搜尋
        keywords=[],  # 關鍵詞, 便於搜尋
        collections=[],  # 集合（分组），一个档案可属于多个集合
        albums=[],  # 相册（专辑），主要用于图片和音乐
        ctime="",  # RFC3339 檔案入庫時間
        utime="",  # RFC3339 檔案更新時間
    )


def Edit_Files_Config() -> dict:
    """用于批量修改档案属性。"""
    return dict(
        ids=[],  # 通过 ID 指定档案
        filenames=[],  # 通过档案名称指定档案
        like=0,  # 點贊
        label="",  # 标签，便於搜尋
        notes="",  # 備註，便於搜尋
        keywords=[],  # 關鍵詞, 便於搜尋
        collections=[],  # 集合（分组），一个档案可属于多个集合
        albums=[],  # 相册（专辑），主要用于图片和音乐
        learn_yaml="https://www.bairesdev.com/tools/json2yaml/",
    )
