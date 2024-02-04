package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"time"
)

type (
	File = util.File
	FileChecked = util.FileChecked
)

const (
	MB = util.MB
)

var (
	FileCheckedPath string
	MainProject     = util.ReadProjectInfo()
)

var (
	renewFlag = flag.Bool("renew", false, "reset checked time of all files")
	listFlag  = flag.Bool("list", false, "list all projects")
	nFlag     = flag.Int("n", 0, "select a project by a number")
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	root := MainProject.Projects[*nFlag]
	FileCheckedPath = root + util.FileCheckedPath
	dbPath := root + util.DatabasePath

	if *listFlag {
		printProjectsList()
		return
	}

	db := lo.Must(openDB(dbPath))
	defer db.Close()

	if *renewFlag {
		fcList := lo.Must(readFileChecked())
		printInfo(root, fcList, db)
		renew(db)
		fmt.Println("renew後待檢查檔案數量:", len(fcList))
		return
	}

	flag.Usage()
}

func printProjectsList() {
	for i, project := range MainProject.Projects {
		fmt.Printf("%d %s\n", i, project)
	}
}

func printInfo(root string, fcList []*FileChecked, db *bolt.DB) {
	fmt.Println("已選擇專案:", root)
	fmt.Println("數據庫檔案數量:", bucketKeysCount(db))
	fmt.Println("待檢查檔案數量:", len(fcList))
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

func openDB(dbPath string) (*bolt.DB, error) {
	return bolt.Open(
		dbPath, util.NormalDirPerm, &bolt.Options{Timeout: 1 * time.Second})
}

func readFileChecked() (fcList []*FileChecked, err error) {
	data := lo.Must(os.ReadFile(FileCheckedPath))
	err = json.Unmarshal(data, &fcList)
	return
}

func renew(db *bolt.DB) {
	fmt.Println("更新 =>", FileCheckedPath)
	if util.PathExists(FileCheckedPath) {
		log.Fatalln("File Exitst:", FileCheckedPath)
	}
	ids := allIDs(db)
	var list []*FileChecked
	for _, id := range ids {
		fc := &FileChecked{id, util.Epoch, false}
		list = append(list, fc)
	}
	_ = lo.Must(
		util.WriteJSON(list, FileCheckedPath))
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

// 注意  fcList 的内容也会改变。
func checkChecksum(root string, fcList []*FileChecked, db *bolt.DB) (checkN, checkedSize int) {
	now := util.Now()
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		for i := range fcList {
			needCheck := isFileNeedCheck(fcList[i].Checked, MainProject.CheckInterval)
			if needCheck {
				f := lo.Must(util.GetFileByID(fcList[i].ID, b))
				fmt.Print(".")
				fcList[i].Damaged = checkFile(root, f)
				fcList[i].Checked = now
				checkN += 1
				checkedSise += f.Size
			}
			// checkN > 0 是为了确保至少检查一个档案
			if checkN > 0 && checkedSize > MainProject.CheckSizeLimit*MB {
				return nil
			}
		}
		return nil
	})
	lo.Must0(err)
	return
}

func checkFile(root string, file File) (damaged bool) {
	fPath := filepath.Join(root, util.FILES, file.Filename)
	sum := lo.Must(util.FileSum512(fPath)
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