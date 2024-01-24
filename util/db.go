package util

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"time"
)

var (
	FilesBucket       = []byte("FilesBucket")
	FilenameBucket    = []byte("FilenameBucket")
	ChecksumBucket    = []byte("ChecksumBucket")
	SizeBucket        = []byte("SizeBucket")
	TypeBucket        = []byte("TypeBucket")
	LikeBucket        = []byte("LikeBucket")
	LabelBucket       = []byte("LabelBucket")
	NotesBucket       = []byte("NotesBucket")
	KeywordsBucket    = []byte("KeywordsBucket")
	CollectionsBucket = []byte("CollectionsBucket")
	AlbumsBucket      = []byte("AlbumsBucket")
	CTimeBucket       = []byte("CTimeBucket")
	UTimeBucket       = []byte("UTimeBucket")
	CheckedBucket     = []byte("CheckedBucket")
	DamagedBucket     = []byte("DamagedBucket")
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

func PutToBucket(key []byte, value []byte, b *bolt.Bucket) error {
	if KeyExistsInBucket(key, b) {
		return fmt.Errorf("key exists: %s", key)
	}
	return b.Put(key, value)
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

func RebuildAllBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		if err := rebuildFilesAndName(tx); err != nil {
			return err
		}
		files, e1 := getAllFiles(tx)
		e2 := rebuildSomeBuckets(files, tx)
		return WrapErrors(e1, e2)
	})
}

func RebuildSomeBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		files, err := getAllFiles(tx)
		if err != nil {
			return err
		}
		return rebuildSomeBuckets(files, tx)
	})
}

func rebuildFilesAndName(tx *bolt.Tx) error {
	filesBuc := reCreateBucket(FilesBuckets)
	nameBuc := reCreateBucket(FilenameBucket)
	files, err := getAllFilesMetadata()
	if err != nil {
		return err
	}
	for _, f := range files {
		e1 := bucketPutJson(f.ID, f, filesBuc)
		e2 := b.Put([]byte(f.Filename), []byte(f.ID))
		if err := WrapErrors(e1, e2); err != nil {
			return err
		}
	}
	return nil
}

func getAllFilesMetadata() ([]*File, error) {
	metaPaths, err := getAllMetadataPaths()
	if err != nil {
		return nil, err
	}
	return metaPathsToFiles(metaPaths)
}

func metaPathsToFiles(paths []string) (files []*File, err error) {
	for _, meta := range paths {
		data, err := os.ReadFile(meta)
		if err != nil {
			return nil, err
		}
		var f File
		if err := json.Unmarshal(data, &f); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	return files, nil
}

func getAllMetadataPaths() ([]string, error) {
	a, b, err := FindOrphans()
	if err != nil {
		return nil, err
	}
	if len(a)+len(b) > 0 {
		return nil, fmt.Errorf("發現孤立檔案，請執行 wuliu-orphan")
	}
	return filepath.Glob(filepath.Join(METADATA, "/*"))
}

func rebuildSomeBuckets(files []*File, tx *bolt.Tx) error {
	csumBuc := reCreateBucket(ChecksumBucket)
	sizeBuc := reCreateBucket(SizeBucket)
	typeBuc := reCreateBucket(TypeBucket)
	likeBuc := reCreateBucket(LikeBucket)
	labelBuc := reCreateBucket(LabelBucket)
	notesBuc := reCreateBucket(NotesBucket)
	kwBuc := reCreateBucket(KeywordsBucket)
	collBuc := reCreateBucket(CollectionsBucket)
	albumBuc := reCreateBucket(AlbumsBucket)
	ctimeBuc := reCreateBucket(CTimeBucket)
	utimeBuc := reCreateBucket(UTimeBucket)
	checkBuc := reCreateBucket(CheckedBucket)
	dmgBuc := reCreateBucket(DamagedBucket)

	for _, f := range files {
		e1 := putStrAndID(f.Checksum, f.ID, csumBuc)
		e2 := putIntAndID(f.Size, f.ID, sizeBuc)
		e3 := putStrAndID(f.Type, f.ID, typeBuc)
		e4 := putIntAndID(f.Like, f.ID, likeBuc)
		e5 := putStrAndID(f.Label, f.ID, labelBuc)
		e6 := putStrAndID(f.Notes, f.ID, notesBuc)
		e7 := putSliceAndID(f.Keywords, f.ID, kwBuc)
		e8 := putSliceAndID(f.Collections, f.ID, collBuc)
		e9 := putSliceAndID(f.Albums, f.ID, albumBuc)
		e10 := putStrAndID(f.CTime, f.ID, ctimeBuc)
		e11 := putStrAndID(f.UTime, f.ID, utimeBuc)
		e12 := putStrAndID(f.Checked, f.ID, checkBuc)
		e13 := putIdAndBool(f.ID, f.Damaged, dmgBuc)

		if err := WrapErrors(e1, e2, e3, e4, e5, e6,
			e7, e8, e9, e10, e11, e12, e13); err != nil {
			return err
		}
	}
	return nil
}

func GetAllFiles(db *bolt.DB) (files []*File, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		files, e := getAllFiles(tx)
		if e != nil {
			return e
		}
	})
	return
}

func getAllFiles(tx *bolt.Tx) (files []*File, err error) {
	b := tx.Bucket(FilesBucket)
	err = b.ForEach(func(_, v []byte) error {
		var f File
		if err := json.Unmarshal(v, &f); err != nil {
			return nil, err
		}
		files = append(files, &f)
	})
	return
}

func reCreateBucket(name []byte, tx *bolt.Tx) (*bolt.Bucket, error) {
	if err := tx.DeleteBucket(name); err != nil {
		return nil, err
	}
	return tx.CreateBucket(name)
}

func putStrAndID(key, id string, b *bolt.Bucket) error {
	ids := bucketGetStrSlice(key, b)
	if ids != nil {
		ids = append(ids, id)
		return bucketPutJson(key, ids, b)
	}
	return bucketPutJson(key, []string{id}, b)
}

func putIntAndID(i int64, id string, b *bolt.Bucket) error {
	key := strconv.FormatInt(i, 10)
	return putKeyAndID(key, f.ID, b)

}

func putSliceAndID(s []string, id string, b *bolt.Bucket) error {
	for _, item := range s {
		if err := putStrAndID(item, id, b); err != nil {
			return err
		}
	}
	return nil
}

func putIdAndBool(id string, v bool, b *bolt.Bucket) error {
	if v {
		return b.Put([]byte(id), []byte("true"))
	}
	return nil
}
