# Wuliu (Golang)

Wuliu File Manager (五柳檔案管理腳本) Golang Scripts

## 名词

- 档案(檔案) = 文件 = file
- 资料夹(資料夾) = 文件夹 = folder = directory
- 专案(專案) = 项目 = project
- 列印 = 打印/显示 = print
- 後綴名(副檔名) = 后缀名/扩展名

## 主要功能

- 給檔案增加更多屬性，例如: 備注、標籤、關鍵詞、点赞。
  有了这些属性，就能做到：
  - 快速找出大文件
  - 快速找出指定日期的文件
  - 快速找出常用或精品文件（利用点赞功能）
  - 根据标签或关键词搜索文件
  - 随手给文件写一句备注或说明（利用备注属性）
  - 多维度分类文件（利用关键词、集合、专辑等属性）
- 如果你懂编程（任何语言！不限语言！）你就能自行添加属性，
  自行编程添加功能，因为每个文件的属性都是一个 json 文件。
- 檢查檔案完整性
- 方便地備份
- 方便地修復受損檔案

## 编程简单第一

- 非常重视编程简单
- 运行效率和使用方便都是次要的

**【注意】**: 
這些腳本為了編程方便，犧牲了易用性，因此在使用過程中必須
一邊閱讀本文 (README.md) 一邊使用。

## 安裝

- 安裝 Golang
- 下載 Wuliu 原始碼
- 添加環境變數

## Scripts

- wuliu-init (新建/初始化一个专案)
- `cat project.info` (查看专案信息)
- wuliu-orphan (检查有无孤立档案)
- wuliu-add (添加档案)
- wuliu-delete (删除档案)
- wuliu-rename (更改檔案名稱)
- wuliu-list (列印档案、标签、备注、关键词等)
- wuliu-db (数据库信息，更新缓存)
- wuliu-checksum (检查档案完整性)
- wuliu-backup (备份专案)
- wuliu-export (導出檔案或檔案屬性)
- wuliu-overwrite (更新單個檔案或檔案屬性)
- wuliu-metadata (批量修改多個檔案的屬性)
- wuliu-like (點讚，方便尋找精品或常用檔案)

## wuliu-init

- 新建/初始化一个专案，主要是新建一些资料夹和数据库、配置等。
- 只能对一个空资料夹进行初始化
- 使用方法: 进入一个空资料夹，执行 `wuliu-init -name [NAME]` 进行初始化。
- 注意，请为不同的专案设定不同的名称，备份时有用。
- 备份专案（详见关于 `wuliu-backup` 的说明）的专案名称必须与源专案一致。
- `wuliu-init -h` 列印帮助信息
- `wuliu-init -v` 列印版本信息
- `wuliu-init -where` 列印 wuliu-init 的位置

### 资料夹

初始化后，会得到一些资料夹。

- **input** (專用於添加新檔案)
- **files** (執行 wuliu-add 後, input 裏的檔案會移動到 files 裏)
- **metadata** (files 裏的每個檔案都有一個對應的屬性檔案 (json檔案))
- **buffer** (用於導出檔案或修改檔案)
- **webpages** (生成網頁便於檢索檔案)
- **recyclebin** (執行 wuliu-delete 刪除的檔案會被移進這裏)

其中，尤其需要注意 input 與 buffer 的區別，
一個是專用於添加新檔案，一個是用於更新檔案（或修改檔案屬性）。

並且, wuliu-add 和 wuliu-import 命令只能操作 input 資料夾,
wuliu-export 和 wuliu-overwrite 命令只能操作 buffer 資料夾。

**【注意】**:
請勿直接修改 files 與 metadata 裏的檔案。
如需修改，請導出後修改，然後再使用 wuliu-overwrite 覆蓋舊檔案。

## project.json

执行 `cat project.json` 可查看专案信息。
当然，也可直接打开 project.json 查看。

