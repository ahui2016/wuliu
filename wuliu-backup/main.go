package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/juju/utils/v4/du"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
)

var (
	MainProjInfo = util.ReadProjectInfo(".")
)

var (
	projectsFlag = flag.Bool("projects", false, "list all projects")
	nFlag        = flag.Int("n", 0, "select a project by a number")
	dangerFlag   = flag.Bool("danger", false, "do backup files")
	fixFlag      = flag.Bool("fix", false, "try to fix files automatically")
)

type (
	File          = util.File
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

	var (
		bkRoot       string
		mainDB, bkDB *bolt.DB
	)

	if *nFlag > 0 {
		bkRoot = getBkRoot()
		mainDB = lo.Must(util.OpenDB("."))
		defer mainDB.Close()
		bkDB = lo.Must(util.OpenDB(bkRoot))
		defer bkDB.Close()

		mainStatus, bkStatus := getProjectsStatus(".", bkRoot, mainDB, bkDB)
		printStatus(mainStatus, bkStatus, *nFlag)
		if err := checkStatus(mainStatus, bkStatus, *fixFlag); err != nil {
			fmt.Println("Error!", err)
			return
		}
	}

	if *dangerFlag || *fixFlag {
		if *nFlag < 1 {
			fmt.Println("請使用參數 '-n' 指定目標專案")
			return
		}
	}

	if *dangerFlag {
		fmt.Printf("\n備份開始\n")

		lo.Must0(syncProjInfo(bkRoot, *nFlag))
		lo.Must0(syncFilesToBK(".", bkRoot, mainDB, bkDB))
		fmt.Printf("\n備份結束\n")
		return
	}

	if *fixFlag {
		fmt.Printf("\n嘗試自動修復受損檔案...\n")
		lo.Must0(autoFix(".", bkRoot, mainDB, bkDB))
		return
	}
}

func getBkRoot() string {
	if *nFlag == 0 {
		log.Fatalln("請使用參數 '-n' 指定備份專案")
	}
	return MainProjInfo.Projects[*nFlag]
}

