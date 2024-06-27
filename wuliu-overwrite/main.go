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
	"path/filepath"
	"slices"
)

type (
	File = util.File
)

var (
	newFlag = flag.String("newjson", "", "create a JSON file for overwriting files")
	cfgPath = flag.String("json", "", "use a JSON file to overwrite files")
	danger  = flag.Bool("danger", false, "really do overwrite files")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	util.CheckNotAllowInBackup()

	if *newFlag != "" {
		if err := newJsonFile(); err != nil {
			fmt.Println(err)
		}
		return
	}

	files := findFiles()
	checkFilesEmpty(files)

	if *danger {
		db := lo.Must(util.OpenDB("."))
		defer db.Close()
		err := overwriteFiles(files, db)
		util.PrintErrorExit(err)
	} else {
		printFiles(files)
	}
}

func newJsonFile() error {
	if util.PathExists(*newFlag) {
		return fmt.Errorf("file exists: %s", *newFlag)
	}
	*cfgPath = ""
	files := findFiles()
	_, err := util.WriteJSON(files, *newFlag)
	return err
}

func readConfig() (cfg map[string]string) {
	data := lo.Must(os.ReadFile(*cfgPath))
	lo.Must0(json.Unmarshal(data, &cfg))
	return
}

func findFiles() map[string]string {
	if *cfgPath != "" {
		return readConfig()
	}
	names := lo.Must(util.NamesInBuffer())
	config := make(map[string]string)
	for _, name := range names {
		filetype := util.TypeByFilename(name)
		target := filetypeToTarget(filetype)
		config[name] = target
	}
	return config
}

// pairs 是指檔案及其屬性同時存在，可以湊成一對。
func splitFiles(files map[string]string) (pairs, normalFiles map[string]string) {
	// 暫時不需要實現該函數，等遇到檔案讀寫衝突再說。
	return
}

func printFiles(files map[string]string) {
	for name, target := range files {
		fmt.Printf("%s <= buffer/%s\n", target, name)
		lo.Must0(checkTarget(target))
	}
}

func overwriteFiles(files map[string]string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		filesBuc := tx.Bucket(util.FilesBucket)
		for name, target := range files {
			if err := overwriteFile(name, target, filesBuc); err != nil {
				return err
			}
		}
		return nil
	})
}

func overwriteFile(name, target string, b *bolt.Bucket) error {
	fmt.Printf("%s <= buffer/%s\n", target, name)
	if err := checkTarget(target); err != nil {
		fmt.Println("Warning!", err)
		return nil
	}
	src := filepath.Join(util.BUFFER, name)
	dst := filepath.Join(target, name)
	if util.PathNotExists(dst) {
		fmt.Println("Warning! 找不到", dst)
		return nil
	}
	if target == util.FILES {
		return overwriteIntoFiles(name, src, dst, b)
	}
	if target == util.METADATA {
		return overwriteIntoMetadata(src, dst, b)
	}
	return nil
}

func overwriteIntoFiles(name, src, dst string, b *bolt.Bucket) error {
	metaPath := filepath.Join(util.METADATA, name+".json")
	f := util.ReadFile(metaPath)

	f.UTime = util.Now()
	sum, err := util.FileSum512(src)
	if err != nil {
		return err
	}
	if f.Checksum == sum {
		fmt.Println("檔案內容沒有變化:", name)
		return nil
	}
	f.Checksum = sum

	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	f.Size = info.Size()

	if err = os.Rename(src, dst); err != nil {
		return err
	}
	data, err := util.WriteJSON(f, metaPath)
	if err != nil {
		return err
	}
	return b.Put([]byte(f.ID), data)
}

func overwriteIntoMetadata(src, dst string, b *bolt.Bucket) error {
	f := util.ReadFile(src)
	old := util.ReadFile(dst)

	if (f.Label+f.Notes+f.CTime == old.Label+old.Notes+old.CTime) &&
		f.Like == old.Like &&
		slices.Equal(f.Keywords, old.Keywords) &&
		slices.Equal(f.Collections, old.Collections) &&
		slices.Equal(f.Albums, old.Albums) {
		fmt.Println("檔案屬性沒有變化:", old.Filename+".json")
		return nil
	}

	f.ID = old.ID
	f.Filename = old.Filename
	f.Checksum = old.Checksum
	f.Size = old.Size
	f.Type = old.Type
	f.UTime = util.Now()

	data, err := util.WriteJSON(f, dst)
	if err != nil {
		return err
	}
	if err = b.Put([]byte(f.ID), data); err != nil {
		return err
	}
	return os.Remove(src)
}

func readOldFile(name, target string) File {
	filePath := filepath.Join(util.METADATA, name)
	if target == util.FILES {
		filePath += ".json"
	}
	return util.ReadFile(filePath)
}

func checkTarget(target string) error {
	if target != util.FILES && target != util.METADATA {
		return fmt.Errorf("不認識目標目錄: %s\n目標目錄只能是 'files' 或 'metadata'", target)
	}
	return nil
}

func checkFilesEmpty(files map[string]string) {
	if len(files) == 0 {
		if *cfgPath != "" {
			log.Fatalln("在指定的 json 中未填寫檔案名稱")
		} else {
			log.Fatalln("在buffer資料夾中未發現檔案")
		}
	}
}

func filetypeToTarget(filetype string) string {
	if filetype == "text/json" {
		return "metadata"
	}
	return "files"
}
