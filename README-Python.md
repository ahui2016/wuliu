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

選擇以下其中一種方法:

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

我自己使用方法二。


## w-add.py (添加檔案)

- 该命令用于添加档案，同时也用于发现新档案
- 需要添加属性 `-danger` 才能真正添加新档案，否则就只是列出新档案
- 如果有一段时间未执行 `wuliu-orphan` 命令，建议先执行 `wuliu-orphan`
- 请把需要添加的档案放到 input 资料夹中，然后执行 `w-add`

### 只添加一部分新档案

- 执行 `w-add --new-yaml add.yaml`
  可在 input 资料夹中生成一个新的 add.yaml, 方便编辑
- 在 add.yaml 中会列出全部待添加的档案名称
- 在 add.yaml 中删除不需要添加的档案名称后，执行命令
  `w-add -yaml add.yaml` 即可只添加指定的新档案
- 如果 add.yaml 檔案已存在, 可添加 '-danger' 參數強制覆蓋:
  `w-add --new-yaml add.yaml -danger`

### 批量修改待添加档案的属性

- 执行 `w-add --new-yaml add.yaml`
  可在 input 资料夹中生成一个新的 add.yaml, 方便编辑
- 执行 `w-add -yaml add.yaml` 列出 add.yaml 中指定的待添加档案，
  同时列出 add.yaml 里的档案属性，该属性将应用于全部待添加档案。
- 注意, add.yaml 应放在专案的根目录。
- 需要添加属性 `-danger` 才能真正添加新档案，否则就只是列印相关信息

### 小技巧 (很實用)

- 生成 add.yaml 后，可删除其中的 filenames 的内容 (修改后是这样 `filenames: []`),
  表示選中 input 资料夹中的全部档案。
- 删除 filenames 的内容后，可执行 `w-add -yaml add.yaml` 预览配置，
  添加参数 `-danger` 正式执行。

### 添加後，修改檔案及其屬性

一旦成功添加檔案，在 metadata 資料夾中會自動生成同名 json, 在該 json 中
含有檔案屬性，但請勿直接修改 json 內容。

也不可直接修改檔案本身。

如需修改檔案本身 或 檔案屬性，請使用 wuliu-export 與 wuliu-overwrite
(詳見後文的相關章節)

如需更改檔案名稱，請使用 wuliu-rename 命令。


## 档案属性

```
{
    id="",           # 档案名称的 CRC32
    filename="",     # 档案名称
    checksum="",     # BLAKE2b
    size=0,          # length in bytes for regular files
    type="",         # 檔案類型, 例: text/js, office/docx
    like=0,          # 點贊
    label="",        # 标签，便於搜尋
    notes="",        # 備註，便於搜尋
    keywords=[],     # 關鍵詞, 便於搜尋
    collections=[],  # 集合（分组），一个档案可属于多个集合
    albums=[],       # 相册（专辑），主要用于图片和音乐
    ctime="",        # RFC3339 檔案入庫時間
    utime=""         # RFC3339 檔案更新時間
}
```

- ID 是档案名称的 CRC32，有冲突的可能性，但可能性较低，
  大不了冲突了再改档案名称，问题不大。
  后续如果档案数量大了，可以考虑改用 CRC64
- 关于 CRC32 <https://softwareengineering.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed>
- Type, Label, Note, Keywords 等都是为了方便搜寻，请大胆灵活使用。
- 为了避免歧义, keywords, collections, albums 等不允许包含半角逗号和空格。
- 請勿直接修改 metadata 裏的檔案。
  如需修改，請導出後修改，然後再使用 w-overwrite 覆蓋舊檔案。
- 手動修改檔案屬性時，請勿直接修改 ID, Filename, Checksum, Size.
- ID 與 Filename 是相關的，修改檔案名稱會改變 ID.
  如需更改檔案名稱，請使用 w-rename 命令。


## w-db

采用 TinyDB <https://github.com/msiemens/tinydb>

- `w-db --create` 創建數據庫
- `w-db --update=add` 發現 metadata 資料夾中的新檔案並添加到數據庫中。
  注意: 該命令只處理新檔案, 忽略檔案被刪除或修改的情況, 如須確保數據庫完全更新, 請刪除數據庫後重新建立。


## w-ls (列印數據庫中的檔案)

- `w-ls` 默認列出 5 個最近更新的檔案
- `w-ls -n 10` 列出 10 個最近更新的檔案
- `w-ls -orderby ctime` 列出 5 個最近創建的檔案
- `w-ls -orderby size -n=3` 列出 3 個體積最大的檔案
- `w-ls -all > utime.yaml` 列出全部檔案, 默認按更新日期排序 (保存到utime.yaml中)
- `w-ls -h` 查看幫助


## w-daily (日記/日報)

- 不使用 markdown, 使用 html, 使用 emmet❗
- 參考 <https://docs.emmet.io/>
- 檔案 ID 特殊處理 (daily-2024-10-01)
- 檔案名: daily-2024-10-01.html
- collections: `["my-daily"]`
- `w-daily -edit today` 新建/編輯今天的日記。
- `w-daily -edit 2024-10-01` 新建/編輯指定日期的日記。
- 使用 `w-daily -edit` 時會檢查 files 資料夾中是否存在同名檔案。
  - 如果不存在, 則新建, 並導出到 buffer 資料夾 (metadata 檔案也一並導出)
  - 如果已存在, 則導出到 buffer 資料夾 (metadata 檔案也一並導出)
  - 編輯 metadata 時請勿修改 labal
  - 編輯 metadata 時, 建議只編輯 keywords
  - 導出到 buffer 資料夾時, 也會檢查是否存在同名檔案。
