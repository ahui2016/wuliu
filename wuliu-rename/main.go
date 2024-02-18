package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	idFlag   = flag.String("id", "", "specify a file ID that already exists")
	nameFlag = flag.String("name", "", "set a new filename")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	if *idFlag == "" && *nameFlag != "" {
		log.Fatalln("Required '-id'")
	}
	if *idFlag != "" && *nameFlag == "" {
		log.Fatalln("Required '-name'")
	}
	if *idFlag != "" && *nameFlag != "" {
		err := checkFilename(*nameFlag)
		printErrorExit(err)

		file, err := util.GetFileInDB(*idFlag, db)
		printErrorExit(err)

		fm, err := renameMeta(file.Filename, *nameFlag)
		printErrorExit(err)

		err = renameFile(file.Filename, *nameFlag)
		printErrorExit(err)

		fmt.Println("Update database...")
		err = renameInDB(*idFlag, fm, db)
		printErrorExit(err)
		fmt.Println("OK")

		return
	}
	flag.Usage()
}

func renameMeta(oldname, newname string) (fm util.FileAndMeta, err error) {
	src := filepath.Join(util.METADATA, oldname+".json")
	dst := filepath.Join(util.METADATA, newname+".json")
	fmt.Printf("Rename: %s => %s\n", src, dst)
	if err = checkExists(src, dst); err != nil {
		return
	}
	file := util.ReadFile(src)
	file.Filename = newname
	file.ID = util.NameToID(newname)
	meta, err := util.WriteJSON(file, dst)
	fm.File = &file
	fm.Metadata = meta

	rcBin := filepath.Join(util.RECYCLEBIN, oldname+".json")
	err = os.Rename(src, rcBin)
	return
}

func renameFile(oldname, newname string) error {
	src := filepath.Join(util.FILES, oldname)
	dst := filepath.Join(util.FILES, newname)
	fmt.Printf("Rename: %s => %s\n", src, dst)
	if err := checkExists(src, dst); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

func checkFilename(name string) (err error) {
	if strings.ContainsAny(name, `\/:*?"<>|`) {
		err = fmt.Errorf(`檔案名稱不允許包含這些字符 \/:*?"<>|`)
	}
	return
}

func checkExists(src, dst string) error {
	if util.PathNotExists(src) {
		return fmt.Errorf("not found: %s", src)
	}
	if util.PathExists(dst) {
		return fmt.Errorf("file exists: %s", dst)
	}
	return nil
}

func renameInDB(oldID string, newfile util.FileAndMeta, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		filesBuc := tx.Bucket(util.FilesBucket)
		if err := filesBuc.Delete([]byte(oldID)); err != nil {
			return err
		}
		if err := util.PutToBucket([]byte(newfile.ID), newfile.Metadata, filesBuc); err != nil {
			return err
		}

		return util.RenameCTime(newfile.CTime, oldID, newfile.ID, tx)
	})
}

func printErrorExit(err error) {
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}
}
