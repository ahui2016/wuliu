package util

import (
	"fmt"
	"hash/crc32"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	GB              = 1 << 30
	MB              = 1024 * 1024
	Day             = 24 * 60 * 60
	RFC3339         = "2006-01-02 15:04:05Z07:00"
	MIMEOctetStream = "application/octet-stream"
	NormalFilePerm  = 0666
	NormalDirPerm   = 0750
	RepoName        = "Wuliu File Manager"
	RepoURL         = "https://github.com/ahui2016/wuliu"
	ProjectInfoPath = "project.json"
	FileCheckedPath = "file_checked.json"
	DatabasePath    = "project.db"
)

const (
	FILES      = "files"
	METADATA   = "metadata"
	INPUT      = "input"
	BUFFER     = "buffer"
	WEBPAGES   = "webpages"
	RECYCLEBIN = "recyclebin"
)

var Folders = []string{
	FILES,
	METADATA,
	INPUT,
	BUFFER,
	WEBPAGES,
	RECYCLEBIN,
}

var (
	Epoch     = time.Unix(0, 0).Format(RFC3339)
	Separator = string(filepath.Separator)
)

type (
	Base64String = string
	HexString    = string
)

type ProjectInfo struct {
	RepoName        string   // 用于判断资料夹是否 Wuliu 专案
	ProjectName     string   // 备份时要求专案名称相同
	IsBackup        bool     // 是否副本（副本禁止添加、删除等）
	Projects        []string // 第一个是主专案，然后是备份专案
	LastBackupAt    []string // 上次备份时间
	CheckInterval   int      // 检查完整性, 单位: day
	CheckSizeLimit  int      // 检查完整性, 单位: MB
	ExportSizeLimit int64    // 導出檔案體積上限，單位: MB
}

func NewProjectInfo(name string) (info ProjectInfo) {
	info.RepoName = RepoName
	info.ProjectName = name
	info.Projects = []string{"."} // 注意必须确保第一个是 "."
	info.LastBackupAt = []string{Epoch}
	info.CheckInterval = 30
	info.CheckSizeLimit = 1024
	info.ExportSizeLimit = 300
	return
}

type ProjectStatus struct {
	*ProjectInfo
	Root         string // 专案根目录
	TotalSize    int    // 全部檔案體積合計
	FilesCount   int    // 檔案數量合計
	DamagedCount int    // 受損檔案數量合計
}

// EditFiles 用于批量修改档案属性。
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

func NewEditFiles(ids, filenames []string) *EditFiles {
	ef := new(EditFiles)
	ef.IDs = ids
	ef.Filenames = filenames
	ef.Keywords = []string{}
	ef.Collections = []string{}
	ef.Albums = []string{}
	return ef
}

func (ef *EditFiles) Check() (err error) {
	if len(ef.IDs) > 0 && len(ef.Filenames) > 0 {
		err = fmt.Errorf("只能指定 ID 或檔案名稱，不可兩者同時指定。")
	}
	return
}

type FileChecked struct {
	ID      string // 档案名称的 CRC32
	Checked string // RFC3339 上次校驗檔案完整性的時間
	Damaged bool   // 上次校驗結果 (檔案是否損壞)
}

type File struct {
	ID          string   `json:"id"`          // 档案名称的 CRC32
	Filename    string   `json:"filename"`    // 档案名称
	Checksum    string   `json:"checksum"`    // BLAKE2b
	Size        int64    `json:"size"`        // length in bytes for regular files
	Type        string   `json:"type"`        // 檔案類型, 例: text/js, office/docx
	Like        int64    `json:"like"`        // 點贊
	Label       string   `json:"label"`       // 标签，便於搜尋
	Notes       string   `json:"notes"`       // 備註，便於搜尋
	Keywords    []string `json:"keywords"`    // 關鍵詞, 便於搜尋
	Collections []string `json:"collections"` // 集合（分组），一个档案可属于多个集合
	Albums      []string `json:"albums"`      // 相册（专辑），主要用于图片和音乐
	CTime       string   `json:"ctime"`       // RFC3339 檔案入庫時間
	UTime       string   `json:"utime"`       // RFC3339 檔案更新時間
}

type FileAndMeta struct {
	*File
	Metadata []byte
}

func NewFile(name string) *File {
	now := Now()
	f := new(File)
	f.ID = NameToID(name)
	f.Filename = name
	f.CTime = now
	f.UTime = now
	return f
}

func NamesToIds(names []string) (ids []string) {
	for _, name := range names {
		ids = append(ids, NameToID(name))
	}
	return
}

// NameToID 目前采用 CRC32Str36
func NameToID(name string) string {
	return CRC32Str36(name)
}

// CRC32Str36 把一个字符串转化为 crc32, 再转化为 36 进制。
func CRC32Str36(s string) string {
	sum := crc32.ChecksumIEEE([]byte(s))
	str36 := strconv.FormatUint(uint64(sum), 36)
	return strings.ToUpper(str36)
}

func Now() string {
	return time.Now().Format(RFC3339)
}
