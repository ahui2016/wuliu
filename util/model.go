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
	Day             = 24 * 60 * 60
	RFC3339         = "2006-01-02 15:04:05Z07:00"
	MIMEOctetStream = "application/octet-stream"
	NormalFilePerm  = 0666
	NormalDirPerm   = 0750
	ProjectInfoPath = "project.json"
	DatabasePath    = "project.db"
)

const (
	FILES      = "files"
	METADATA   = "metadata"
	INPUT      = "input"
	OUTPUT     = "output"
	WEBPAGES   = "webpages"
	RECYCLEBIN = "recyclebin"
)

var Separator = string(filepath.Separator)

type (
	Base64String = string
	HexString    = string
)

type ProjectInfo struct {
	RepoName         string
	RepoURL          string
	OrphanLastCheck  string // 上次检查孤立档案的时间
	OrphanFilesCount int    // 孤立的档案数量
	OrphanMetaCount  int    // 孤立的 metadata 数量
}

var DefaultWuliuInfo = ProjectInfo{
	RepoName: "Wuliu File Manager",
	RepoURL:  "https://github.com/ahui2016/wuliu",
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
	Checked     string   `json:"checked"`     // RFC3339 上次校驗檔案完整性的時間
	Damaged     bool     `json:"damaged"`     // 上次校驗結果 (檔案是否損壞)
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
	f.Checked = now
	return f
}

type FilesToDelete struct {
	IDs   []string
	Names []string
}

func (files *FilesToDelete) Check() error {
	idsLength := len(files.IDs)
	namesLength := len(files.Names)
	if idsLength+namesLength == 0 {
		return fmt.Errorf("未指定要刪除的檔案")
	}
	if idsLength > 0 && namesLength > 0 {
		return fmt.Errorf("只能指定 ID 或檔案名稱，不可兩者同時指定。")
	}
	return nil
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
