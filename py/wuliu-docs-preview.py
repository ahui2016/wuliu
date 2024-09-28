import sys
import argparse
from pathlib import Path
from wuliu.albums import *
from wuliu.const import *
from wuliu.common import print_err_exit


def get_docs_metadata() -> list:
    """獲取全部可預覽檔案的屬性

    Docs_msgp 是由 "wuliu-db -dump docs" 導出的數據，檔名固定為 docs.msgp
    :return: 返回 File 列表 (參考 util/model.go 裏的 File)
    """
    return get_files_metadata(DOCS_MSGP)


def read_docs_msgp(album_path: Path) -> dict:
    return read_album_msgp(album_path, DOCS_MSGP)


def make_album(album_info: dict):
    files = get_docs_metadata()
    files = filter_files(files, album_info)
    docs = {f["ID"]: f for f in files}
    album_path = Path(WEBPAGES).joinpath(album_info["name"])
    album_path.mkdir(exist_ok=True)
    write_album_msgp(docs, album_info, album_path, DOCS_MSGP, "docs_index.html")


# ↓↓↓ main ↓↓↓

parser = argparse.ArgumentParser()

parser.add_argument("-json", type=str, help="use a json file as album info")

parser.add_argument(
    "--new-json", type=str, help="create a new json file to input album info"
)

args = parser.parse_args()

if args.new_json:
    create_new_album_info(args.new_json)
    sys.exit()

if args.json:
    album_info, err = read_album_info(args.json)
    print_err_exit(err, front_msg=f'{args.json} "name"')
    make_album(album_info)
    sys.exit()

parser.print_help()
