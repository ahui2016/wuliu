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
	files := util.NewFilesFrom(names)
	for _, f := range files {
		fmt.Println(f.ID, f.Filename)
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
