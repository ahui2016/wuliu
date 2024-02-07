package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	"log"
	"path/filepath"
)

var (
	MainProjInfo = util.ReadProjectInfo(".")
)

var (
	projectsFlag = flag.Bool("projects", false, "list all projects")
	nFlag        = flag.Int("n", 0, "select a project by a number")
	backupFlag   = flag.Bool("backup", false, "do backup files")
)

type (
	FileChecked   = util.FileChecked
	ProjectInfo   = util.ProjectInfo
	ProjectStatus = util.ProjectStatus
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	if *projectsFlag {
		printProjectsList()
		return
	}

	if *nFlag > 0 {
		bkRoot := getBkRoot()
		mainStatus, bkStatus := getProjectsStatus(".", bkRoot)
		printStatus(mainStatus, bkStatus)
		if err := checkStatus(mainStatus, bkStatus); err != nil {
			log.Fatalln(err)
		}
	}

	if *backupFlag {
		// bkRoot := getBkRoot()
		return
	}

}

func getBkRoot() string {
	if *nFlag == 0 {
		log.Fatalln("請使用參數 '-n' 指定備份專案")
	}
	return MainProjInfo.Projects[*nFlag]
}

func getProjectsStatus(mainRoot, bkRoot string) (mainStatus, bkStatus ProjectStatus) {
	mainDB := lo.Must(util.OpenDB("."))
	defer mainDB.Close()

	bkDB := lo.Must(util.OpenDB(bkRoot))
	defer bkDB.Close()

	mainProjInfo := util.ReadProjectInfo(mainRoot)
	fileN, totalSize := lo.Must2(util.DatabaseFilesSize(mainDB))
	fcList := lo.Must(util.ReadFileChecked("."))
	damagedFiles := lo.Filter(fcList, func(fc *FileChecked, _ int) bool {
		return fc.Damaged
	})
	mainStatus.ProjectInfo = &mainProjInfo
	mainStatus.Root = "."
	mainStatus.TotalSize = totalSize
	mainStatus.FilesCount = fileN
	mainStatus.DamagedCount = len(damagedFiles)

	bkProjInfo := util.ReadProjectInfo(bkRoot)
	fileN, totalSize = lo.Must2(util.DatabaseFilesSize(bkDB))
	fcList = lo.Must(util.ReadFileChecked(bkRoot))
	damagedFiles = lo.Filter(fcList, func(fc *FileChecked, _ int) bool {
		return fc.Damaged
	})
	bkStatus.ProjectInfo = &bkProjInfo
	bkStatus.Root = bkRoot
	bkStatus.TotalSize = totalSize
	bkStatus.FilesCount = fileN
	bkStatus.DamagedCount = len(damagedFiles)

	return
}

func printProjectsList() {
	bkProjects := MainProjInfo.Projects[1:]
	if len(bkProjects) == 0 {
		fmt.Println("無備份專案。")
		fmt.Println("添加備份專案的方法請參閱", util.RepoURL)
		return
	}
	for i, project := range bkProjects {
		fmt.Printf("%d %s\n", i+1, project)
	}
}

func syncProjInfo(bkRoot string) error {
	bkProjInfo := MainProjInfo
	bkProjInfo.IsBackup = true
	bkProjInfoPath := filepath.Join(bkRoot, util.ProjectInfoPath)
	_, err := util.WriteJSON(bkProjInfo, bkProjInfoPath)
	return err
}

// 检查 ProjectName 相同，检查 IsBakcup == true, 列印两个数据库的档案数量、
// 上次备份日期、损坏档案，有损坏档案禁止备份。
func checkStatus(mainStatus, bkStatus ProjectStatus) error {
	if mainStatus.ProjectName != bkStatus.ProjectName {
		return fmt.Errorf("專案名稱不一致: '%s' ≠ '%s'\n", mainStatus.ProjectName, bkStatus.ProjectName)
	}
	if !bkStatus.IsBackup {
		return fmt.Errorf("不是備份專案: %s 裏的 IsBackup 是 false\n")
	}
	if mainStatus.DamagedCount+bkStatus.DamagedCount > 0 {
		return fmt.Errorf("發現受損檔案，必須修復後纔能備份。\n")
	}
	return nil
}

func printStatus(mainStatus, bkStatus ProjectStatus) {
	fmt.Printf("源專案\t\t%s\n", mainStatus.Root)
	fmt.Printf("檔案數量\t%d\n", mainStatus.FilesCount)
	fmt.Printf("體積合計\t%s\n", mainStatus.TotalSize)
	fmt.Printf("受損檔案\t%d\n", mainStatus.DamagedCount)
	fmt.Println()
	fmt.Printf("目標專案\t%s\n", bkStatus.Root)
	fmt.Printf("檔案數量\t%d\n", bkStatus.FilesCount)
	fmt.Printf("體積合計\t%s\n", bkStatus.TotalSize)
	fmt.Printf("受損檔案\t%d\n", bkStatus.DamagedCount)
}
