# Wuliu File Manager (五柳档案管理脚本)

最近我精力下降得厉害，没有能力维护程序代码，因此追求程序的极致简单、易维护，其他都不重要。

注意，写着写着很容易不小心变复杂，一定要警惕，保持简单。

这是一些档案管理的零散脚本，不是一个完整的程序。

## 名词

- 档案 = 文件 = file
- 资料夹 = 文件夹 = folder = directory
- 专案 = 项目 = project

## 主要目的

- 给档案增加更多属性，比如分类、标签、点赞
- 更方便地备份档案
- 检查档案完整性

### 档案属性

- 每个档案与一个 json 关联, json 里有档案的属性（例如分类、标签、备注等）
- 属性可自由添加，只要脚本能处理即可
- 档案属性缓冲使用数据库 https://github.com/etcd-io/bbolt
- 一切以实际档案与 json 为准，建议每隔一段时间就重建缓冲数据库

## 编程简单第一

- 非常重视编程简单
- 运行效率和使用方便都是次要的

## 基本原理

- 可以处理一个普通的资料夹，但只能处理根目录，不能处理子资料夹（为了编程简单）
- 使用 JSON 作为数据库（编程简单，效率次要）
- 加密使用 VeraCrypt

## 基本功能

- 添加档案：直接添加
- 更改档案名称：使用脚本更改档案名称，可保留档案属性
- 删除档案：直接删除或使用脚本，建议先移动至 `recyclebin`
- 检查档案完整性：提供缓存
- 搜索档案：提供缓存
- 可方便地查看缓存新限度（上次缓存时间）

### 专案与子资料夹

- 没有子资料夹
- 如果需要加密，应独立新建一个专案
- 如果档案太多遇到性能问题，应另外新建一个专案
- 这样做是为了保持代码简单

## 使用方法

### 新建专案

对一个空资料夹新建专案（初始化），会生成几个资料夹：
  
  - `files`
  - `metadata`
  - `input`
  - `output`
  - `webpages`
  - `recyclebin`

其中 metadata 里的 json 档案与 files 里的普通档案一一对应，
其他任何档案请勿放进 metadata 里，并且，metadata 里也不允许有子资料夹。

注意: files, metadata, input 均不允许包含子资料夹！

### 添加档案

- 直接向 `files` 资料夹添加档案，然后执行脚本去发现新档案
- 先把新档案放入 `input` 资料夹，可提高程序效率
- 提供一个脚本，方便新档案的批量处理（比如批量设定分类、标签）
- 档案 ID 使用年份加档案名的 https://pkg.go.dev/hash/adler32

### 更改档案名称

- 使用脚本更改档案名称，可保留档案属性
- 如果直接修改档案名称，可执行脚本发现数据库与实际文件的差异

### 删除档案

- 使用脚本删除档案，可同时更新数据库
- 如果直接删除档案，可执行脚本发现数据库与实际文件的差异

### 批量处理

- 添加档案时，修改档案属性时（分类、标签）可进行批量处理
- 配置文件: JSON
- 批量处理提供预览功能
- Windows 资源管理器可以复制档案路径（可多选），利用这个特性，读取剪贴板进行批量处理
- 也可以通过静态网页 javascript 进行多选

### 生成静态网页

- 生成一些静态网页，方便浏览和复制档案名
- 专为 markdown 档案生成一些静态网页
- 网页放在 `wuliu-webpages` 资料夹中

### 相册

- 建议建立一个以图片为主的专案
- 使用第三方图片软件浏览图片
- 后期考虑专为图片生成一些静态网页

## 数据库

- TODO: 等数据库大了之后，试试 `compact` 后 `os.Rename` 看看能否缩小体积。
  - 大概率不需要做这个，因为提供了 "重建数据库" 功能。

## 备份

- 单向备份
- 只需要找出新增的档案和删除的档案
- 分别同步档案和数据库
- 直接复制数据库会有一些问题（比如查错时间），但为了简化程序。

### 直接复制数据库

- 「源数据库」与「目标数据库」的内容几乎完全一致，只有「上次查错时间」不同
- 「上次查错时间」单独处理? (file_checked.json)

## 查错（档案完整性）与修复

- 只查错，不修复，只能手动修复（下载档案、重新上传，或者手动复制粘贴）
- 注意添加、删除、更新档案内容时要更新 file_checked.json

## 编写程序

- 尝试使用最简单的文本编辑器 (EmEditor)
- 配合 ripgrep 使用 <https://github.com/BurntSushi/ripgrep>
- https://pkg.go.dev/flag

