package util

import (
	"encoding/json"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

var Separator = string(filepath.Separator)

func GetCwd() string {
	return lo.Must(os.Getwd())
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
