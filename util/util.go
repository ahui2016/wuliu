package util

import (
	"encoding/json"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

const (
	NormalFilePerm  = 0666
	NormalDirPerm   = 0750
	ProjectInfoPath = "metadata/project.json"
	DatabasePath    = "metadata/project.db"
)

var Separator = string(filepath.Separator)

func GetCwd() string {
	return lo.Must(os.Getwd())
}

// GetExePath returns the path name for the executable
// that started the current process.
func GetExePath() string {
	return lo.Must1(os.Executable())
}

func DirIsEmpty(dirpath string) (ok bool, err error) {
	items, err := filepath.Glob(dirpath + Separator + "*")
	ok = len(items) == 0
	return
}

func DirIsNotEmpty(dirpath string) (ok bool, err error) {
	ok, err = DirIsEmpty(dirpath)
	return !ok, err
}

func PathNotExists(name string) (ok bool) {
	_, err := os.Lstat(name)
	if os.IsNotExist(err) {
		ok = true
		err = nil
	}
	lo.Must0(err)
	return
}

func PathExists(name string) bool {
	return !PathNotExists(name)
}

// MkdirIfNotExists 创建資料夹, 忽略 ErrExist.
// 在 Windows 里, 文件夹的只读属性不起作用, 为了统一行为, 不把資料夹设为只读.
func MkdirIfNotExists(name string) error {
	if PathExists(name) {
		return nil
	}
	return os.Mkdir(name, NormalDirPerm)
}

// WriteFile 写檔案, 使用权限 0666
func WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, NormalFilePerm)
}

// WriteJSON 把 data 转换为漂亮格式的 JSON 并写入檔案 filename 中。
func WriteJSON(data interface{}, filename string) error {
	dataJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return WriteFile(filename, dataJSON)
}
