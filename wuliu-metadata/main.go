package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
)

type (
	File        = util.File
	FileAndMeta = util.FileAndMeta
	EditFiles   = util.EditFiles
)

var (
	newFlag   = flag.String("newjson", "", "create a JSON file for modifying metadata")
	cfgPath   = flag.String("json", "", "use a JSON file to modify metadata")
	omitempty = flag.Bool("omitempty", true, "ignore empty values")
	danger    = flag.Bool("danger", false, "really do modifay metadata")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	util.CheckNotAllowInBackup()

	if *cfgPath+*newFlag == "" {
		flag.Usage()
		return
	}
	if *newFlag != "" {
		if err := newJsonFile(); err != nil {
			fmt.Println(err)
		}
		return
	}

	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	cfg, files := readConfig(db)
	files = updateFiles(cfg, files)

	if *danger {
		err := overwriteMetadata(files, db)
		util.PrintErrorExit(err)
	} else {
		printMetadata(files)
	}
}

func newJsonFile() error {
	if util.PathExists(*newFlag) {
		return fmt.Errorf("file exists: %s", *newFlag)
	}
	v := util.NewEditFiles([]string{}, []string{})
	_, err := util.WriteJSON(v, *newFlag)
	return err
}

func overwriteMetadata(files []*File, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(util.FilesBucket)
		for _, f := range files {
			metaPath := filepath.Join(util.METADATA, f.Filename+".json")
			fmt.Println("Update =>", metaPath)
			if util.PathNotExists(metaPath) {
				fmt.Println("Warning! 找不到", metaPath)
			}
			data, err := util.WriteJSON(f, metaPath)
			if err != nil {
				return err
			}
			if err = b.Put([]byte(f.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func printMetadata(files []*File) {
	fmt.Printf("\n檔案屬性修改預覽:\n")
	fmt.Printf("(尚未實際執行，使用參數 '-danger' 纔會實際執行)\n\n")
	util.PrintFilesMore(files)
}

func readConfig(db *bolt.DB) (cfg EditFiles, files []*File) {
	data := lo.Must(os.ReadFile(*cfgPath))
	err := json.Unmarshal(data, &cfg)
	util.PrintErrorExit(err)
	if len(cfg.Filenames) > 0 {
		log.Fatalln("批量修改檔案屬性時不可通過 Filenames 指定檔案")
	}
	filenames, err := util.IdsToNames(cfg.IDs, db)
	util.PrintErrorExit(err)

	for _, name := range filenames {
		metaPath := filepath.Join(util.METADATA, name+".json")
		if util.PathNotExists(metaPath) {
			log.Fatalln("Warning! 找不到", metaPath)
		}
		f := util.ReadFile(metaPath)
		files = append(files, &f)
	}
	return
}

func updateFiles(cfg EditFiles, files []*File) []*File {
	if *omitempty {
		return updateFilesOmitEmpty(cfg, files)
	}
	now := util.Now()
	for i := range files {
		files[i].Like = cfg.Like
		files[i].Label = cfg.Label
		files[i].Notes = cfg.Notes
		files[i].Keywords = cfg.Keywords
		files[i].Collections = cfg.Collections
		files[i].Albums = cfg.Albums
		files[i].UTime = now
	}
	return files
}

func updateFilesOmitEmpty(cfg EditFiles, files []*File) []*File {
	now := util.Now()
	for i := range files {
		if cfg.Like != 0 {
			files[i].Like = cfg.Like
		}
		if cfg.Label != "" {
			files[i].Label = cfg.Label
		}
		if cfg.Notes != "" {
			files[i].Notes = cfg.Notes
		}
		if len(cfg.Keywords) != 0 {
			files[i].Keywords = cfg.Keywords
		}
		if len(cfg.Collections) != 0 {
			files[i].Collections = cfg.Collections
		}
		if len(cfg.Albums) != 0 {
			files[i].Albums = cfg.Albums
		}
		files[i].UTime = now
	}
	return files
}
