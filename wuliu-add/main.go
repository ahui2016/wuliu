package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	"os"
)

type (
	File = util.File
)

var (
	do = flag.Bool("do", false, "really do add files")
)

func main() {
	flag.Parse()

	util.MustInWuliu()
	checkOrphan()
	names := lo.Must(util.FindNewFiles())
	files := lo.Must(util.NewFilesFromInput(names))
	checkExist(files)

	if *do {
		addNew(files)
	} else {
		findNew(files)
	}
}

func findNew(files []*File) {
	for _, f := range files {
		fmt.Println(f.ID, f.Filename)
	}
}

func addNew(files []*File) {
	for _, f := range files {
		metaPath := util.METADATA + "/" + f.Filename + ".json"
		fmt.Println("Create =>", metaPath)
		lo.Must0(util.WriteJSON(f, metaPath))
		src := util.INPUT + "/" + f.Filename
		dst := util.FILES + "/" + f.Filename
		fmt.Println("Add =>", dst)
		lo.Must0(os.Rename(src, dst))
	}
}

func checkExist(files []*File) {
	var dstFiles []string
	for _, f := range files {
		dst := util.FILES + "/" + f.Filename
		meta := util.METADATA + "/" + f.Filename + ".json"
		dstFiles = append(dstFiles, dst, meta)
	}
	var existFiles []string
	for _, f := range dstFiles {
		if util.PathExists(f) {
			existFiles = append(existFiles, f)
		}
	}
	if len(existFiles) > 0 {
		fmt.Println("【注意！】同名檔案已存在：")
		util.PrintList(existFiles)
		os.Exit(0)
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
