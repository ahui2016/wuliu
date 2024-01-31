package main

import (
	"encoding/json"
	"flag"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
)

type File = util.File

var (
	nFlag = flag.Int("n", 15, "default: 15")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	db := lo.Must(util.OpenDB())
	defer db.Close()

	listFiles(*nFlag, db)
}

func listFiles(limitN int, db *bolt.DB) error {
	var fileIDs []string
	return db.View(func(tx *bolt.Tx) error {
		n := 0
		b := tx.Bucket(util.CTimeBucket)
		c := b.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var ids []string
			if err := json.Unmarshal(v, &ids); err != nil {
				return err
			}
			fileIDs = append(fileIDs, ids...)
			n += len(ids)
			if n >= limitN {
				break
			}
		}
		if len(fileIDs) > limitN {
			fileIDs = fileIDs[:limitN]
		}
		files, err := util.GetFilesByIDs(fileIDs, tx)
		if err != nil {
			return err
		}
		util.PrintFilesSimple(files)
		return nil
	})
}
