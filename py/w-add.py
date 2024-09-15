import sys
import tomllib
import tomli_w
import argparse
from pathlib import Path

from wuliu.const import *
from wuliu.common import print_err, print_err_exit, read_project_info, \
    check_not_in_backup, get_filenames, time_now, name_to_id, file_sum512, \
    type_by_filename


def new_file(name: str) -> dict:
    f = New_File()
    f[ID] = name_to_id(name)
    f[FILENAME] = name
    now = time_now()
    f[CTIME] = now
    f[UTIME] = now
    return f


def new_files_from(names: list[str], folder: str) -> list:
    """
    把档案名 (names) 转换为 files, 此时假设档案在 folder 资料夹内。
    """
    files = []
    folder_path = Path(folder)
    for name in names:
        file_path = folder_path.joinpath(name)
        if file_path.is_dir():
            print(f"{file_path} 是資料夾, 自動忽略")
            continue
        file_stat = file_path.lstat()
        checksum = file_sum512(file_path)
        f = new_file(name)
        f[CHECKSUM] = checksum
        f[SIZE] = file_stat.st_size
        f[TYPE] = type_by_filename(name)
        files.append(f)
    return files


def find_new_files() -> list:
    """
    寻找 input 资料夹里的全部档案。
    """
    names_in_input = get_filenames(Path(INPUT))
    files = new_files_from(names_in_input, INPUT)
    return files


def create_config_toml(filename:str, allow_danger:bool):
    file_path = Path(filename)
    if file_path.exists() and not allow_danger:
        print_err(f"file exists: {filename}")
        return
    print(f"Create => {filename}")
    text = tomli_w.dumps(Edit_Files_Config())
    file_path.write_text(text, encoding="utf-8")


if __name__ == "__main__":

    parser = argparse.ArgumentParser()

    parser.add_argument("-danger", action="store_true",
        help="allow dangerous operations")

    parser.add_argument("--new-toml", type=str,
        help="create a TOML file for adding files")

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    if args.new_toml:
        create_config_toml(args.new_toml, args.danger)
        sys.exit()

    parser.print_help()
