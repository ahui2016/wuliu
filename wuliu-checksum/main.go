package main

import (
	"flag"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
)

type (
	FileChecked = util.FileChecked
)

var (
	FileCheckedPath = util.FileCheckedPath
	renewFlag       = flag.Bool("renew", false, "reset checked time of all files")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	db := lo.Must(util.OpenDB())
	defer db.Close()

	if *renewFlag {
		renew(db)
		return
	}

	flag.Usage()
}

func renew(db *bolt.DB) {
	if util.PathExists(FileCheckedPath) {
		log.Fatalln("File Exitst:", FileCheckedPath)
	}
	ids := allIDs(db)
	var list []FileChecked
	for _, id := range ids {
		fc := FileChecked{id, util.Epoch, false}
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
