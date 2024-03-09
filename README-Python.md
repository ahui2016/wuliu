# Wuliu (Python)

Wuliu File Manager (五柳檔案管理腳本) Python Scripts

本文假設讀者已閱讀 <README.md>

## 支持一切程式語言

- 本軟件的本質是對 JSON 進行操作，因此使用任何語言均可為本軟件編寫功能。
- 建議使用命令 `wuliu-db -dump` (詳見 <README.md>) 把整個數據庫導出為一個 JSON

## 安裝 Python

- 本文假設讀者(用戶)有 Python 基礎
- <https://www.python.org/downloads/>
- 我在編寫本文介紹的腳本時，使用 Python 3.12

## 安裝 Wuliu Python Scripts

- 先進入原始碼中的 `py` 資料夾
- 執行 `python3 -m pip install -r requirements.txt` 或 `pip install -r requirements.txt`

然後選擇以下其中一種方法:

### 方法一

1. 複製原始碼中的 `py` 資料夾到專案根目錄
2. 用命令 `python3 py/wuliu-thumb.py` 的形式執行 Python 腳本

### 方法二

1. 參考 <README.md> 中的說明，下載原始碼，把原始碼中的 `py` 資料夾添加到系統的環境變數中。
2. 直接執行命令 `wuliu-thumb.py` 或 `wuliu-thumb`
  - Linux 系統請參考 [Executable Python Scripts](https://docs.python.org/3/tutorial/appendix.html#executable-python-scripts)


python -m pip freeze will produce a similar list of the installed packages, but the output uses the format that python -m pip install expects. 
https://docs.python.org/3/tutorial/venv.html
