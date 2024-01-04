# wuliu

Wuliu File Manager (五柳档案管理脚本)

## 名词

- 档案 = 文件 = file
- 资料夹 = 文件夹 = folder = directory
- 专案 = 项目 = project

## Scripts

- wuliu-init (新建/初始化一个专案)

## wuliu-init

- 新建/初始化一个专案，主要是新建一些资料夹和数据库、配置等。
- 只能对一个空资料夹进行初始化
- 使用方法: 进入一个空资料夹，执行 `wuliu-init` (没有任何参数，只能初始化当前目录)

## wuliu-add

- 该命令用于添加档案，同时也用于发现新档案
- 需要添加属性 `--do` 才能真正添加新档案，否则就只是列出新档案
- 执行 `wuliu-add` 发现 files 和 input 里的新档案
  - 执行 `wuliu-add --from=files` 只发现 files 里的新档案 (取消，为了编程简单)
  - 执行 `wuliu-add --from=input` 只发现 input 里的新档案 (取消，为了编程简单)
  - 差点又搞复杂了，一定要非常警惕，保持编程简单

### 只添加一部分新档案

- 执行命令 `wuliu-add --writenames=files.txt`
  可以把发现的新档案名称列印到 files.txt 中
- files.txt 在 `input` 资料夹中（注意防止覆盖）
- 在 files.txt 中删除不需要添加的档案名称
- 执行命令 `wuliu-add --files=files.txt` 只添加指定的新档案

### wuliu-add --json=common.json

- 执行 `wuliu-add --json=common.json` 发现 files 和 input 里的新档案，
  同时列出 common.json 里的档案属性，该属性将应用于待添加的新档案。
- 注意, common.json 应放在 input 资料夹中。
- 执行 `wuliu-add --newjson=common.json`
  可在 input 资料夹中生成一个新的 common.json, 方便编辑
- `--json` 与 `--files` 可组合使用
- 需要添加属性 `--do` 才能真正添加新档案，否则就只是列印相关信息

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
- 关于 CRC32 <https://softwareengineering.stackexchange.com/questions/49550/which-hashing-algorithm-is-best-for-uniqueness-and-speed>
- Type, Label, Note, Keywords 等都是为了方便搜寻，请大胆灵活使用。
- Keywords, Collections 等 `[]string` 类型，都排序，排序后转为纯字符
  （用逗号空格 `, ` 分隔）方便保存到 kv 数据库。
- 因此 `[]string` 类型在用户输入时不允许包含逗号、顿号和空格。




