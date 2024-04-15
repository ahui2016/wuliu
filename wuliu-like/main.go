package main

type (
	File = util.File
)

var (
	idFlag = flag.String("id", "", "which file to like or unlike")
	unlikeFlag = flag.Bool("unlike", false, "set like to zero")
	nFlag = flag.Int("n", 1, "how much do you like it")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	util.CheckNotAllowInBackup()

	err := requireIdFlag(*idFlag)
	util.PrintErrorExit(err)

	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	file, err := util.GetFileInDB(*idFlag, db)
	util.PrintErrorExit(err)

}

func requireIdFlag(id string) error {
	if id == "" {
		return fmt.Errorf("require 'id' flag")
	}
	return nil
}

func unlike(file File, db *bolt.db) error {
	file.Like = 0
	return updateMetadata(file, db)
}

func like(n int, file File, db *bolt.db) error {
	file.Like = n
	return updateMetadata(file, db)
}

func updateMetadata(file File, db *bolt.DB) error {
	file.UTime = util.Now()
	metaPath := filepath.Join(util.Metadata, file.Name+".json")
	data, err := util.WriteJSON(file, metaPath)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b1 := tx.Bucket(util.FilesBucket)
		e1 := b.Put([]byte(file.ID), data)

		b2 := tx.Bucket(util.LikeBucket)
		e2 = util.PutIntAndIDs(int64(file.Like), file.ID, b2)

		b3 := tx.Bucket(util.UTimeBucket)
		e3 := util.PutStrAndIDs(file.UTime, file.ID, b3)

		return util.WrapErrors(e1, e2, e3)
	})
	return err
}