- 編輯 buffer 資料夾中的檔案後, 使用 `w-overwrite` 命令更新檔案。
- `w-daily -list 2024-10` 列出 2024 年 10 月已創建的日記
- `w-daily -list 2024` 列出 2024 年內已創建的日記
- `w-daily -list all` 列出全部已創建的日記
- `w-daily -list=all -web` 生成 daily-index.html
- `w-daily` 功能以檔案 ID 為準, ID 以 `daily-` 開頭的即視為日記。
- `w-daily --create-website` 把 collections 包含 "my-daily" 的檔案全部複製到
  webpages 資料夾中, 創建一個簡單的網站。 因此, 日記中引用的圖片等資源請加入 my-daily 集合。


## w-overwrite

- 執行 `w-overwrite` 查看待覆蓋檔案清單。
  (注意，待覆蓋檔案應存放在 buffer 資料夾中。)
- 在該清單中可以看到，凡是 *非json* 檔案都將覆蓋 files,
  凡是 *json* 檔案都將覆蓋 metadata.
- 如果其中有 json 檔案想覆蓋 files, 請執行 `w-overwrite --new-yaml overwrite.yaml`
  然後編輯 overwrite.yaml, 根據需要把其中的 "metadata" 改為 "files".
  (此時，還可以刪除 overwrite.json 裏的一部分檔案名稱，只有保留在清單中的檔案纔會被覆蓋。)
- 經過上述操作後，執行 `w-overwrite -yaml overwrite.yaml` 查看待覆蓋檔案清單
- 執行 `w-overwrite -yaml overwrite.yaml -danger` 或
  `w-overwrite -danger` 正式覆蓋。
- 如果不使用 `-danger` 參數，則只是查看待覆蓋檔案清單，不會真正發生覆蓋。
- 【注意】手動修改檔案屬性時，請勿直接修改 ID, Filename, Checksum, Size, Type, UTime.
- 請勿直接修改 metadata 資料夾中的檔案。
  如需修改，請導出後修改，再使用 w-overwrite 進行更新。
  另外，可以使用 w-metadata 命令批量修改属性。
- ID 與 Filename 是相關的，請勿直接修改檔案名稱，如需更改檔案名稱，
  請使用 w-rename 命令。


## wuliu-photo-album (創建相簿網頁)

- 執行命令 `wuliu-db -dump pics` 導出全部圖片的屬性 (pics.msgp)
- 執行命令 `wuliu-photo-album --new-json photo-album.json` 生成 photo-album.json
- photo-album.json 的內容如下所示:

```python
New_Album_Info = {  # 用於創建新相簿
    name: '',       # 相簿名稱，必填，只允許使用 0-9, a-z, A-Z, _, -, .(点)
    ids: [],        # 通過 ID 指定圖片 (一旦指定ID, 其他條件無效)
    label: '',      # label, notes, keyword, collections, albums
    notes: '',      # 這五項可取併集(默認), 也可取交集, 通過下面的 union 控制。
    keywords: [],   # 這五項其中留空的，則被忽略。
    collections: [],  # 如果這五項及 ids 都留空，則表示 "全部圖片"。
    albums: [],
    union: True,       # True: 併集(聯集), False: 交集
    orderby: 'utime',  # 排序依據: utime/ctime/filename/like
    ascending: False,  # False: 降序, True: 昇序, 如果指定 ids, 則以 ids 為準
}
```

- 請編輯 photo-album.json, 其中 **name** 必須填寫，用來作為相簿資料夾名稱。
- **ids** 通常不填寫，一旦填寫, label, notes, keywords, collections, albums 和
  union, orderby, ascending 都無效。【注意，該功能暫未實現】
- **label, notes, keywords, collections, albums** 這五項默認取併集(聯集),
  (取交集的功能暫時不做),
  如果這五項及 ids 都留空，則生成一個包含全部圖片的相簿。
- 命令 `wuliu-photo-album -json photo-album.json` 生成相簿。
  該命令會自動讀取 pics.msgp, 生成的相簿在 webpages 資料夾中。
- 其中排序功能是在前端實現的，因此生成相簿後如果想改變排序，可進入相簿資料夾修改
  pics.js 中的 orderby 和 ascending, 保存後刷新頁面即可生效。
- 縮略圖尺寸可在 project.json 中修改。
- 相簿的大標題和副標題在相簿資料夾的 index.html 中修改。

## wuliu-docs-preview (文檔預覽網頁)

該腳本可創建一個網頁，便於預覽文檔 (pdf/html/txt 等瀏覽器可直接預覽的格式)。

- 執行命令 `wuliu-db -dump docs` 導出全部可預覽檔案的屬性 (docs.msgp)
- 執行命令 `wuliu-docs-preview --new-json docs-album.json` 生成 docs-album.json
- docs-album.json 的內容與前述 wuliu-photo-album (創建相簿網頁) 的
  photo-album.json 相同，填寫方法也相同。
- 命令 `wuliu-docs-preview -json docs-album.json` 生成網頁。
  該命令會自動讀取 docs.msgp, 生成的網頁在 webpages 資料夾中。


## TODO

- 重建數據庫時檢查孤立檔案
- w-checksum.py

## notes

- https://github.com/wintermute-cell/magick.css
- https://pypi.org/project/PyYAML/
- https://www.bairesdev.com/tools/json2yaml/
- https://www.cloudbees.com/blog/yaml-tutorial-everything-you-need-get-started
- https://reorx.com/blog/python-yaml-tips/ （有用！）


python -m pip freeze will produce a similar list of the installed packages, but the output uses the format that python -m pip install expects. 
https://docs.python.org/3/tutorial/venv.html