func getProjectsStatus(mainRoot, bkRoot string, mainDB, bkDB *bolt.DB) (mainStatus, bkStatus ProjectStatus) {
	mainProjInfo := util.ReadProjectInfo(mainRoot)
	fileN, totalSize := lo.Must2(util.DatabaseFilesSize(mainDB))
	fcMap := lo.Must(util.ReadFileChecked("."))
	damagedFiles := util.DamagedOfFileChecked(fcMap)
	mainStatus.ProjectInfo = &mainProjInfo
	mainStatus.Root = "."
	mainStatus.TotalSize = totalSize
	mainStatus.FilesCount = fileN
	mainStatus.DamagedCount = len(damagedFiles)

	bkProjInfo := util.ReadProjectInfo(bkRoot)
	fileN, totalSize = lo.Must2(util.DatabaseFilesSize(bkDB))
	fcMap = lo.Must(util.ReadFileChecked(bkRoot))
	damagedFiles = util.DamagedOfFileChecked(fcMap)
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

func syncProjInfo(bkRoot string, n int) error {
	now := util.Now()
	MainProjInfo.LastBackupAt[0] = now
	MainProjInfo.LastBackupAt[n] = now
	fmt.Println("Update =>", util.ProjectInfoPath)
	if err := util.WriteProjectInfo(MainProjInfo); err != nil {
		return err
	}

	bkProjInfo := MainProjInfo
	bkProjInfo.IsBackup = true
	bkProjInfoPath := filepath.Join(bkRoot, util.ProjectInfoPath)
	fmt.Println("Update =>", bkProjInfoPath)
	_, err := util.WriteJSON(bkProjInfo, bkProjInfoPath)
	return err
}

// 检查 ProjectName 相同，检查 IsBakcup == true, 列印两个数据库的档案数量、
// 上次备份日期、损坏档案，有损坏档案禁止备份。
func checkStatus(mainStatus, bkStatus ProjectStatus, fix bool) error {
	if mainStatus.ProjectName != bkStatus.ProjectName {
		return fmt.Errorf("專案名稱不一致: '%s' ≠ '%s'\n", mainStatus.ProjectName, bkStatus.ProjectName)
	}
	if !bkStatus.IsBackup {
		bkProjInfoPath := filepath.Join(bkStatus.Root, util.ProjectInfoPath)
		return fmt.Errorf("不是備份專案: %s 裏的 IsBackup 是 false\n", bkProjInfoPath)
	}
	if !fix && mainStatus.DamagedCount+bkStatus.DamagedCount > 0 {
		return fmt.Errorf("發現受損檔案，必須修復後纔能備份。\n")
	}
	sizeDiff := mainStatus.TotalSize - bkStatus.TotalSize
	return checkBackupDiskUsage(bkStatus.Root, sizeDiff)
}

func checkBackupDiskUsage(volumePath string, addUpSize int64) error {
	usage := du.NewDiskUsage(volumePath)
	if addUpSize <= 0 {
		return nil
	}
	var margin uint64 = 1 << 30 // 1GB
	available := usage.Available()
	sizeStr := util.FileSizeToString(float64(addUpSize), 2)
	availableStr := util.FileSizeToString(float64(available), 2)
	if uint64(addUpSize)+margin > available {
		return fmt.Errorf("磁盤空間不足: %s\nwant %s, available %s\n", volumePath, sizeStr, availableStr)
	}
	return nil
}

func printStatus(mainStatus, bkStatus ProjectStatus, n int) {
	totalSize := util.FileSizeToString(float64(mainStatus.TotalSize), 2)
	mainBackupAt := mainStatus.LastBackupAt[0]
	fmt.Printf("源專案\t\t%s\n", mainStatus.Root)
	fmt.Printf("檔案數量\t%d\n", mainStatus.FilesCount)
	fmt.Printf("體積合計\t%s\n", totalSize)
	fmt.Printf("受損檔案\t%d\n", mainStatus.DamagedCount)
	fmt.Printf("上次備份時間\t%s\n", mainBackupAt)
	fmt.Println()
	totalSize = util.FileSizeToString(float64(bkStatus.TotalSize), 2)
	bkBackupAt := mainStatus.LastBackupAt[n]
	fmt.Printf("目標專案\t%s\n", bkStatus.Root)
	fmt.Printf("檔案數量\t%d\n", bkStatus.FilesCount)
	fmt.Printf("體積合計\t%s\n", totalSize)
	fmt.Printf("受損檔案\t%d\n", bkStatus.DamagedCount)
	fmt.Printf("上次備份時間\t%s\n", bkBackupAt)
	fmt.Println()
	sizeDiff := mainStatus.TotalSize - bkStatus.TotalSize
	diff := util.FileSizeToString(float64(sizeDiff), 2)
	backupAtDiff := lo.Ternary(mainBackupAt == bkBackupAt, "相同", "不同")
	fmt.Printf("源專案檔案數量 - 目標專案檔案數量 = %d\n", mainStatus.FilesCount-bkStatus.FilesCount)
	fmt.Printf("源專案檔案體積 - 目標專案檔案體積 = %s\n", diff)
	fmt.Printf("上次備份時間: %s\n", backupAtDiff)
}

func syncFilesToBK(mainRoot, bkRoot string, mainDB, bkDB *bolt.DB) error {
	files, err := getChangedFiles(mainRoot, bkRoot, mainDB, bkDB)
	if err != nil {
		return err
	}
	return files.Sync()
}

type ChangedFiles struct {
	MainRoot   string
	BkRoot     string
	Deleted    []string
	Updated    []string
	Overwrited []string
	Added      []string
}

func (files ChangedFiles) Sync() (err error) {
	if err = files.syncDelete(); err != nil {
		fmt.Println("Error: delete", err)
		return
	}
	if err = files.syncUpdate(); err != nil {
		fmt.Println("Error: upddate", err)
		return
	}
	if err = files.syncOverwrite(); err != nil {
		fmt.Println("Error: overwrite", err)
		return
	}
	if err = files.syncAdd(); err != nil {
		fmt.Println("Error: add", err)
		return
	}
	return nil
}

func (files ChangedFiles) syncDelete() error {
	for _, name := range files.Deleted {
		fmt.Print(".")
		filePath := filepath.Join(files.BkRoot, util.FILES, name)
		metaPath := filepath.Join(files.BkRoot, util.METADATA, name+".json")
		e1 := os.Remove(metaPath)
		e2 := os.Remove(filePath)
		if err := util.WrapErrors(e1, e2); err != nil {
			return err
		}
	}
	return nil
}

func (files ChangedFiles) syncUpdate() error {
	for _, name := range files.Updated {
		fmt.Print(".")
		if err := overwriteMetadata(name, files.BkRoot, files.MainRoot); err != nil {
			return err
		}
	}
	return nil
}

func overwriteMetadata(name, bkRoot, mainRoot string) error {
	src := filepath.Join(mainRoot, util.METADATA, name+".json")
	dst := filepath.Join(bkRoot, util.METADATA, name+".json")
	return util.CopyFile(dst, src)
}

func (files ChangedFiles) syncOverwrite() error {
	for _, name := range files.Overwrited {
		fmt.Print(".")
		if err := overwriteFileAndMeta(name, files.BkRoot, files.MainRoot); err != nil {
			return err
		}
	}
	return nil
}

func overwriteFileAndMeta(name, bkRoot, mainRoot string) error {
	if err := overwriteMetadata(name, bkRoot, mainRoot); err != nil {
		return err
	}
	return overwriteFile(name, bkRoot, mainRoot)
}

func overwriteFile(name, bkRoot, mainRoot string) error {
	src := filepath.Join(mainRoot, util.FILES, name)
	dst := filepath.Join(bkRoot, util.FILES, name)
	return util.CopyFile(dst, src)
}

func (files ChangedFiles) syncAdd() error {
	for _, name := range files.Added {
		fmt.Print(".")
		if err := overwriteFileAndMeta(name, files.BkRoot, files.MainRoot); err != nil {
			return err
		}
	}
	return nil
}

func getChangedFiles(mainRoot, bkRoot string, mainDB, bkDB *bolt.DB) (files ChangedFiles, err error) {
	files.MainRoot = mainRoot
	files.BkRoot = bkRoot

	err = bkDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		return b.ForEach(func(k, v []byte) error {
			bkFile, err := unmarshalFile(v)
			if err != nil {
				return err
			}

			mainFile, err := getFileByID(bkFile.ID, mainDB)
			if err != nil {
				return err
			}

			// 已被刪除的檔案
			if mainFile == nil {
				files.Deleted = append(files.Deleted, bkFile.Filename)
				return nil
			}

			// 更新了內容的檔案
			if bkFile.Checksum != mainFile.Checksum {
				files.Overwrited = append(files.Overwrited, bkFile.Filename)
				return nil
			}

			// 更新了屬性(metadata/json)的檔案
			if bkFile.UTime != mainFile.UTime {
				files.Updated = append(files.Updated, bkFile.Filename)
			}
			return nil
		})
	})
	if err != nil {
		return
	}

	// 新增的檔案
	err = mainDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		return b.ForEach(func(k, v []byte) error {
			mainFile, err := unmarshalFile(v)
			if err != nil {
				return err
			}
			bkFile, err := getFileByID(mainFile.ID, bkDB)
			if err != nil {
				return err
			}
			if bkFile == nil {
				files.Added = append(files.Added, mainFile.Filename)
			}
			return nil
		})
	})
	return
}

