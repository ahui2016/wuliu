import sys
import argparse

from pathlib import Path
from tinydb import TinyDB, Query
from wuliu.const import *
from wuliu.common import print_err, print_err_exit, read_project_info, check_not_in_backup


def create_database():
	if Path(Project_PY_DB).exists():
		print_err_exit(f"file exists: {Project_PY_DB}")
	print(f"Create {Project_PY_DB}")
	db = TinyDB(Project_PY_DB)
	db.close()


# ↓↓↓ main ↓↓↓ 

parser = argparse.ArgumentParser()

parser.add_argument("--create", action="store_true",
	help="create the database")

args = parser.parse_args()

info = read_project_info()
check_not_in_backup(info)

if args.create:
	create_database()
	sys.exit()

parser.print_help()
