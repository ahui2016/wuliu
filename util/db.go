package util

import (
	"fmt"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"time"
)

var (
	FilesBucket    = []byte("FilesBucket")
	FilenameBucket = []byte("FilenameBucket")
)

var buckets = [][]byte{
	FilesBucket,
	FilenameBucket,
}

func OpenDB() (*bolt.DB, error) {
	return bolt.Open(
		DatabasePath, NormalDirPerm, &bolt.Options{Timeout: 1 * time.Second})
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
			id := []byte(f.ID)
			filename := []byte(f.Filename)
			if err := PutToBucket(id, f.Metadata, filesBuc); err != nil {
				return err
			}
			if err := PutToBucket(filename, id, filenameBuc); err != nil {
				return err
			}
		}
		return nil
	})
}

func PutToBucket(key []byte, value []byte, b *bolt.Bucket) error {
	if KeyExistsInBucket(key, b) {
		return fmt.Errorf("key exists: %s", key)
	}
	return b.Put(key, value)
}

func DeleteInDB(ids []string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		var filenames []string
		filesBuc := tx.Bucket(FilesBucket)
		for _, id := range ids {
			// 如果找不到 id, 则忽略，不报错。
			if v := filesBuc.Get([]byte(id)); v != nil {
				var f File
				if err := json.Unmarshal(v, &f); err != nil {
					filenames = append(filenames, f.Filename)
				}
			}
		}

		filenameBuc := tx.Bucket(FilenameBucket)
		for _, name := range filenames {
			if err := filenameBuc.Delete([]byte(name)); err != nil {
				return err
			}
		}
		for _, id := range ids {
			if err := filesBuc.Delete([]byte(id)); err != nil {
				return err
			}
		}
		return nil
	})
}

func FilesExistInDB(files []*File, db *bolt.DB) (existFiles []string) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(FilesBucket)
		for _, f := range files {
			if KeyExistsInBucket([]byte(f.ID), b) {
				existFiles = append(existFiles, f.Filename)
			}
		}
		return nil
	})
	return
}

func KeyExistsInBucket(key []byte, b *bolt.Bucket) bool {
	return b.Get(key) != nil
}