// 如果 err == nil && f == nil, 则意味着 id 不存在。
func getFileByID(id string, db *bolt.DB) (f *File, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		data := b.Get([]byte(id))
		if data == nil {
			return nil
		}
		file, e1 := unmarshalFile(data)
		if e1 != nil {
			return e1
		}
		f = &file
		return nil
	})
	return
}

func unmarshalFile(data []byte) (f File, err error) {
	err = json.Unmarshal(data, &f)
	return
}

func autoFix(mainRoot, bkRoot string, mainDB, bkDB *bolt.DB) error {
	if err := autoFixOneWay(mainRoot, bkRoot, mainDB); err != nil {
		return err
	}
	return autoFixOneWay(bkRoot, mainRoot, bkDB)
}

// 從 root1 和 db 中找出受損檔案, 再從 root2 中尋找有用檔案。
// 有用檔案是指與受損檔案對應的完整檔案。
func autoFixOneWay(root1, root2 string, db *bolt.DB) error {
	fcMap, err := util.ReadFileChecked(root1)
	if err != nil {
		return err
	}
	ids := util.DamagedOfFileChecked(fcMap)
	if len(ids) == 0 {
		fmt.Println("無受損檔案 =>", root1)
		return nil
	}
	damagedFiles, err := getFilesByIDs(ids, db)
	if err != nil {
		return err
	}
	changed, err := fixFiles(root1, root2, damagedFiles, fcMap)
	if err != nil {
		return err
	}
	if changed {
		fileCheckedPath := filepath.Join(root1, util.FileCheckedPath)
		fmt.Println("Update =>", fileCheckedPath)
		_, err := util.WriteJSON(fcMap, fileCheckedPath)
		return err
	}
	return nil
}

// files 是 root1 里的受损档案, fcMap 是 root1 的档案检查列表。
// 如果 changed==true, 說明 fcMap 的内容已改變。
func fixFiles(root1, root2 string, files []*File, fcMap map[string]*FileChecked) (changed bool, err error) {
	var fixedIDs []string
	for _, f := range files {
		filepath1 := filepath.Join(root1, util.FILES, f.Filename)
		filepath2 := filepath.Join(root2, util.FILES, f.Filename)
		sum, err := util.FileSum512(filepath2)
		if err != nil {
			return false, err
		}
		if sum != f.Checksum {
			fmt.Println("未修復 =>", filepath1)
			continue
		}
		fmt.Println("發現有用檔案 =>", filepath2)
		fmt.Println("自動修復 =>", filepath1)
		if err = util.CopyFile(filepath1, filepath2); err != nil {
			return false, err
		}
		fixedIDs = append(fixedIDs, f.ID)
		fcMap[f.ID].Damaged = false
		changed = true
	}
	return
}

func getFilesByIDs(ids []string, db *bolt.DB) (files []*File, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		files, err = util.GetFilesByIDs(ids, tx)
		return err
	})
	return
}
