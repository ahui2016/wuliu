import sys
import json
import shutil
import argparse
from sqlite3 import Connection as Conn
from pathlib import Path
from operator import itemgetter
from jinja2 import Environment, select_autoescape, FileSystemLoader, Template
from wuliu.const import *
from wuliu.common import (
    time_now,
    file_sum512,
    type_by_filename,
    read_project_info,
    check_not_in_backup,
    json_dumps,
)
from wuliu.db import open_db, db_insert_file, db_cache


input_folder = Path(INPUT)
files_folder = Path(FILES)
buffer_folder = Path(BUFFER)
meta_folder = Path(METADATA)


def today() -> str:
    return time_now()[:10]


def daily_file_exists(file_id: str) -> bool:
    filename = file_id + ".html"
    file_path = files_folder.joinpath(filename)
    meta_path = meta_folder.joinpath(filename + ".json")
    file_exists = file_path.exists()
    meta_exists = meta_path.exists()
    if file_exists and not meta_exists:
        raise WuliuError(f"{file_path} 存在, 但 {meta_path} 不存在")
    if meta_exists and not file_exists:
        raise WuliuError(f"{meta_path} 存在, 但 {file_path} 不存在")
    return file_exists  # and meta_exists (此时两者数值相同，返回其中之一即可)


def export_file(file_id: str):
    filename = file_id + ".html"
    src = files_folder.joinpath(filename)
    dst = buffer_folder.joinpath(filename)
    if dst.exists():
        print(f"[warning] file exists: {dst}")
    else:
        print(f"Export => {dst}")
        shutil.copyfile(src, dst)

    meta_file = filename + ".json"
    src = meta_folder.joinpath(meta_file)
    dst = buffer_folder.joinpath(meta_file)
    if dst.exists():
        print(f"[warning] file exists: {dst}")
        return
    print(f"Export => {dst}")
    shutil.copyfile(src, dst)


def new_file(file_id: str) -> dict:
    filename = file_id + ".html"
    file_path = files_folder.joinpath(filename)
    f = New_File()
    f[ID] = file_id
    f[FILENAME] = filename
    now = time_now()
    f[CTIME] = now
    f[UTIME] = now
    f[CHECKSUM] = file_sum512(file_path)
    file_stat = file_path.lstat()
    f[SIZE] = file_stat.st_size
    f[TYPE] = type_by_filename(filename)
    f[COLLECTIONS] = [MY_DAILY]
    return f


def copyfile_or_not(src, dst):
    print(f"Copy => {dst}")
    if dst.exists():
        print(f"未執行複製! file exists: {dst}")
    else:
        shutil.copyfile(src, dst)


def create_daily(doc: dict, tmpl: Template, db: Conn):
    """doc = {'id': '', 'date': ''}"""
    # 剛剛檢查過檔案不存在, 因此這裡不再檢查。
    filename = doc["id"] + ".html"
    html = tmpl.render(doc=doc)
    src = files_folder.joinpath(filename)
    print(f"Create => {src}")
    src.write_text(html, encoding="utf-8")
    dst = buffer_folder.joinpath(filename)
    copyfile_or_not(src, dst)

    meta_filename = filename + ".json"
    file_meta = new_file(doc["id"])
    meta = json.dumps(file_meta, ensure_ascii=False, indent=4)
    src = meta_folder.joinpath(meta_filename)
    print(f"Create => {src}")
    src.write_text(meta, encoding="utf-8")
    dst = buffer_folder.joinpath(meta_filename)
    copyfile_or_not(src, dst)

    print("Insert into the database...")
    db_insert_file(file_meta, db)
    print("OK")


def create_export(day: str, jinja_env: Environment, db: Conn):
    file_id = DAILY_PREFIX + day
    if daily_file_exists(file_id):
        export_file(file_id)
    else:
        tmpl = jinja_env.get_template(DAILY_NEW_HTML)
        doc = dict(id=file_id, date=day)
        create_daily(doc, tmpl, db)


def get_daily_by_date(date: str, cache: dict) -> list:
    prefix = DAILY_PREFIX + date
    dates = [v for (k, v) in cache.items() if k.startswith(prefix)]
    dates.sort(key=itemgetter(ID), reverse=True)
    return dates


def get_all_daily(cache: dict) -> list:
    return get_daily_by_date("", cache)


def print_daily(files: list):
    for f in files:
        date = f[ID].removeprefix(DAILY_PREFIX)
        if f[KEYWORDS]:
            keywords = ", ".join(f[KEYWORDS])
            print(f"- {date} ({keywords})")
        else:
            print(f"- {date}")


def show_daily_list(date: str, cache: dict, webpage: bool):
    if date == "all":
        files = get_all_daily(cache)
        if not files:
            print("[warning] 沒有日記, 請創建日記", file=sys.stderr)
            return
        if webpage:
            create_index_page(files)
            return
        print("【全部日記】")
        print_daily(files)
    else:
        files = get_daily_by_date(args.list, cache)
        if not files:
            print(f"[warning] 未找到 {args.list} 的日記", file=sys.stderr)
            return
        if webpage:
            create_index_page(files)
            return
        print(f"【{args.list} 的日記】")
        print_daily(files)


def create_index_page(files: list):
    src = Path(WEBPAGES).joinpath(TEMPLATES, DAILY_INDEX_HTML)
    dst = Path(DAILY_INDEX_HTML)
    if not dst.exists():
        print(f"Copy => {dst}")
        shutil.copyfile(src, dst)

    text = json_dumps(files)
    text = "files = " + text
    daily_js = Path(DAILY_JS)
    print(f"Create => {daily_js}")
    daily_js.write_text(text, encoding="utf-8")

    print(f"請用瀏覽器打開 {dst}")


if __name__ == "__main__":
    # 在 Windows 中使用 `>` 重定向打印到文件时可能会遇到编码问题，因此需要这行设置。
    # sys.stdout.reconfigure(encoding='utf-8')  # type: ignore

    parser = argparse.ArgumentParser()

    parser.add_argument("-list", type=str, help="'-list all' or '-list 2014-09'")
    parser.add_argument("-web", action="store_true", help="'w-daily -list=all -web'")
    parser.add_argument("-edit", type=str, help="'-edit today' or '-edit 1970-12-31'")

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    jinja_env = Environment(
        loader=FileSystemLoader("webpages/templates"), autoescape=select_autoescape()
    )

    if args.list:
        db = open_db(Project_PY_DB)
        cache = db_cache(db)
        show_daily_list(args.list, cache, args.web)
        db.close()
        sys.exit()

    if args.edit == "today":
        args.edit = today()

    if args.edit:
        db = open_db(Project_PY_DB)
        create_export(args.edit, jinja_env, db)
        db.close()
        sys.exit()

    parser.print_help()
