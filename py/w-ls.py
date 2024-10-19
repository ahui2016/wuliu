import sys
import argparse
import humanize
from wuliu.const import *
from wuliu.common import (
    print_err,
    print_err_exit,
    read_project_info,
    yaml_dump_all,
)
from wuliu.db import open_db, db_new_files, db_all_files


def trim_empty_items(files: list) -> list:
    """
    为了简化打印结果，删除内容为空的属性。
    同时也删除不需要显示的属性 (例如 checksum)
    """
    for f in files:
        del f[CHECKSUM]
        f[SIZE] = humanize.naturalsize(f[SIZE])
        del f[TYPE]
        if f[LIKE] == 0:
            del f[LIKE]
        else:
            f[LIKE] = "❤️" * f[LIKE]
        if not f[LABEL]:
            del f[LABEL]
        if not f[NOTES]:
            del f[NOTES]
        if not f[KEYWORDS]:
            del f[KEYWORDS]
        if not f[COLLECTIONS]:
            del f[COLLECTIONS]
        if not f[ALBUMS]:
            del f[ALBUMS]
        del f[CTIME]
    return files


def dump_files(files: list, title: str):
    if not files:
        print_err("空空如也")
        return
    files = trim_empty_items(files)
    text = yaml_dump_all(files)
    print(title)
    print(text)


if __name__ == "__main__":
    # 在 Windows 中使用 `>` 重定向打印到文件时可能会遇到编码问题，因此需要这行设置。
    sys.stdout.reconfigure(encoding="utf-8")  # type: ignore

    parser = argparse.ArgumentParser()

    parser.add_argument("-n", type=int, help="print N files, default N=5")
    parser.add_argument("-all", action="store_true", help="print all files")
    parser.add_argument("-orderby", type=str, help="size/like/ctime/utime")

    args = parser.parse_args()
    info = read_project_info()
    # check_not_in_backup(info)

    if args.orderby == "size":
        title = "# 【大體積檔案】"
    elif args.orderby == "like":
        title = "# 【精選檔案】"
    elif args.orderby == "ctime":
        title = "# 【最近創建檔案】"
    else:  # args.orderby == "utime":
        title = "# 【最近更新檔案】"

    if args.all:
        args.n = -1

    with open_db(Project_PY_DB) as db:
        cache = db_cache()
        files = db_get_files(cache, args.n, args.orderby)
        dump_files(files, title)
