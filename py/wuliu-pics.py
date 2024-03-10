from pathlib import Path
import sys
import argparse
import msgpack


parser = argparse.ArgumentParser()
parser.add_argument('-msgp', type=str, default='',
    help='the msgp file created by "wuliu-db -dump"')
args = parser.parse_args()
print(args)

if args.msgp == '':
    print('wuliu-pics.py: error: required argument -msgp')
    sys.exit(1)

with Path(args.msgp).open(mode='rb') as f:
    files = msgpack.load(f)
    print(files)
