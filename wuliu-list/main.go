package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"slices"
	"strconv"
)

type File = util.File

var (
	nFlag       = flag.Int("n", 15, "default: 15")
	moreFlag    = flag.Bool("more", false, "show more information")
	ascFlag     = flag.Bool("asc", false, "sort in ascending order")
	orderbyFlag = flag.String("orderby", "ctime", "size/like/utime")
	labelFlag   = flag.Bool("labels", false, "print all labels")
	notesFlag   = flag.Bool("notes", false, "print all notes")
	kwFlag      = flag.Bool("keywords", false, "print all keywords")
	collFlag    = flag.Bool("collections", false, "print all collections")
	albumFlag   = flag.Bool("albums", false, "print all albums")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	if *labelFlag {
		printLabels(db)
		return
	}
	if *notesFlag {
		printNotes(db)
		return
	}
	if *kwFlag {
		printKeywords(db)
		return
	}
	if *collFlag {
		printCollections(db)
		return
	}
	if *albumFlag {
		printAlbums(db)
		return
	}

	files := lo.Must(sortBy(*orderbyFlag, *nFlag, !*ascFlag, db))
	if *moreFlag {
		util.PrintFilesMore(files)
		return
	}
	util.PrintFilesSimple(files)
}

func sortBy(orderby string, limitN int, descending bool, db *bolt.DB) (files []*File, err error) {
	bucketName := bucketNameFrom(orderby)
	fmt.Printf("\n檔案排序依據: %s, %s\n\n", bucketName, ascOrDesc(descending))
	return sortedFiles(bucketName, limitN, descending, db)
}

func ascOrDesc(descending bool) string {
	if descending {
		return "descending"
	}
	return "ascending"
}

func bucketNameFrom(orderby string) []byte {
	switch orderby {
	case "size":
		return util.SizeBucket
	case "like":
		return util.LikeBucket
	case "utime":
		return util.UTimeBucket
	default:
		return util.CTimeBucket
	}
}

func sortedFiles(bucketName []byte, limitN int, descending bool, db *bolt.DB) (files []*File, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		fileIDs, err := sortedIDs(tx, bucketName, limitN, descending)
		if err != nil {
			return err
		}
		files, err = util.GetFilesByIDs(fileIDs, tx)
		return err
	})
	return
}

func sortedIDs(tx *bolt.Tx, bucketName []byte, limitN int, descending bool) (fileIDs []string, err error) {
	b := tx.Bucket(bucketName)
	if string(bucketName) == string(util.SizeBucket) {
		return orderBySizeLimit(limitN, descending, b)
	}
	if descending {
		return getIdsDescending(limitN, b)
	}
	return getIdsAscending(limitN, b)
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

// 假設該 bucket 中的每個 key 都對應多個 fileID.
// 並且假設每個 fileID 只能對應一個 key.
func getIdsAscending(limitN int, b *bolt.Bucket) (fileIDs []string, err error) {
	c := b.Cursor()
	n := 0
	for k, v := c.First(); k != nil; k, v = c.Next() {
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

func orderBySizeLimit(limitN int, descending bool, sizeBuc *bolt.Bucket) (fileIDs []string, err error) {
	sizeToIds, err := sizeOfFiles(sizeBuc)
	if err != nil {
		return nil, err
	}
	keys := orderBySize(descending, sizeToIds)
	n := 0
	for _, size := range keys {
		ids := sizeToIds[size]
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

// 對 sizeToIds 的 keys 進行排序, 並返回 keys.
func orderBySize(descending bool, sizeToIds map[int64][]string) []int64 {
	keys := lo.Keys(sizeToIds)
	slices.Sort(keys)
	if !descending {
		return keys
	}
	slices.Reverse(keys)
	return keys
}

func sizeOfFiles(sizeBuc *bolt.Bucket) (map[int64][]string, error) {
	sizeToIds := make(map[int64][]string)
	err := sizeBuc.ForEach(func(k, v []byte) error {
		size, err := strconv.ParseInt(string(k), 10, 64)
		if err != nil {
			return err
		}
		var ids []string
		if err := json.Unmarshal(v, &ids); err != nil {
			return err
		}
		sizeToIds[size] = ids
		return nil
	})
	return sizeToIds, err
}

func printLabels(db *bolt.DB) {
	printKeysAndLength(util.LabelBucket, db)
}

func printNotes(db *bolt.DB) {
	printKeysAndLength(util.NotesBucket, db)
}

func printKeywords(db *bolt.DB) {
	printKeysAndLength(util.KeywordsBucket, db)
}

func printCollections(db *bolt.DB) {
	printKeysAndLength(util.CollectionsBucket, db)
}

func printAlbums(db *bolt.DB) {
	printKeysAndLength(util.AlbumsBucket, db)
}

func printKeysAndLength(bucketName []byte, db *bolt.DB) {
	keywords, err := util.GetKeysAndIdsLength(bucketName, db)
	util.PrintErrorExit(err)
	fmt.Println()
	if len(keywords) == 0 {
		fmt.Println("(none)")
	}
	for k, n := range keywords {
		fmt.Printf("%s (%d)\n", k, n)
	}
	fmt.Println()
}
