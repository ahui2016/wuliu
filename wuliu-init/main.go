package main

import (
	"fmt"
	"github.com/ahui2016/wuliu/util"
)

func main() {
	fmt.Println(checkCWD())
}

func checkCWD() error {
	cwd := util.GetCwd()
	ok, err := util.DirIsEmpty(cwd)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("當前目錄不為空: %s", cwd)
	}
	return nil
}