```
type ProjectInfo struct {
    RepoName        string   // 用于判断资料夹是否 Wuliu 专案
    ProjectName     string   // 备份时要求专案名称相同
    IsBackup        bool     // 是否副本（副本禁止添加、删除等）
    Projects        []string // 第一个是主专案，然后是备份专案
    LastBackupAt    []string // 上次备份时间
    CheckInterval   int      // 检查完整性, 单位: day
    CheckSizeLimit  int      // 检查完整性, 单位: MB
    ExportSizeLimit int      // 導出檔案體積上限，單位: MB
    ThumbSize       [2]int   // 縮略圖尺寸
}
```

其中 `Projects` 要注意必须确保第一个是 "."。

## wuliu-orphan

- 孤立档案: files 资料夹里的每个档案都有一个同名 json 在 metadata 资料夹中，
  如果缺少同名 json, 就是孤立档案。
- `wuliu-orphan --check` 命令会检查有无孤立档案，如果有，就提示处理
- 如果在 metadata 中有 json, 但在 files 中找不到对应的档案，也会提示处理
- 建议在某些操作（例如添加档案）之前先检查有无孤立档案

## wuliu-add

- 该命令用于添加档案，同时也用于发现新档案
- 需要添加属性 `--danger` 才能真正添加新档案，否则就只是列出新档案
- 如果有一段时间未执行 `wuliu-orphan` 命令，建议先执行 `wuliu-orphan`
- 请把需要添加的档案放到 input 资料夹中，然后执行 `wuliu-add`

### 只添加一部分新档案

- 执行 `wuliu-add --newjson add.json`
  可在 input 资料夹中生成一个新的 add.json, 方便编辑
- 在 add.json 中会列出全部待添加的档案名称
- 在 add.json 中删除不需要添加的档案名称
- 执行命令 `wuliu-add --json add.json` 只添加指定的新档案

### 批量修改待添加档案的属性

- 执行 `wuliu-add --newjson add.json`
  可在 input 资料夹中生成一个新的 add.json, 方便编辑
- 执行 `wuliu-add --json add.json` 列出 add.json 中指定的待添加档案，
  同时列出 add.json 里的档案属性，该属性将应用于全部待添加档案。
- 注意, add.json 应放在专案的根目录。
- 需要添加属性 `--danger` 才能真正添加新档案，否则就只是列印相关信息

### 档案属性

```
{
    ID          string    `json:"id"`          // 档案名称的 CRC32
    Filename    string    `json:"filename"`    // 档案名称
    Checksum    string    `json:"checksum"`    // BLAKE2b
    Size        int64     `json:"size"`        // length in bytes for regular files
    Type        string    `json:"type"`        // 檔案類型, 例: text/js, office/docx
    Like        int64     `json:"like"`        // 點贊
    Label       string    `json:"label"`       // 标签，便於搜尋
    Notes       string    `json:"notes"`       // 備註，便於搜尋
    Keywords    []string  `json:"keywords"`    // 關鍵詞, 便於搜尋
    Collections []string  `json:"collections"` // 集合（分组），一个档案可属于多个集合
    Albums      []string  `json:"albums"`      // 相册（专辑），主要用于图片和音乐
    CTime       string    `json:"ctime"`       // RFC3339 檔案入庫時間
    UTime       string    `json:"utime"`       // RFC3339 檔案更新時間
    Checked     string    `json:"checked"`     // RFC3339 上次校驗檔案完整性的時間
    Damaged     bool      `json:"damaged"`     // 上次校驗結果 (檔案是否損壞)
}
```

- ID 是档案名称的 CRC32，有冲突的可能性，但可能性较低，
  大不了冲突了再改档案名称，问题不大。
  后续如果档案数量大了，可以考虑改用 CRC64
- 关于 CRC32 <https://softwareengineering.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed>
- Type, Label, Note, Keywords 等都是为了方便搜寻，请大胆灵活使用。
- Keywords, Collections 等 `[]string` 类型，都排序，排序后转为纯字符
  （用逗号空格 `, ` 分隔）方便保存到 kv 数据库。
- 因此 `[]string` 类型在用户输入时不允许包含逗号、顿号和空格。
- 請勿直接修改 metadata 裏的檔案。
  如需修改，請導出後修改，然後再使用 wuliu-overwrite 覆蓋舊檔案。
