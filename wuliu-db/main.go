package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"slices"
)

var (
	infoFlag   = flag.String("info", "", "count/size")
	updateFlag = flag.String("update", "", "cache/rebuild")
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	var db *bolt.DB
	if *infoFlag != "" || *updateFlag == "cache" {
		db = lo.Must(util.OpenDB())
		defer db.Close()
	}

	if *infoFlag != "" && !slices.Contains([]string{"count", "size"}, *infoFlag) {
		log.Fatalln("不認識", *infoFlag)
	}
	if *updateFlag != "" && !slices.Contains([]string{"cache", "rebuild"}, *updateFlag) {
		log.Fatalln("不認識", *updateFlag)
	}
	if *infoFlag+*updateFlag == "" {
		flag.Usage()
		os.Exit(0)
	}

	if *infoFlag == "count" {
		lo.Must0(printDatabaseCount(db))
		return
	}
	if *updateFlag == "cache" {
		lo.Must0(updateCache(db))
		return
	}
	if *updateFlag == "rebuild" {
		util.RebuildDatabase()
		return
	}
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
