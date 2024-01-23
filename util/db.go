package util

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"time"
)

var (
	FilesBucket    = []byte("FilesBucket")
	FilenameBucket = []byte("FilenameBucket")
	ChecksumBucket = []byte("ChecksumBucket")
	SizeBucket = []byte("SizeBucket")
	TypeBucket = []byte("TypeBucket")
	LikeBucket = []byte("LikeBucket")
	LabelBucket = []byte("LabelBucket")
	NotesBucket = []byte("NotesBucket")
	KeywordsBucket = []byte("KeywordsBucket")
	CollectionsBucket = []byte("CollectionsBucket")
	AlbumsBucket = []byte("AlbumsBucket")
	CTimeBucket = []byte("CTimeBucket")
	UTimeBucket = []byte("UTimeBucket")
	CheckedBucket = []byte("CheckedBucket")
	DamagedBucket = []byte("DamagedBucket")
)

var buckets = [][]byte{
	FilesBucket,
	FilenameBucket,
	ChecksumBucket,
	SizeBucket,
	TypeBucket,
	LikeBucket,
	LabelBucket,
	NotesBucket,
	KeywordsBucket,
	CollectionsBucket,
	AlbumsBucket,
	CTimeBucket,
	UTimeBucket,
	CheckedBucket,
	DamagedBucket,
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

func bucketPutJson(k string, v any, b *bolt.Bucket) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.Put([]byte(k), data)
}

func bucketGetStrSlice(key string, b *bolt.Bucket) ([]string, error) {
	data := b.Get([]byte(key))
	if data == nil {
		return nil, nil
	}
	strSlice := []string{}
	err := json.Unmarshal(data, &strSlice)
	return strSlice, err
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

func putChecksum(f *File, b *bolt.Bucket) error {
	ids := bucketGetStrSlice(f.Checksum, b)
	if ids != nil {
		ids = append(ids, f.ID)
		return bucketPutJson(f.Checksum, ids, b)
	}
	return bucketPutJson(f.Checksum, []string{f.ID}, b)
}

func rebuildChecksumBucket(files []*File, tx *bolt.Tx) error {
	b := tx.Bucket(ChecksumBucket)
	for _, f := range files {
		
	}
}

func PutToBucket(key []byte, value []byte, b *bolt.Bucket) error {
	if KeyExistsInBucket(key, b) {
		return fmt.Errorf("key exists: %s", key)
	}
	return b.Put(key, value)
}

func DeleteInDB(ids []string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		filesBuc := tx.Bucket(FilesBucket)
		filenames, err := idsToNames(ids, filesBuc)
		if err != nil {
			return err
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

func idsToNames(ids []string, filesBuc *bolt.Bucket) (names []string, err error) {
	for _, id := range ids {
		// 如果找不到 id, 则忽略，不报错。
		v := filesBuc.Get([]byte(id))
		if v == nil {
			continue
		}
		var f File
		if err := json.Unmarshal(v, &f); err != nil {
			return nil, err
		}
		names = append(names, f.Filename)
	}
	return
}

func IdsToNames(ids []string, db *bolt.DB) (names []string, err error) {
	dbErr := db.View(func(tx *bolt.Tx) error {
		filesBuc := tx.Bucket(FilesBucket)
		names, err = idsToNames(ids, filesBuc)
		return err
	})
	return names, dbErr
}

func FilesExistInDB(files []*File, db *bolt.DB) (existFiles []*File) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(FilesBucket)
		for _, f := range files {
			if KeyExistsInBucket([]byte(f.ID), b) {
				existFiles = append(existFiles, f)
			}
		}
		return nil
	})
	return
}

func KeyExistsInBucket(key []byte, b *bolt.Bucket) bool {
	return b.Get(key) != nil
}
