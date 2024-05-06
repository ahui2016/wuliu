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
	dumpFlag   = flag.String("dump", "", "all/pics/docs")
	kwFlag    = flag.String("keyword", "", "the keyword to be renamed")
	collFlag  = flag.String("collection", "", "the collection to be renamed")
	albumFlag = flag.String("album", "", "the album to be renamed")
	newNameFlag = flag.String("new-name", "", "a new name for keyword/collection/album")
)

func main() {
	flag.Parse()
	util.MustInWuliu()

	var db *bolt.DB
	db = lo.Must(util.OpenDB("."))
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
	if *renameFlag != "" && !slices.Contains([]string{"keyword", "collection", "album"}, *renameFlag) {
		log.Fatalln("不認識 rename:", *renameFlag)
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
		util.RebuildDatabase(".")
		return
	}

	if *toFlag != "" && *renameFlag == "" {
		log.Fatalln("參數 'to' 必須與 'rename' 同時使用。")
	}
	if *renameFlag != "" && *toFlag == "" {
		log.Fatalln("參數 'reanme' 必須與 'to' 同時使用。")
	}

}

func renameKwCollAlbum(db *bolt.DB) error {
	files, err := searchFiles(
}

func RebuildBucket(kw, coll, album string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		files, err := getAllFiles(tx)
		if err != nil {
			return err
		}
		if kw != "" {
			return util.RebuildKeywordsBucket(files, tx)
		}
		if coll != "" {
			return util.RebuildCollectionsBucket(files, tx)
		}
		return util.RebuildAlbumsBucket(files, tx)
	})
}

func updateDB(metadatas map[string][]byte, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		for id, meta := range metadatas {
			if err := util.PutToBucket([]byte(id), meta, b); err != nil {
				return err
			}
		}
		return nil
	})
}

func updateMetaFiles(files []*File, old, kw, coll, album string) (map[string][]byte, error) {
	metadatas := make(map[string][]byte)
	for _, f := range files {
		meta, err := updateMetaJson(f, old, kw, coll, album)
		if err != nil {
			return metadatas, err
		}
		metadatas[f.ID] = meta
	}
}

// updateMetaJson 更新 metadata 資料夾中的 json 檔案。
// 每次只能更新 kw/coll/album 其中之一。
func updateMetaJson(file *File, old, kw, coll, album string) ([]byte, error) {
	if kw != "" {
		file.Keywords = lo.Replace(file.Keywords, old, kw, 1)
	} else if coll != "" {
		file.Collections = lo.Replace(file.Collections, old, coll, 1)
	} else if album != "" {
		file.Albums = lo.Replace(file.Albums, old, album, 1)
	}
	metaPath := filepath.Join(util.METADATA, file.Filename+".json")
	return util.WriteJSON(file, metaPath)
}

func searchFiles(pattern, bucketName string, db *bolt.DB) ([]*File, error) {
	bucket := getBucket(bucketName)
	return util.GetFilesInBucket(pattern, "exactly", bucket, db)
}

func getBucket(name string) []byte {
	if name == "keyword" {
		return util.KeywordsBucket
	}
	if name == "collection" {
		return util.CollectionsBucket
	}
	if name == "album" {
		return util.AlbumsBucket
	}
	return []byte{}
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
