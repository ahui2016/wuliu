# Wuliu (Python)

Wuliu File Manager (五柳檔案管理腳本) Python Scripts

本文假設讀者已閱讀 <README.md>

## 支持一切程式語言

本軟件的本質是對 JSON 進行操作，因此使用任何語言均可為本軟件編寫功能。

## 安裝 Python

- 本文假設讀者(用戶)有 Python 基礎
- <https://www.python.org/downloads/>
- 我在編寫以下介紹的腳本時，使用 Python 3.12

## 安裝 Wuliu Python Scripts

- 先進入原始碼中的 `py` 資料夾
- 執行 `python3 -m pip install -r requirements.txt` 或 `pip install -r requirements.txt`

然後選擇以下其中一種方法:

### 方法一

0. 創建一個虛擬環境
1. 複製原始碼中的 `py` 資料夾到專案根目錄
2. 進入專案根目錄的 `py` 資料夾，執行 `python3 -m pip install -r requirements.txt`
3. 用命令 `python3 py/wuliu-thumb.py` 的形式執行 Python 腳本

其中第 0 步是可選的，參考 https://docs.python.org/3/tutorial/venv.html

### 方法二

1. 參考 <README.md> 中的說明，下載原始碼，把原始碼中的 `py` 資料夾添加到系統的環境變數中。
2. 進入原始碼中的 `py` 資料夾，執行 `python3 -m pip install -r requirements.txt`
3. 直接執行命令 `wuliu-thumb.py` 或 `wuliu-thumb`

Linux 系統請參考 [Executable Python Scripts](https://docs.python.org/3/tutorial/appendix.html#executable-python-scripts)

## wuliu-thumbs (更新縮略圖)

- 執行命令 `wuliu-db -dump pics` 導出全部圖片的屬性 (pics.msgp)
- 執行命令 `wuliu-thumbs -msgp pics.msgp` 即可生成縮略圖。
  同時，該命令還會生成檔案 thumbs.msgp, 其中記錄了縮略圖的資訊。
- 該命令會自動對比 pics.msgp 與 thumbs.msgp, 發現新增圖片及修改過的圖片，
  沒變化的圖片會被忽略，發現已刪除的圖片也會自動刪除縮略圖。
- 縮略圖尺寸可在 project.json 中修改。
- 如果想重新生成全部縮略圖，可以刪除 thumbs.msgp
- `wuliu-thumbs --delete-orphans` 尋找並刪除漏網之魚 (應刪除但未刪除的縮略圖)
- 注意，在執行 `wuliu-thumbs --delete-orphans` 時，必須確保 thumbs.msgp 是正確的。
- 多數情況下不需要執行 `wuliu-thumbs --delete-orphans`

## wuliu-photo-album (創建相簿網頁)


New_Album_Info = {  # 用於創建新相簿
    ids: [],        # 通過 ID 指定圖片 (一旦指定ID, 其他條件無效)
    label: '',      # label, notes, keyword, collections, albums
    notes: '',      # 這五項可取交集(默認), 也可取併集, 通過下面的 union 控制。
    keywords: [],
    collections: [],
    albums: [],
    union: False,  # False: 交集, True: 併集(聯集)
    orderby: 'utime',  # 排序依據: utime/ctime/filename/like
    ascending: False,  # False: 降序, True: 昇序, 如果指定 ids, 則以 ids 為準
}


python -m pip freeze will produce a similar list of the installed packages, but the output uses the format that python -m pip install expects. 
https://docs.python.org/3/tutorial/venv.html
