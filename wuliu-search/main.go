package main

import (
	"cmp"
	"flag"
	"fmt"
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
	labelFlag   = flag.String("label", "", "search by label")
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
	matchMode := ""
	pattern := ""

	if *labelFlag != "" {
		mode = "Label"
		pattern = *labelFlag
		files, matchMode = lo.Must2(searchByLabel(*labelFlag, *matchFlag, db))
	} else if *kwFlag != "" {
		mode = "Keyword"
		pattern = *kwFlag
		files, matchMode = lo.Must2(searchByKeyword(*kwFlag, *matchFlag, db))
	} else if *collFlag != "" {
		mode = "Collection"
		pattern = *collFlag
		files, matchMode = lo.Must2(searchByCollection(*collFlag, *matchFlag, db))
	} else if *albumFlag != "" {
		mode = "Album"
		pattern = *albumFlag
		files, matchMode = lo.Must2(searchByAlbum(*albumFlag, *matchFlag, db))
	}

	files, orderBy := sortFilesLimit(*orderbyFlag, *nFlag, !*ascFlag, files)
	fmt.Printf("\nSearch %s(%s):%s, order by %s, %s\n\n", mode, matchMode, pattern, orderBy, util.AscOrDesc(!*ascFlag))

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

func searchByLabel(pattern, matchMode string, db *bolt.DB) ([]*File, string, error) {
	modes := []string{"exactly", "contains", "suffix"}
	if !slices.Contains(modes, matchMode) {
		matchMode = "prefix"
	}
	files, err := util.GetFilesInBucket(pattern, matchMode, util.LabelBucket, db)
	return files, matchMode, err

}

func searchByKeyword(pattern, matchMode string, db *bolt.DB) ([]*File, string, error) {
	return searchKwCollAlbum(pattern, matchMode, util.KeywordsBucket, db)
}

func searchByCollection(pattern, matchMode string, db *bolt.DB) ([]*File, string, error) {
	return searchKwCollAlbum(pattern, matchMode, util.CollectionsBucket, db)
}

func searchByAlbum(pattern, matchMode string, db *bolt.DB) ([]*File, string, error) {
	return searchKwCollAlbum(pattern, matchMode, util.AlbumsBucket, db)
}

// searchKwCollAlbum search by keyword, collection name or album name.
func searchKwCollAlbum(pattern, matchMode string, bucket []byte, db *bolt.DB) ([]*File, string, error) {
	modes := []string{"prefix", "contains", "suffix"}
	if !slices.Contains(modes, matchMode) {
		matchMode = "exactly"
	}
	files, err := util.GetFilesInBucket(pattern, matchMode, bucket, db)
	return files, matchMode, err
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