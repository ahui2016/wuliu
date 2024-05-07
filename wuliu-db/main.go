package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"path/filepath"
	"slices"
)

type File = util.File

var (
	infoFlag    = flag.String("info", "", "count/size")
	updateFlag  = flag.String("update", "", "cache/rebuild")
	dumpFlag    = flag.String("dump", "", "all/pics/docs")
	kwFlag      = flag.String("keyword", "", "the keyword to be renamed")
	collFlag    = flag.String("collection", "", "the collection to be renamed")
	albumFlag   = flag.String("album", "", "the album to be renamed")
	newNameFlag = flag.String("rename-to", "", "a new name for keyword/collection/album")
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	if *infoFlag != "" && !slices.Contains([]string{"count", "size"}, *infoFlag) {
		log.Fatalln("不認識 info:", *infoFlag)
	}
	if *updateFlag != "" && !slices.Contains([]string{"cache", "rebuild"}, *updateFlag) {
		log.Fatalln("不認識 update:", *updateFlag)
	}
	if *dumpFlag != "" && !slices.Contains([]string{"all", "pics", "docs"}, *dumpFlag) {
		log.Fatalln("不認識 dump:", *dumpFlag)
	}

	if *dumpFlag != "" {
		err := dump(*dumpFlag, db)
		util.PrintErrorExit(err)
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
		db.Close()
		util.RebuildDatabase(".")
		return
	}

	if *newNameFlag != "" && *kwFlag+*collFlag+*albumFlag == "" {
		log.Fatalln("參數 '--rename-to' 必須與 keyword/collection/album 其中之一同時使用。")
	}
	if *kwFlag != "" && *newNameFlag == "" {
		log.Fatalln("參數 '-keyword' 必須與 '--rename-to' 同時使用。")
	}
	if *collFlag != "" && *newNameFlag == "" {
		log.Fatalln("參數 '-collection' 必須與 '--rename-to' 同時使用。")
	}
	if *albumFlag != "" && *newNameFlag == "" {
		log.Fatalln("參數 '-album' 必須與 '--rename-to' 同時使用。")
	}
	err := renameKwCollAlbum(*kwFlag, *collFlag, *albumFlag, *newNameFlag, db)
	util.PrintErrorExit(err)
}

func renameKwCollAlbum(kw, coll, album, newName string, db *bolt.DB) error {
	files, err := searchFiles(kw, coll, album, db)
	if err != nil {
		return err
	}

	metadatas, err := updateMetaFiles(files, kw, coll, album, newName)
	if err != nil {
		return err
	}

	fmt.Println("Update database...")
	if err = updateFilesBucket(metadatas, db); err != nil {
		return err
	}
	return rebuildBucket(kw, coll, album, db)
}

func rebuildBucket(kw, coll, album string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		files, err := util.GetAllFilesTx(tx)
		if err != nil {
			return err
		}
		if kw != "" {
			return util.RebuildKeywordsBucket(files, tx)
		}
		if coll != "" {
			return util.RebuildCollBucket(files, tx)
		}
		return util.RebuildAlbumsBucket(files, tx)
	})
}

func updateFilesBucket(metadatas map[string][]byte, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		for id, meta := range metadatas {
			if err := b.Put([]byte(id), meta); err != nil {
				return err
			}
		}
		return nil
	})
}

func updateMetaFiles(files []*File, kw, coll, album, newName string) (map[string][]byte, error) {
	metadatas := make(map[string][]byte)
	for _, f := range files {
		meta, err := updateMetaJson(f, kw, coll, album, newName)
		if err != nil {
			return metadatas, err
		}
		metadatas[f.ID] = meta
	}
	return metadatas, nil
}

// updateMetaJson 更新 metadata 資料夾中的 json 檔案。
// 每次只能更新 kw/coll/album 其中之一。
func updateMetaJson(file *File, kw, coll, album, newName string) ([]byte, error) {
	if kw != "" {
		file.Keywords = lo.Replace(file.Keywords, kw, newName, 1)
	} else if coll != "" {
		file.Collections = lo.Replace(file.Collections, coll, newName, 1)
	} else if album != "" {
		file.Albums = lo.Replace(file.Albums, album, newName, 1)
	}
	metaPath := filepath.Join(util.METADATA, file.Filename+".json")
	fmt.Println("Update =>", metaPath)
	return util.WriteJSON(file, metaPath)
}

func searchFiles(kw, coll, album string, db *bolt.DB) ([]*File, error) {
	bucket, pattern := getBucketPattern(kw, coll, album)
	return util.GetFilesInBucket(pattern, "exactly", bucket, db)
}

func getBucketPattern(kw, coll, album string) ([]byte, string) {
	if kw != "" {
		return util.KeywordsBucket, kw
	}
	if coll != "" {
		return util.CollectionsBucket, coll
	}
	if album != "" {
		return util.AlbumsBucket, album
	}
	return []byte{}, ""
}

func dump(what string, db *bolt.DB) error {
	if what == "docs" {
		return dumpDocs(db)
	}
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

func dumpDocs(db *bolt.DB) error {
	docs, err := util.GetAllDocs(db)
	filename := "docs.msgp"
	return dumpSelectedFiles(filename, docs, err)
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
