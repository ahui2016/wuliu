package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"slices"
)

var (
	infoFlag   = flag.String("info", "", "count/size")
	updateFlag = flag.String("update", "", "cache/rebuild")
	dumpFlag   = flag.String("dump", "", "all/pics")
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	var db *bolt.DB
	if *dumpFlag+*infoFlag != "" || *updateFlag == "cache" {
		db = lo.Must(util.OpenDB("."))
		defer db.Close()
	}

	if *dumpFlag != "" {
		err := dump(*dumpFlag, db)
		util.PrintErrorExit(err)
		return
	}

	if *infoFlag != "" && !slices.Contains([]string{"count", "size"}, *infoFlag) {
		log.Fatalln("不認識", *infoFlag)
	}
	if *updateFlag != "" && !slices.Contains([]string{"cache", "rebuild"}, *updateFlag) {
		log.Fatalln("不認識", *updateFlag)
	}
	if *infoFlag+*updateFlag == "" {
		flag.Usage()
		return
	}

	if *infoFlag == "count" {
		lo.Must0(printDatabaseCount(db))
		return
	}
	if *infoFlag == "size" {
		printTotalSize(db)
		return
	}
	if *updateFlag == "cache" {
		lo.Must0(updateCache(db))
		return
	}
	if *updateFlag == "rebuild" {
		util.RebuildDatabase(".")
		return
	}
}

func dump(what string, db *bolt.DB) error {
	if what == "pics" {
		return dumpPics(db)
	}
	if what == "all" {
		return dumpAll(db)
	}
	return fmt.Errorf("Unknown value: %s", what)
}

func dumpAll(db *bolt.DB) error {
	files, err := util.GetAllFiles(db)
	filename := "all.msgp"
	return dumpSelectedFiles(filename, files, err)
}

func dumpPics(db *bolt.DB) error {
	pics, err := util.GetAllPics(db)
	filename := "pics.msgp"
	return dumpSelectedFiles(filename, pics, err)
}

func dumpSelectedFiles(filename string, files []*util.File, err error) error {
	if err != nil {
		return err
	}
	fmt.Println("Write ->", filename)
	return util.WriteMSGP(files, filename)
}

func printDatabaseCount(db *bolt.DB) error {
	fmt.Println("number of keys/value pairs in the database")
	return db.View(func(tx *bolt.Tx) error {
		for _, name := range util.Buckets {
			b := tx.Bucket(name)
			stats := b.Stats()
			fmt.Printf("%s: %d\n", name, stats.KeyN)
		}
		return nil
	})
}

func updateCache(db *bolt.DB) error {
	return util.RebuildSomeBuckets(db)
}

func printTotalSize(db *bolt.DB) {
	fileN, totalSize := lo.Must2(util.DatabaseFilesSize(db))
	totalSizeStr := util.FileSizeToString(float64(totalSize), 2)
	fmt.Printf("Total: %d files, %s\n", fileN, totalSizeStr)
}
