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

func main() {
	customFlagUsage()
	flag.Parse()
	util.PrintVersionExit(*vFlag)
	util.PrintWhereExit(*wFlag)

	if *nameFlag == "" {
		flag.Usage()
		return
	}
	util.FolderMustEmpty(".")
	util.MakeFolders(true)
	writeProjectInfo(*nameFlag)
	util.InitFileChecked()
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

func writeProjectInfo(name string) {
	fmt.Println("Create", util.ProjectInfoPath)
	info := util.NewProjectInfo(name)
	lo.Must0(util.WriteProjectInfo(info))
}
