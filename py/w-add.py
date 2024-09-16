import sys
import yaml
import argparse
from pathlib import Path

from wuliu.const import *
from wuliu.common import print_err, print_err_exit, read_project_info, \
    check_not_in_backup, get_filenames, time_now, name_to_id, file_sum512, \
    type_by_filename, yaml_dump, yaml_load_file


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


def find_input_files() -> list:
    """
    寻找 input 资料夹里的全部档案。
    """
    names_in_input = get_filenames(Path(INPUT))
    files = new_files_from(names_in_input, INPUT)
    return files


def add_files_config(ids: list, filenames: list) -> dict:
    cfg = Edit_Files_Config()
    cfg[IDS] = ids
    cfg[FILENAMES] = filenames
    return cfg


def create_config_yaml(filename:str, allow_danger:bool):
    file_path = Path(filename)
    if file_path.exists() and not allow_danger:
        print_err(f"file exists: {filename}")
        return
    names_in_input = get_filenames(Path(INPUT))
    text = yaml_dump(add_files_config([], names_in_input))
    print(f"Create => {filename}")
    file_path.write_text(text, encoding="utf-8")


def read_config(file_path: Path) -> dict:
    """返回 None 表示有错误"""
    cfg = yaml_load_file(file_path)
    if len(cfg[IDS]) > 0:
        print_err_exit("添加新檔案時不可通過 ID 指定檔案")
    return cfg


if __name__ == "__main__":

    parser = argparse.ArgumentParser()

    parser.add_argument("-danger", action="store_true",
        help="allow dangerous operations")

    parser.add_argument("--new-yaml", type=str,
        help="create a YAML file for adding files")

    parser.add_argument('-yaml', type=str,
        help='use a YAML file to add files')

    args = parser.parse_args()
    info = read_project_info()
    check_not_in_backup(info)

    if args.new_yaml:
        create_config_yaml(args.new_yaml, args.danger)
        sys.exit()

    if args.yaml:
        cfg = read_config(Path(args.yaml))
        print(cfg)
        sys.exit()

    parser.print_help()
