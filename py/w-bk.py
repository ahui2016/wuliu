import sys
import argparse
from pathlib import Path
from wuliu.const import *
from wuliu.common import (
    print_err,
    read_project_info,
    check_not_in_backup,
)
# from wuliu.db import open_db


Projects = "Projects"


def print_projects_list(info: dict):
    bk_projects = info[Projects][1:]
    if len(bk_projects) == 0:
        print("無備份專案。")
        print(f"添加備份專案的方法請參閱 {Repo_URL}")
        return
    for i, project in enumerate(bk_projects):
        print(f"{i+1}. {project}")


def get_bk_path(n: int, info: dict) -> Path:
    """
    raise WuliuError
    """
    if n <= 0:
        raise WuliuError("參數 '-n' 必須大於零")

    bk_count = len(info[Projects]) - 1  # 減去源專案得到備份專案數量
    if n >= bk_count:
        raise WuliuError(f"'-n={n}', 但在目前只有 {bk_count} 個備份專案。")

    return info[Projects][n]


if __name__ == "__main__":
    # 在 Windows 中使用 `>` 重定向打印到文件时可能会遇到编码问题，因此需要这行设置。
    sys.stdout.reconfigure(encoding="utf-8")  # type: ignore

    parser = argparse.ArgumentParser()

    parser.add_argument("-projects", action="store_true", help="list all projects")
    parser.add_argument("-n", type=int, help="select a project by a number")

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    if args.projects:
        print_projects_list(info)
        sys.exit()

    if args.n is not None:
        try:
            bk_path = get_bk_path(args.n, info)
        except WuliuError as err:
            print_err(str(err))
        else:
            print(bk_path)
        sys.exit()

    parser.print_help()
