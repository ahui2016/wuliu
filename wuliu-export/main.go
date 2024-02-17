package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"path/filepath"
)

const (
	MB = util.MB
)

var (
	fileFlag = flag.String("file", "", "specify a file ID and export the file")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	projInfo := util.ReadProjectInfo(".")

	if *fileFlag != "" {
		if err := exportFile(*fileFlag, db, projInfo); err != nil {
			fmt.Println("Error!", err)
		}
	}

}

func exportFile(id string, db *bolt.DB, info util.ProjectInfo) error {
	f, err := getFileByID(id, db)
	if err != nil {
		return err
	}
	if err = checkSizeLimit(f.Size, info); err != nil {
		return err
	}
	src := filepath.Join(util.FILES, f.Filename)
	dst := filepath.Join(util.BUFFER, f.Filename)
	fmt.Println("Export =>", dst)
	if util.PathExists(dst) {
		return fmt.Errorf("file exists: %s", dst)
	}
	return util.CopyFile(dst, src)
}

func getFileByID(id string, db *bolt.DB) (f util.File, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		f, err = util.GetFileByID(id, b)
		return err
	})
	return
}

func checkSizeLimit(size int64, info util.ProjectInfo) error {
	limit := info.ExportSizeLimit * MB
	if size > limit {
		sizeStr := util.FileSizeToString(float64(size), 2)
		return fmt.Errorf("檔案體積(%s) 超過上限(%d MB)", sizeStr, info.ExportSizeLimit)
	}
	return nil
}