- 手動修改檔案屬性時，請勿直接修改 ID, Filename, Checksum, Size.
- ID 與 Filename 是相關的，修改檔案名稱會改變 ID.
  如需更改檔案名稱，請使用 wuliu-rename 命令。

### 添加後修改檔案及其屬性

一旦成功添加檔案，在 metadata 資料夾中會自動生成同名 json, 在該 json 中
含有檔案屬性，但請勿直接修改 json 內容。

也不可直接修改檔案本身。

如需修改檔案本身 或 檔案屬性，請使用 wuliu-export 與 wuliu-overwrite
(詳見後文的相關章節)

如需更改檔案名稱，請使用 wuliu-rename 命令。

## wuliu-delete

- 该命令删除添加档案，包括删除对应的 json 档案和数据库中的条目
- `wuliu-delete -id [ID]` 通过档案 ID 指定需要删除的档案（只能指定一个）
- `wuliu-delete -name [NAME]` 通过档案名称指定需要删除的档案（只能指定一个）
- `wuliu-delete --newjson delete.json` 在专案根目录生成一个 delete.json 档案模板，
  方便批量填写需要删除的档案。
- 在 delete.json 中填写要删除的一个或多个档案的 id
- `wuliu-delete --json delete.json ` 通过 delete.json 指定需要删除的档案（可指定多个）
- 需要添加属性 `--danger` 才能真正删除档案，否则就只是列出 delete.json 的内容

## wuliu-rename

ID 與 Filename 是相關的，修改檔案名稱會改變 ID.
而且, files 與 metadata 的檔案名稱也要同時更改，因此，
如需更改檔案名稱，請使用 wuliu-rename 命令，不要手動更改。

- `wuliu-rename -id=[ID] -name [NAME]` 其中 ID 是舊ID, NAME 是新檔名。
- 注意檔名包括後綴名。
- 更改檔名不會修改 UTime(檔案更新時間)

## wuliu-list

- `wuliu-list` 列印最近 15 个档案 (ID, 体积, 档案名称)
- `wuliu-list n=100` 列印最近 100 个档案，按 CTime 倒序排列 (CTime 是入库时间)
- 默認按 CTime 排序，使用參數 `-orderby [INDEX]` 可按其他維度排序
  (例如 size, like, utime 等)
  (注意，有時需要先執行 `wuliu-db -update=cache` 更新數據庫緩存。)
  - 例: `wuliu-list -orderby utime` 列印最近修改過的 15 个档案
- 默認從大到小排序 (descending), 使用參數 `-asc` 改為從小到大排序 (ascending)。
  - 例: `wuliu-list -orderby size` 列印體積最大的 15 个档案
  - 例: `wuliu-list -orderby=size -asc` 列印體積最小的 15 个档案
- 默認列印簡單信息 (ID, 体积, 档案名称), 使用參數 `-more` 列印詳細信息。
- `wuliu-list > list.txt` 可把結果保存到一個檔案中。

上面是 wuliu-list 列印檔案的功能，另外, wuliu-list 還有其他功能，如下所示:
(注意，有時需要先執行 `wuliu-db -update=cache` 更新數據庫緩存。)

- `wuliu-list -labels` 列印全部標籤
- `wuliu-list -notes` 列印全部備註
- `wuliu-list -keywords` 列印全部關鍵詞
- `wuliu-list -collections` 列印全部集合
- `wuliu-list -albums` 列印全部相冊(專輯)

建議使用 `wuliu-list -labels > labels.txt` 的方式把結果保存到一個檔案中。

## wuliu-search

- `wuliu-search -keyword 小米` 搜尋 keyword 為 "小米" 的檔案。
- 可通過 filename/notes/label/keyword/collection/album 搜尋檔案。
- 其中 keyword/collection/album 默認精確匹配 (-match=exactly)
- 其中 filename/notes/label 默認前綴匹配 (-match=prefix)
- 匹配方式可選擇 exactly/prefix/contains/suffix
- 例如 `wuliu-search -match=contains -filename 偵探` 搜尋檔名包含 "偵探" 的檔案。
- 搜尋結果 (檔案清單) 默認按檔案創建時間排序 (-orderby=ctime)
- 排序方式可選擇 ctime/utime/filename
- 默認從大到小排序 (descending), 使用參數 `-asc` 改為從小到大排序 (ascending)。
- 例如

