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

func isRegularFile(name string) (ok bool, err error) {
	info, err := os.Lstat(name)
	if err != nil {
		return
	}
	return info.Mode().IsRegular(), nil
}

// GetFilenamesBase 假设 folder 里全是普通档案，没有资料夹。
func GetFilenamesBase(folder string) ([]string, error) {
	pattern := filepath.Join(folder, "*")
	names, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	baseNames := lo.Map(names, func(name string, _ int) string {
		return filepath.Base(name)
	})
	return baseNames, nil
}

func getRegularFiles(folder string) (files []string, err error) {
	pattern := filepath.Join(folder, "*")
	items, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	for _, file := range items {
		ok, err := isRegularFile(file)
		if err != nil {
			return nil, err
		}
		if ok {
			files = append(files, file)
		}
	}
	return files, nil
}
