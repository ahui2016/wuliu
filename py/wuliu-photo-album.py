import sys
import json
import argparse
from pathlib import Path
from typing import Tuple
from wuliu.const import *
from wuliu.common import check_filename, print_err_exit


def create_new_album_info(filename: str):
    target_path = Path(filename)
    if target_path.exists():
        print(f'Error! File Exists: {filename}')
        return

    print(f'Create => {filename}')
    blob = json.dumps(New_Album_Info, indent=4)
    target_path.write_text(blob, encoding='utf8')


def read_album_info(filename: str):
    """
    :return: Tuple[dict|None, str]
    """
    data = Path(filename).read_text(encoding='utf8')
    info = json.loads(data)
    if info['name'] == '':
        return None, f'請在 {filename} 中填寫 name, 不可留空。'
    err = check_filename(info['name'])
    return info, err


# ↓↓↓ main ↓↓↓ 

parser = argparse.ArgumentParser()

parser.add_argument('-json', type=str, default='',
        help='read album info')

parser.add_argument('--new-json', type=str, default='',
        help='a filename of the json file of album info')

args = parser.parse_args()

if args.new_json != '':
    create_new_album_info(args.new_json)
    sys.exit()

if args.json != '':
    info, err = read_album_info(args.json)
    print_err_exit(err, front_msg=f'{args.json}, name: {info['name']}')
    print(f'name: {info['name']}')
