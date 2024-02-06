package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	"log"
	"os"
)

var (
	nameFlag = flag.String("name", "", "set a unique name for the project")
	vFlag    = flag.Bool("v", false, "print the version of the command")
	wFlag    = flag.Bool("where", false, "print where is the command")
)

var Folders = []string{
	util.FILES, util.METADATA, util.INPUT, util.OUTPUT, util.WEBPAGES, util.RECYCLEBIN,
}

func main() {
	customFlagUsage()
	flag.Parse()
	util.PrintVersionExit(*vFlag)
	util.PrintWhereExit(*wFlag)

	if *nameFlag == "" {
		flag.Usage()
		return
	}
	checkCWD()
	makeFolders()
	writeProjectInfo(*nameFlag)
	writeFileChecked()
	util.CreateDatabase()
}

// customFlagUsage 必须在 `flag.Parse()` 之前执行才有效。
func customFlagUsage() {
	cmdUsage := "在空资料夹内执行 `wuliu-init -name` 初始化专案。"
	flag.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(), "Usage of %s:\n%s\n", os.Args[0], cmdUsage)
		flag.PrintDefaults()
	}
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

func writeProjectInfo(name string) {
	fmt.Println("Create", util.ProjectInfoPath)
	info := util.NewProjectInfo(name)
	lo.Must0(util.WriteProjectInfo(info))
}

func writeFileChecked() {
	_ = lo.Must(
		util.WriteJSON([]int{}, util.FileCheckedPath))
}