## 数据库 (bolt)

- https://github.com/etcd-io/bbolt
- Please note that Bolt obtains a file lock on the data file so multiple processes cannot
  open the same database at the same time. 
- If the key doesn't exist then it will return nil. It's important to note that you can have a
  zero-length value set to a key which is different than the key not existing.
- Please note that values returned from `Get()` are only valid while the transaction is open.
  If you need to use a value outside of the transaction then you must use `copy()` to copy
  it to another byte slice.
- Note that, while RFC3339 is sortable, the Golang implementation of RFC3339Nano does
  not use a fixed number of digits after the decimal point and is therefore not sortable.

## wuliu-db

- `wuliu-db --info=count` 查看数据库条目数量
- `wuliu-db --info=size` 查看全部档案的总体积
- `wuliu-db -dump all` 導出整個數據庫到一個 msgpack 格式的檔案
- `wuliu-db -dump pics` 導出圖片到 msgpack
- `wuliu-db -dump docs` 導出用瀏覽器能直接預覽的文檔到 msgpack

### 更新数据库

- 更新数据库，是指以 metadata 为准更新数据库，因此如果一段时间没执行 wuliu-orphan,
  请先执行一次 wuliu-orphan 再更新数据库。
- 执行 `wuliu-db --update=rebuild` 根据 metadata(真实的 json 档案) 重建整个数据库。
  执行 `wuliu-db --update=cache` 根据缓存更新索引（不需要读取硬盘里的 json 档案）。
- 由于数据库缓存（即 files 索引和 filename 索引）在添加文件、修改文件属性、删除文件时
  会自动更新，因此多数情况下只需要 `--update=cache`, 不需要重建数据库。

### keyword/collection/album 改名

例如 `wuliu-db -keyword 叮噹貓 --rename-to 多啦A梦` 把數據庫裡名為 "叮噹貓"
的 keyword 批量改為 "多啦A梦"。

collection 或 album 的更改也類似，例如 `wuliu-db -album 歌曲 --rename-to 音樂`
把數據庫裡名為 "歌曲" 的 album 批量改為 "音樂"。

## wuliu-checksum

- `wuliu-checksum --renew` 将全部文件的 damaged 设为 false, 上次检查时间设为 epoch
- `wuliu-checksum --check` 校验文件完整性（看文件是否损坏）
- `wuliu-checksum --projects` 列印全部专案
- `wuliu-checksum -n [N]` 通过序号选择专案，默认是 0 (即当前专案)

在执行 `wuliu-checksum` 命令时，有时会显示以下信息：

```
已選擇專案: .
數據庫檔案數量: 5
待檢查檔案數量: 5
```

如果发现 "數據庫檔案數量" 与 "待檢查檔案數量" 不一致，
可以执行 `wuliu-checksum --renew` 进行修正。

执行 `wuliu-checksum --check` 时，会根据 project.json 中的 CheckInterval (检查周期)
自动判断档案是否需要检查，根据 CheckSizeLimit (检查体积上限) 自动终止检查，防止
单次检查时间太长。

### wuliu-checksum -same (找出重複檔案)

`wuliu-checksum -same` 該命令不可與參數 `-n` 同時使用，只能檢查當前專案。

## wuliu-backup

创建新备份专案的方法：
（差点就为这个功能写程序了，写到一半发现违反“编程简单第一”原则，就改成手动操作。）

- 请先进入一个空资料夹（新的备份专案的根目录，以下称为 backupRoot）
- 把现有专案（以下称为“主专案”）的 files, metadata 这两个资料夹，
  以及 project.db, project.json 这两个档案复制到 backupRoot 中。
  其他资料夹和档案可复制可不复制。
