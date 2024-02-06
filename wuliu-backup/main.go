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
	projectsFlag  = flag.Bool("projects", false, "list all projects")
	nFlag         = flag.Int("n", 0, "select a project by a number")
	backupFlag    = flag.Bool("backup", false, "do backup files")
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	if *backupFlag {
		if *nFlag == 0 {
			log.Fatalln("請使用參數 '-n' 指定備份專案")
		}
		bkRoot := MainProject.Projects[*nFlag]
		return
	}

}

func syncProjInfo(bkRoot string) error {
	bkProjInfo := MainProject
	bkProjInfo.IsBackup = true
	bkProjInfoPath := filepath.Join(bkRoot, util.ProjectInfoPath)
	_, err := WriteJSON(bkProjInfo, bkProjInfoPath)
	return err
}

func backupProjectInfo() {}
