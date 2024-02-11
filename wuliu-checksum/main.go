package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"path/filepath"
	"time"
)

type (
	File        = util.File
	FileChecked = util.FileChecked
)

const (
	MB = util.MB
)

var (
	MainProject = util.ReadProjectInfo(".")
)

var (
	renewFlag    = flag.Bool("renew", false, "reset checked time and status of all files")
	projectsFlag = flag.Bool("projects", false, "list all projects")
	nFlag        = flag.Int("n", 0, "select a project by a number (default: 0)")
	checkFlag    = flag.Bool("check", false, "check if files are corrupted")
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	if !(*renewFlag || *projectsFlag || *checkFlag) {
		flag.Usage()
	}

	root := MainProject.Projects[*nFlag]

	if *projectsFlag {
		printProjectsList()
		return
	}

	db := lo.Must(util.OpenDB(root))
	defer db.Close()

	fcMap := lo.Must(util.ReadFileChecked(root))

	if *renewFlag {
		printInfo(root, len(fcMap), db)
		n := renewFileChecked(root, db)
		fmt.Println("renew後待檢查檔案數量:", n)
		return
	}

	if *checkFlag {
		printInfo(root, len(fcMap), db)
		doCheck(root, fcMap, db)
		return
	}
}

func doCheck(root string, fcMap map[string]*FileChecked, db *bolt.DB) {
	checkN, checkedSize := checkChecksum(root, fcMap, db)
	totalSize := util.FileSizeToString(float64(checkedSize), 2)
	fmt.Println("本次檢查檔案數量:", checkN)
	fmt.Println("本次檢查檔案體積:", totalSize)
	printDamaged(fcMap, db)
	if checkN > 0 {
		fileCheckedPath := filepath.Join(root, util.FileCheckedPath)
		fmt.Println("Update =>", fileCheckedPath)
		_ = lo.Must(
			util.WriteJSON(fcMap, fileCheckedPath))
	}
}

func printProjectsList() {
	for i, project := range MainProject.Projects {
		fmt.Printf("%d %s\n", i, project)
	}
}

func printInfo(root string, n int, db *bolt.DB) {
	fmt.Println("已選擇專案:", root)
	fmt.Println("數據庫檔案數量:", bucketKeysCount(db))
	fmt.Println("待檢查檔案數量:", n)
}

func printDamaged(fcMap map[string]*FileChecked, db *bolt.DB) {
	ids := util.DamagedOfFileChecked(fcMap)
	names := lo.Must(util.IdsToNames(ids, db))
	fmt.Println("已損壞的檔案:", len(ids))
	for i := range ids {
		fmt.Println(ids[i], names[i])
	}
}

func bucketKeysCount(db *bolt.DB) (n int) {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		n = b.Stats().KeyN
		return nil
	})
	lo.Must0(err)
	return
}

func allIDs(db *bolt.DB) (ids []string) {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		return b.ForEach(func(id, _ []byte) error {
			ids = append(ids, string(id))
			return nil
		})
	})
	lo.Must0(err)
	return
}

// 注意，該函數運行後, fcMap 的内容也会改变。
func checkChecksum(root string, fcMap map[string]*FileChecked, db *bolt.DB) (checkN int, checkedSize int64) {
	now := util.Now()
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		for id := range fcMap {
			needCheck := isFileNeedCheck(fcMap[id].Checked, MainProject.CheckInterval)
			if needCheck {
				f := lo.Must(util.GetFileByID(id, b))
				fmt.Print(".")
				fcMap[id].Damaged = checkFile(root, f)
				fcMap[id].Checked = now
				checkN += 1
				checkedSize += f.Size
			}
			// checkN > 0 是为了确保至少检查一个档案
			if checkN > 0 && checkedSize > int64(MainProject.CheckSizeLimit*MB) {
				return nil
			}
		}
		fmt.Println()
		return nil
	})
	lo.Must0(err)
	return
}

func checkFile(root string, file File) (damaged bool) {
	fPath := filepath.Join(root, util.FILES, file.Filename)
	sum := lo.Must(util.FileSum512(fPath))
	if sum != file.Checksum {
		damaged = true
	}
	return
}

func isFileNeedCheck(checked string, intervalDay int) bool {
	interval := intervalDay * util.Day
	needCheckUnix := time.Now().Unix() - int64(interval)
	needCheckDate := time.Unix(needCheckUnix, 0).Format(util.RFC3339)
	// 如果上次校验日期小于(早于) needCheckDate, 就需要再次校验。
	return checked < needCheckDate
}

func renewFileChecked(root string, db *bolt.DB) int {
	fileCheckedPath := filepath.Join(root, util.FileCheckedPath)
	fmt.Println("更新 =>", fileCheckedPath)
	if util.PathExists(fileCheckedPath) {
		log.Fatalln("File Exitst:", fileCheckedPath)
	}
	ids := allIDs(db)
	m := make(map[string]*FileChecked)
	for _, id := range ids {
		fc := &FileChecked{ID: id, Checked: util.Epoch, Damaged: false}
		m[id] = fc
	}
	_ = lo.Must(
		util.WriteJSON(m, fileCheckedPath))
	return len(ids)
}
