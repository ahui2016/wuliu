import sys
import zlib
import json
import yaml
import arrow
import hashlib
import mimetypes
import msgpack
from typing import Union
from pathlib import Path
from PIL import Image, ImageOps
from .const import *


def print_err(err: str | None):
    """如果有错误就打印, 没错误就忽略."""
    if err:
        print(f"Error! {err}", file=sys.stderr)


def print_err_exit(err: str | None, front_msg: str = ""):
    """若有错误则打印并结束程序, 无错误则忽略.
    如果提供了 front_msg, 则在 err 之前显示。
    """
    if err:
        if front_msg:
            print(f"Error! {front_msg}", file=sys.stderr)
            print(err, file=sys.stderr)
            sys.exit()
        else:
            print(f"Error! {err}", file=sys.stderr)
            sys.exit()


# https://reorx.com/blog/python-yaml-tips/
class IndentDumper(yaml.Dumper):
    def increase_indent(self, flow=False, indentless=False):
        return super(IndentDumper, self).increase_indent(flow, False)


def yaml_dump(doc) -> str:
    return yaml.dump(doc, Dumper=IndentDumper, allow_unicode=True, sort_keys=False)


def yaml_dump_all(docs: list) -> str:
    return yaml.dump_all(
        docs,
        Dumper=IndentDumper,
        allow_unicode=True,
        sort_keys=False,
        explicit_start=True,
    )


def yaml_load_file(f: Path):
    text = f.read_text(encoding="utf-8")
    return yaml.safe_load(text)


# https://github.com/numpy/numpy/blob/main/numpy/core/numeric.py
def base_repr(number: int, base: int = 10, padding: int = 0) -> str:
    """
    Return a string representation of a number in the given base system.
    """
    digits = "0123456789abcdefghijklmnopqrstuvwxyz"
    if base > len(digits):
        raise ValueError("Bases greater than 36 not handled in base_repr.")
    elif base < 2:
        raise ValueError("Bases less than 2 not handled in base_repr.")

    num = abs(number)
    res = []
    while num:
        res.append(digits[num % base])
        num //= base
    if padding:
        res.append("0" * padding)
    if number < 0:
        res.append("-")
    return "".join(reversed(res or "0"))


def base36(number: int) -> str:
    return base_repr(number, 36)


def time_now() -> str:
    return arrow.now().format(arrow.FORMAT_RFC3339)


def crc32_str36(s: str) -> str:
    """把一个字符串转化为 crc32, 再转化为 36 进制。"""
    sum = zlib.crc32(s.encode())
    str36 = base36(sum)
    return str36.upper()


def name_to_id(name: str) -> str:
    """根据文件名计算出文件 ID, 确保相同的文件名拥有相同的 ID"""
    return crc32_str36(name)


# BLAKE2b is faster than MD5, SHA-1, SHA-2, and SHA-3, on 64-bit x86-64 and ARM architectures.
# https://en.wikipedia.org/wiki/BLAKE_(hash_function)#BLAKE2
# https://blog.min.io/fast-hashing-in-golang-using-blake2/
def file_sum512(name: str | Path) -> str:
    with open(name, "rb") as f:
        digest = hashlib.file_digest(f, hashlib.blake2b)
        return digest.hexdigest()


def my_type_by_filename(ext: str) -> str | None:
    if not ext:
        return None
    if ext[0] == ".":
        ext = ext[1:]
    if ext in ["doc", "docx", "ppt", "pptx", "rtf", "xls", "xlsx"]:
        return "office/" + ext
    if ext in ["epub", "mobi", "azw", "azw3", "djvu"]:
        return "ebook/" + ext
    if ext in ["zip", "rar", "7z", "gz", "tar", "bz", "bz2", "xz"]:
        return "compressed/" + ext
    if ext in [
        "md",
        "json",
        "xml",
        "html",
        "xhtml",
        "htm",
        "atom",
        "rss",
        "yaml",
        "js",
        "ts",
        "go",
        "py",
        "cs",
        "dart",
        "rb",
        "c",
        "h",
        "cpp",
        "rs",
    ]:
        return "text/" + ext
    return None


def type_by_filename(filename: str) -> str:
    ext = Path(filename).suffix.lower()
    my_type = my_type_by_filename(ext)
    if my_type is None:
        # mimetypes.add_type("application/msword", ".docx")
        # mimetypes.add_type("application/vnd.ms-excel", ".xlsx")
        # mimetypes.add_type("application/vnd.ms-powerpoint", ".pptx")
        return mimetypes.types_map[ext]
    return my_type


def read_project_info():
    project = Path(Project_JSON)
    if not project.exists():
        print_err_exit(f"not found: {project}")
    data = project.read_text(encoding="utf-8")
    info = json.loads(data)
    if info["RepoName"] != Repo_Name:
        print_err_exit(f"RepoName ({info["RepoName"]}) != {Repo_Name}")
    return info


def check_not_in_backup(info: dict):
    if info["IsBackup"]:
        print_err_exit("這是備份專案, 不可使用該功能")


def get_filenames(folder: Path) -> list[str]:
    """
    假设 folder 里全是普通档案，没有资料夹。
    """
    files = folder.glob("*")
    return [f.name for f in files]


def read_thumbs_msgp() -> dict:
    file = Path(THUMBS_MSGP)
    if not file.exists():
        return dict()
    data = file.read_bytes()
    return msgpack.unpackb(data)


def open_image(file: str | Path) -> Union[Image.Image, None]:
    try:
        img = Image.open(file)
    except OSError:
        img = None
    return img


def create_thumb_img(img: Image.Image, thumb_path: Path, thumb_size):
    """請使用 create_thumb"""
    img = ImageOps.exif_transpose(img)
    img = img.convert("RGB")
    img = ImageOps.fit(img, thumb_size)
    img.save(thumb_path)


def create_thumb(pic_path, thumb_path, thumb_size) -> str | None:
    img = open_image(pic_path)
    if img is None:
        return f"Not Image: {pic_path}"
    create_thumb_img(img, thumb_path, thumb_size)
    return None


def files_to_pics(files):
    return {file[ID]: file for file in files if file[Type].startswith("image")}


def check_filename(name: str) -> str:
    """
    :return: 有错返回 err: str, 无错返回空字符串。
    """
    if Filename_Forbid_Pattern.search(name) is None:
        return ""
    else:
        return (
            "只能使用 0-9, a-z, A-Z, _(下劃線), -(連字號), .(點)"
            "\n注意：不能使用空格，請用下劃線或連字號替代空格。"
        )


def check_keywords(keywords: list[str]) -> str | None:
    """
    :return: 有错返回 err:str, 无措返回 None 或空字符串。
    """
    joined = "".join(keywords)
    if Keywords_Forbid_Pattern.search(joined) is None:
        return None
    else:
        return "keywords/collections/albums 禁止包含半角逗號或空格"
