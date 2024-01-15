package util

import (
	"fmt"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"time"
)

const (
	FilesBucket    = []byte("FilesBucket")
	FilenameBucket = []byte("FilenameBucket")
)

var buckets = []string{
	FilesBucket,
	FilenameBucket,
}

func OpenDB() (*bolt.DB, error) {
	return bolt.Open(
		DatabasePath, NormalDirPerm, &bolt.Options{Timeout: 1 * time.Second}))
}

func CreateDatabase() {
	fmt.Println("Create", DatabasePath)
	db := lo.Must(OpenDB())
	defer db.Close()
	lo.Must0(createBuckets(db))
}

func createBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		for _, name := range buckets {
			lo.Must(tx.CreateBucketIfNotExists(name))
		}
		return nil
	})
}

func AddFilesToDB(files []FileAndMeta, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		filesBuc := tx.Bucket(FilesBucket)
		filenameBuc := tx.Bucket(FilenameBucket)
		for _, f := range files {
			if err := filesBuc.Put([]byte(f.ID), f.Metadata); err != nil {
				return err
			}
			if err := filenameBuc.Put([]byte(f.Filename), []byte(f.ID)); err != nil {
				return err
			}
		}
		return nil
	})
}