- 编辑 backupRoot 里的 project.json, 把 IsBackup 的值改为 true
- 编辑主专案中的 project.json, 把 backupRoot 的路径添加到 Projects 清單中。
  注意，路径里的反斜杠改要为 "`\\`" 或 "`/`"
- 编辑主专案中的 project.json, 在 LastBackupAt 中添加一个时间
  （可以是空字符串 ""）

以上是创建新备份专案的方法，以下是 wuliu-backup 的其他命令：

- `wuliu-backup --projects` 列印全部备份专案（目标专案）
- “目标专案”是指专门用于备份的专案
- 本软件采用单向备份方式，备份时以“源专案”为准，
  使目标专案里的档案变成与源专案一样。
- `wuliu-backup -n [N]` 通过序号选择目标专案
- `wuliu-backup  -n [N] -danger` 正式执行备份
- 例如执行命令 `wuliu-backup -n 1` 会列印第 1 个备份专案的信息，但不会执行备份。
  而执行命令  `wuliu-backup -n=1 -danger` 则会正式执行备份。
- 建议在执行 `wuliu-backup -n [N]` 查看信息前，先执行 `wuliu-db -update=cache`
- 有时还可能需要去目标专案的根目录里执行  `wuliu-db -update=cache`
- 当档案数量较少时，建议先在源专案与目标专案两边都执行
  `wuliu-orphan --check` 和 `wuliu-db -update=rebuild`
  因为备份时需要使用数据库，而重建数据库有助于确保数据库与实际档案信息保持一致。
- 備份成功後，會自動更新數據庫。

### 修復受損檔案

- 如果發現受損檔案，可使用 `wuliu-backup -fix` 命令嘗試自動修復。
- 使用該命令時，需要同時使用參數 `-n` 指定目標專案。
- 如果自動修復失敗（通常因為兩邊專案裏的同一檔案都受損了），
  可換一個目標專案再嘗試修復。
- 如果仍無法修復，則需要手動修復。

手動修復方法如下：

- 方法一：使用 wuliu-export 命令導出受損檔案，然後刪除受損檔案。
- 方法二：使用 wuliu-overwrite 覆蓋受損檔案。

wuliu-export 與 wuliu-overwrite 的使用方法詳見本文的其他章節。

## wuliu-export

- `wuliu-export -file [ID]` 通過檔案 ID 指定要導出的檔案
- `wuliu-export -meta [ID]` 通過檔案 ID 指定要導出的檔案屬性 (json)
- `wuliu-export -id [ID]` 通過檔案 ID 導出的一個檔案及其屬性
- `wuliu-export -batch [FILENAME]` 通過一個 json 檔案進行批量導出。
  (批量導出功能暫時不做，因為預估該功能需求不大)
  【小技巧】手动把 metadata 里的 json 复制到 buffer 里，相当于批量导出。
- **注意** 默認只能導出 300MB 以下的檔案 (單檔案體積上限)。
  - 修改 project.json 中的 ExportSizeLimit 可更改該限制 (單位:MB)
- 如需導出大體積檔案，建議手動複製。

被導出的檔案一律導出到 buffer 資料夾中。

## wuliu-overwrite

- 執行 `wuliu-overwrite` 查看待覆蓋檔案清單。
  (注意，待覆蓋檔案應存放在 buffer 資料夾中。)
- 在該清單中可以看到，凡是 *非json* 檔案都將覆蓋 files,
  凡是 *json* 檔案都將覆蓋 metadata.
- 如果其中有 json 檔案想覆蓋 files, 請執行 `wuliu-overwrite -newjson overwrite.json`
  然後編輯 overwrite.json, 根據需要把其中的 "metadata" 改為 "files".
  (此時，還可以刪除 overwrite.json 裏的一部分檔案名稱，只有保留在清單中的檔案纔會被覆蓋。)
- 經過上述操作後，執行 `wuliu-overwrite -json overwrite.json` 查看待覆蓋檔案清單
- 執行 `wuliu-overwrite -json overwrite.json -danger` 或
  `wuliu-overwrite -danger` 正式覆蓋。
