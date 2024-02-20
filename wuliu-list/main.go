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
	ascFlag = flag.Bool("asc", false, "sort in ascending order")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	files := lo.Must(sortByCTime(*nFlag, !*ascFlag, db))
	util.PrintFilesSimple(files)
}

func sortByCTime(limitN int, descending bool, db *bolt.DB) (files []*File, err error) {
	return sortedFiles(util.CTimeBucket, limitN, descending, db)
}

func sortedFiles(bucketName []byte, limitN int, descending bool, db *bolt.DB) (files []*File, err error) {
	var fileIDs []string
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if descending {
			fileIDs, err = getIdsDescending(limitN, b)
		}
		if err != nil {
			return err
		}
		files, err = util.GetFilesByIDs(fileIDs, tx)
		return err
	})
	return
}

// 假設該 bucket 中的每個 key 都對應多個 fileID.
// 並且假設每個 fileID 只能對應一個 key.
func getIdsDescending(limitN int, b *bolt.Bucket) (fileIDs []string, err error) {
	c := b.Cursor()
	n := 0
	for k, v := c.Last(); k != nil; k, v = c.Prev() {
		var ids []string
		if err := json.Unmarshal(v, &ids); err != nil {
			return nil, err
		}
		// 假設每個 fileID 只能對應一個 key, 因此 fileIDs 裏沒有重複項，不需要去重處理。
		fileIDs = append(fileIDs, ids...)
		n += len(ids)
		if n >= limitN {
			break
		}
	}
	if len(fileIDs) > limitN {
		fileIDs = fileIDs[:limitN]
	}
	return
}
