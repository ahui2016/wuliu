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
	MainProject = util.ReadProjectInfo()
)

var (
	projectsFlag = flag.Bool("projects", false, "list all projects")
	nFlag        = flag.Int("n", 0, "select a project by a number")
	addFlag      = flag.String("add", "", "add a new backup-project")
	backupFlag   = flag.Bool("backup", false, "do backup files")
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	if *addFlag != "" {
		util.FolderMustEmpty(*addFlag)
		util.MakeFolders(true)
		writeProjectInfo(*nameFlag)
		util.InitFileChecked()
		util.CreateDatabase()
		return
	}

	if *backupFlag {
		if *nFlag == 0 {
			log.Fatalln("請使用參數 '-n' 指定備份專案")
		}
		backupRoot := MainProject.Projects[*nFlag]
		return
	}

}

func backupProjectInfo()
