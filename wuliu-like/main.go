package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"path/filepath"
)

type (
	File = util.File
)

var (
	idFlag = flag.String("id", "", "which file to like or unlike")
	nFlag  = flag.Int("n", 1, "how much do you like it")
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

	err = like(*nFlag, file, db)
	util.PrintErrorExit(err)
}

func requireIdFlag(id string) error {
	if id == "" {
		return fmt.Errorf("require 'id' flag")
	}
	return nil
}

func like(n int, file File, db *bolt.DB) error {
	if file.Like == n {
		fmt.Printf("❤️=%d [%s] %s\n", file.Like, file.ID, file.Filename)
		fmt.Println("無變化")
		return nil
	}
	file.Like = n
	if err := updateMetadata(file, db); err != nil {
		return err
	}
	fmt.Printf("❤️=%d [%s] %s\n", file.Like, file.ID, file.Filename)
	fmt.Println("UTime =", file.UTime)
	return nil
}

func updateMetadata(file File, db *bolt.DB) error {
	file.UTime = util.Now()
	metaPath := filepath.Join(util.METADATA, file.Filename+".json")
	data, err := util.WriteJSON(file, metaPath)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		return b.Put([]byte(file.ID), data)
	})
}
