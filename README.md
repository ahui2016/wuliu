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
- 执行 `wuliu-add` 默认发现 files 和 input 里的新档案
  - 执行 `wuliu-add --from=files` 只发现 files 里的新档案
  - 执行 `wuliu-add --from=input` 只发现 input 里的新档案
- 需要添加属性 `--write` 才能真正添加新档案，否则就只是列出新档案

### 只添加一部分新档案

- 执行命令 `wuliu-add --output=files.txt` 可以把发现的新档案列印到 files.txt 中
- 注意, files.txt 在 `input` 资料夹中（注意防止覆盖）
- 在 files.txt 中删除不需要添加的档案名称
- 执行命令 `wuliu-add --files=files.txt` 只添加指定的新档案
- 使用 `--files` 参数时, `--from` 参数无效

### wuliu-add --json=common.json

- 执行 `wuliu-add --json=common.json` 默认发现 files 和 input 里的新档案，
  同时列出 common.json 里的档案属性，该属性将应用于全部被发现的新档案。
- 注意, common.json 应放在 input 资料夹中。
- `--json` 与 `--files` 可组合使用
- 需要添加属性 `--write` 才能真正添加新档案，否则就只是列印相关信息

