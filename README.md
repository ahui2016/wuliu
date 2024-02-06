# wuliu

Wuliu File Manager (五柳档案管理脚本)

## 名词

- 档案 = 文件 = file
- 资料夹 = 文件夹 = folder = directory
- 专案 = 项目 = project
- 列印 = 打印/显示 = print

## Scripts

- wuliu-init (新建/初始化一个专案)
- `cat project.info` (查看专案信息)
- wuliu-orphan (检查有无孤立档案)
- wuliu-add (添加档案)
- wuliu-delete (删除档案)
- wuliu-list (列印档案、标签、备注、关键词等)
- wuliu-db (数据库信息，更新缓存)
- wuliu-checksum (检查档案完整性)
- wuliu-replace (更新同名档案，覆盖旧档案)
- wuliu-backup (备份专案)

## wuliu-init

- 新建/初始化一个专案，主要是新建一些资料夹和数据库、配置等。
- 只能对一个空资料夹进行初始化
- 使用方法: 进入一个空资料夹，执行 `wuliu-init -name` 进行初始化。
- 注意，请为不同的专案设定不同的名称，备份时有用。
- 备份专案（详见关于 `wuliu-backup` 的说明）的专案名称必须与源专案一致。
- `wuliu-init -h` 列印帮助信息
- `wuliu-init -v` 列印版本信息
- `wuliu-init -where` 列印 wuliu-init 的位置

## project.json

建议经常执行 `cat project.json` 查看专案信息。
当然，也可直接打开 project.json 查看。

```
type ProjectInfo struct {
	RepoName         string
	RepoURL          string
	IsBackup         bool     // 是否副本（副本禁止添加、删除等）
	Projects         []string // 第一个是主专案，然后是备份专案
	LastBackupAt     []string // 上次备份时间
	CheckInterval    int      // 检查完整性, 单位: day
	CheckSizeLimit   int      // 检查完整性, 单位: MB
	OrphanLastCheck  string   // 上次检查孤立档案的时间
	OrphanFilesCount int      // 孤立的档案数量
	OrphanMetaCount  int      // 孤立的 metadata 数量
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
	ID			string	`json:"id"`				// 档案名称的 CRC32
	Filename	string	`json:"filename"`			// 档案名称
	Checksum	string	`json:"checksum"`		// BLAKE2b
	Size		int64	`json:"size"`				// length in bytes for regular files
	Type		string	`json:"type"`				// 檔案類型, 例: text/js, office/docx
	Like		int64	`json:"like"`				// 點贊
	Label		string	`json:"label"`			// 标签，便於搜尋
	Notes		string	`json:"notes"`			// 備註，便於搜尋
	Keywords	[]string	`json:"keywords"`		// 關鍵詞, 便於搜尋
	Collections	[]string	`json:"collections"`		// 集合（分组），一个档案可属于多个集合
	Albums		[]string	`json:"albums"`			// 相册（专辑），主要用于图片和音乐
	CTime		string	`json:"ctime"`			// RFC3339 檔案入庫時間
	UTime		string	`json:"utime"`			// RFC3339 檔案更新時間
	Checked	string	`json:"checked"`			// RFC3339 上次校驗檔案完整性的時間
	Damaged	bool	`json:"damaged"`		// 上次校驗結果 (檔案是否損壞)
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

## wuliu-delete

- 该命令删除添加档案，包括删除对应的 json 档案和数据库中的条目
- `wuliu-delete -id` 通过档案 ID 指定需要删除的档案（只能指定一个）
- `wuliu-delete -name` 通过档案名称指定需要删除的档案（只能指定一个）
- `wuliu-delete --newjson delete.json` 在专案根目录生成一个 delete.json 档案模板，
  方便批量填写需要删除的档案。
- 例如想删除 files/aaa.txt 档案，请在 delete.json 中指定档案名称 aaa.txt,
  如果想删除 metadata/bbb.txt.json, 请在 delete.json 中指定档案名称 bbb.txt (不需要 `.json`)
- `wuliu-delete --json delete.json ` 通过 delete.json 指定需要删除的档案（可指定多个）
- 需要添加属性 `--danger` 才能真正删除档案，否则就只是列出 delete.json 的内容

## wuliu-list

- `wuliu-list` 列印最近 15 个档案 (ID, 体积, 档案名称)
- `wuliu-list n=100` 列印最近 100 个档案，按 CTime 倒序排列 (CTime 是入库时间)

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

### 更新数据库

- 更新数据库，是指以 metadata 为准更新数据库，因此如果一段时间没执行 wuliu-orphan,
  请先执行一次 wuliu-orphan 再更新数据库。
- 执行 `wuliu-db --update=rebuild` 根据 metadata(真实的 json 档案) 重建整个数据库。
  执行 `wuliu-db --update=cache` 根据缓存更新索引（不需要读取硬盘里的 json 档案）。
- 由于数据库缓存（即 files 索引和 filename 索引）在添加文件、修改文件属性、删除文件时
  会自动更新，因此多数情况下只需要 `--update=cache`, 不需要重建数据库。

## wuliu-checksum

- `wuliu-checksum --renew` 将全部文件的 damaged 设为 false, 上次检查时间设为 epoch
- `wuliu-checksum --check` 校验文件完整性（看文件是否损坏）
- `wuliu-checksum --projects` 列印全部专案
- `wuliu-checksum -n` 通过序号选择专案，默认是 0 (即当前专案)

在执行 `wuliu-checksum` 命令时，有时会显示以下信息：

```
已選擇專案: ./
數據庫檔案數量: 5
待檢查檔案數量: 5
```

如果发现 "數據庫檔案數量" 与 "待檢查檔案數量" 不一致，
可以执行 `wuliu-checksum --renew` 进行修正。

执行 `wuliu-checksum --check` 时，会根据 project.json 中的 CheckInterval (检查周期)
自动判断档案是否需要检查，根据 CheckSizeLimit (检查体积上限) 自动终止检查，防止
单次检查时间太长。

## wuliu-backup

创建新备份专案的方法：

- 请先进入一个空资料夹（新的备份专案的根目录，以下称为 backupRoot）
- 把现有专案（以下称为“主专案”）的 files, metadata 这两个资料夹，
  以及 project.db, project.json 这两个文件复制到 backupRoot 中。
- 编辑 backupRoot 里的 project.json, 把 IsBackup 的值改为 true
- 编辑主专案中的 project.json, 把 backupRoot 的路径添加到 Projects 列表中。
  注意，路径里的反斜杠改要为 "\\" 或 "/"

以上是创建新备份专案的方法，以下是 wuliu-backup 的其他命令：

- `wuliu-backup --projects` 列印全部备份专案
- “备份专案”是指专门用于备份的专案
- 本软件采用单向备份方式，备份时以“主专案”为准，
  使备份专案里的档案变成与主专案一样。
- `wuliu-backup -n` 通过序号选择备份专案
- `wuliu-backup -backup` 正式执行备份
- 例如执行命令 `wuliu-backup -n 1` 会列印第 1 个备份专案的信息，但不会执行备份。
  而执行命令  `wuliu-backup -n=1 -backup` 则会正式执行备份。
