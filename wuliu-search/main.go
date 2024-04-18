package main

import (
	"fmt"
	"cmp"
	"flag"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"slices"
)

type File = util.File

var (
	nFlag       = flag.Int("n", 15, "default: 15")
	moreFlag    = flag.Bool("more", false, "show more information")
	ascFlag     = flag.Bool("asc", false, "sort in ascending order")
	orderbyFlag = flag.String("orderby", "ctime", "ctime/utime/filename")
	matchFlag   = flag.String("match", "", "exactly/prefix/contains/suffix")
	kwFlag      = flag.String("keyword", "", "search by a keyword")
	collFlag    = flag.String("collection", "", "search by a collection name")
	albumFlag   = flag.String("album", "", "search by a album name")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	files := []*File{}
	mode := ""
	pattern := ""
	if *kwFlag != "" {
		mode = "Keyword"
		pattern = *kwFlag
		files = lo.Must(searchByKeyword(*kwFlag, *matchFlag, db))
	} else if *collFlag != "" {
		mode = "Collection"
		pattern = *collFlag
		files = lo.Must(searchByCollection(*collFlag, *matchFlag, db))
	}

	files, orderBy := sortFilesLimit(*orderbyFlag, *nFlag, !*ascFlag, files)
	fmt.Printf("\nSearching %s:%s, order by %s, %s\n\n", mode, pattern, orderBy, util.AscOrDesc(!*ascFlag))

	if len(files) == 0 {
		fmt.Println("找不到符合條件的檔案。")
		return
	}

	if *moreFlag {
		util.PrintFilesMore(files)
		return
	}
	util.PrintFilesSimple(files)
}

func searchByKeyword(pattern, matchMode string, db *bolt.DB) (files []*File, err error) {
	return searchKwCollAlbum(pattern, matchMode, util.KeywordsBucket, db)
}

func searchByCollection(pattern, matchMode string, db *bolt.DB) (files []*File, err error) {
	return searchKwCollAlbum(pattern, matchMode, util.CollectionsBucket, db)
}

// searchKwCollAlbum search by keyword, collection name or album name.
func searchKwCollAlbum(pattern, matchMode string, bucket []byte, db *bolt.DB) (files []*File, err error) {
	if matchMode == "prefix" {
		return
	}
	if matchMode == "contains" {
		return
	}
	if matchMode == "suffix" {
		return
	}
	return util.GetFilesInBucket(pattern, bucket, db)
}

func sortFilesLimit(orderBy string, n int, desc bool, files []*File) ([]*File, string) {
	if orderBy == "filename" {
		return files, orderBy
	}
	if orderBy == "utime" {
		return files, orderBy
	}
	files = orderByCTimeLimit(n, desc, files)
	return files, "ctime"
}

func orderByCTimeLimit(n int, desc bool, files []*File) []*File {
	if desc {
		slices.SortFunc(files, func(a, b *File) int {
			return cmp.Compare(b.CTime, a.CTime)
		})
	} else {
		slices.SortFunc(files, func(a, b *File) int {
			return cmp.Compare(a.CTime, b.CTime)
		})
	}
	if len(files) > n {
		files = files[:n]
	}
	return files
}