- 如果不使用 `-danger` 參數，則只是查看待覆蓋檔案清單，不會真正發生覆蓋。
- 【注意】手動修改檔案屬性時，請勿直接修改 ID, Filename, Checksum, Size, Type, UTime.
- 請勿直接修改 metadata 裏的檔案。
  如需修改，請導出後修改，然後再使用 wuliu-overwrite 覆蓋舊檔案。
  另外，可以使用 wuliu-metadata 命令批量修改属性。
  如果 wuliu-metadata 也无法满足要求，可以直接修改 metadata 资料夹里的 json 档案，
  然后执行 `wuliu-db -update=rebuild` 重建数据库。
- 【注意】如果进入 metadata 资料夹直接修改 json, 不会自动更新 UTime (可手动修改)。
- 【注意】进入 metadata 资料夹修改 json 后请立即重建数据库。
- 【小技巧】也可以把 metadata 里的 json 复制到 buffer 里，修改后执行 wuliu-overwrite
- ID 與 Filename 是相關的，修改檔案名稱會改變 ID.
  如需更改檔案名稱，請使用 wuliu-rename 命令。

## wuliu-metadata

該命令用於批量修改多個檔案的屬性。

執行 `wuliu-metadata --newjson metadata.json` 可生成檔案 metadata.json, 其結構如下:

```
type EditFiles struct {
    IDs         []string `json:"ids"`         // 通过 ID 指定档案
    Filenames   []string `json:"filenames"`   // 通过档案名称指定档案
    Like        int64    `json:"like"`        // 點贊
    Label       string   `json:"label"`       // 标签，便於搜尋
    Notes       string   `json:"notes"`       // 備註，便於搜尋
    Keywords    []string `json:"keywords"`    // 關鍵詞, 便於搜尋
    Collections []string `json:"collections"` // 集合（分组），一个档案可属于多个集合
    Albums      []string `json:"albums"`      // 相册（专辑），主要用于图片和音乐
}
```

請在 metadata.json 中填寫內容，其中，通過 ids 指定要修改的檔案 ID,
請勿填寫 filenames. 然後執行 `wuliu-metadata -json metadata.json`
可預覽批量修改後的檔案屬性，尚未實際執行。

使用參數 `-danger` 纔會實際執行，
例如 `wuliu-metadata -json metadata.json -danger`

默認只有填寫了內容的項目會被修改，空值項目保持不變 (不會被改為空值)。
如果想讓空值也生效，請使用參數 `-omitempty=false`,
例如 `wuliu-metadata -json metadata.json -omitempty=false`

## wuliu-like (點讚，方便尋找精品或常用檔案)

- `wuliu-like -id ID -n=3` 把一个文件的 like (小心心/点赞) 设为 3,
  其中 ID 是文件的 ID, n 是一个整数，数字越大表示越喜欢/越重要。
- `wuliu-like -id ID -n=0` 把 n 设为零，取消点赞。
- 也可以不輸入 n, 默認 `-n=1`
- 点赞或取消点赞后，需要执行 `wuliu-db -update=cache` 更新索引缓冲。

## 未为视频文件优化

- 视频文件通常较大
- 大文件比较麻烦
- 尤其是比较大的视频文件，不建议使用本软件来管理
- 体积不大的视频文件，可以使用本软件，但没有优化，只当作普通文件处理
  （没有预览和播放视频等功能）
- 可以考虑单独建立一个专案，专门用来管理视频文件，这样在 files 文件夹里
  就只有视频文件，比较方便预览和播放

## 相册功能

- 相册功能，就是针对图片文件的功能，方便预览和管理图片
- 相册功能采用 Python 语言实现，详情请看 <README-Python.md>
- 同时，相册功能也是一个演示，证明可以使用任何编程语言来扩展本软件的功能

## TODO

- wuliu-checksum -same
- wuliu-db -dump 导出整个数据库到 msgpack, 方便其他编程语言使用
- wuliu-list -ctime="2024-02-01" 通過日期前綴後列印檔案
