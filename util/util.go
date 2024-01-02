package util

import (
	"os"
	"path/filepath"
	"github.com/samber/lo"
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
