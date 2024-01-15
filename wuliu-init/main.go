package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	"log"
)

var (
	vFlag = flag.Bool("v", false, "print the version of the command")
	wFlag = flag.Bool("where", false, "print where is the command")
)

var Folders = []string{
	util.FILES, util.METADATA, util.INPUT, util.OUTPUT, util.WEBPAGES, util.RECYCLEBIN,
}

func main() {
	flag.Parse()
	util.PrintVersionExit(*vFlag)
	util.PrintWhereExit(*wFlag)

	checkCWD()
	makeFolders()
	writeProjectInfo()
	util.CreateDatabase()
}

func checkCWD() {
	cwd := util.GetCwd()
	if lo.Must(util.DirIsNotEmpty(cwd)) {
		log.Fatalln("當前目錄不為空:", cwd)
	}
}

func makeFolders() {
	for _, folder := range Folders {
		fmt.Println("Create folder:", folder)
		lo.Must0(util.MkdirIfNotExists(folder))
	}
}

func writeProjectInfo() {
	fmt.Println("Create", util.ProjectInfoPath)
	lo.Must0(util.WriteProjectInfo(util.DefaultWuliuInfo))
}
