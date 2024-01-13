package main

import (
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	"os"
)

func main() {
	util.MustInWuliu()
	checkOrphan()
	findNewFiles()
}

func findNewFiles() {
	names := lo.Must(util.FindNewFiles())
	for _, name := range names {
		fmt.Println(name)
	}
}

func checkOrphan() {
	info := util.ReadProjectInfo()
	if info.OrphanFilesCount+info.OrphanMetaCount > 0 {
		fmt.Println("發現孤立檔案, 請執行 wuliu-orphan 進行檢查")
		fmt.Println("上次檢查時間:", info.OrphanLastCheck)
		os.Exit(0)
	}
}
