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
	"slices"
)

type (
	File        = util.File
	FileAndMeta = util.FileAndMeta
	EditFiles   = util.EditFiles
)

var (
	newFlag = flag.String("newjson", "", "create a JSON file for adding files")
	cfgPath = flag.String("json", "", "use a JSON file to add files")
	danger  = flag.Bool("danger", false, "really do add files")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	util.CheckNotAllowInBackup()

	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	files := findNewFiles()
	checkExist(files, db)

	if *newFlag != "" {
		newJsonFile()
		return
	}

	if *danger {
		addNewFiles(files, db)
	} else {
		printNewFiles(files)
	}
}

func newJsonFile() {
	if util.PathExists(*newFlag) {
		log.Fatalln("file exists:", *newFlag)
	}
	names := lo.Must(util.NamesInInput())
	v := util.NewEditFiles([]string{}, names)
	lo.Must(
		util.WriteJSON(v, *newFlag))
}

func readConfig() (cfg EditFiles) {
	data := lo.Must(os.ReadFile(*cfgPath))
	err := json.Unmarshal(data, &cfg)
	util.PrintErrorExit(err)
	if len(cfg.IDs) > 0 {
		log.Fatalln("添加新檔案時不可通過 ID 指定檔案")
	}
	return
}

func findNewFiles() []*File {
	inputNames := lo.Must(util.NamesInInput())
	if *cfgPath == "" {
		return lo.Must(util.NewFilesFrom(inputNames, util.INPUT))
	}
	cfg := readConfig()
	if len(cfg.Filenames) == 0 {
		cfg.Filenames = inputNames
	}
	var filenames []string
	for _, name := range cfg.Filenames {
		if slices.Contains(inputNames, name) {
			filenames = append(filenames, name)
		} else {
			fmt.Println("Not Found:", name)
		}
	}
	files := lo.Must(util.NewFilesFrom(filenames, util.INPUT))
	for i := range files {
		files[i].Like = cfg.Like
		files[i].Label = cfg.Label
		files[i].Notes = cfg.Notes
		files[i].Keywords = cfg.Keywords
		files[i].Collections = cfg.Collections
		files[i].Albums = cfg.Albums
	}
	return files
}

func printNewFiles(files []*File) {
	if len(files) == 0 {
		fmt.Println("在input資料夾中未發現新檔案")
		return
	}
	if *cfgPath != "" {
		data := lo.Must(os.ReadFile(*cfgPath))
		fmt.Printf("%s\n", data)
		return
	}
	for _, f := range files {
		size := util.FileSizeToString(float64(f.Size), 2)
		size = fmt.Sprintf("(%s)", size)
		size = util.PaddingRight(size, " ", 11)
		fmt.Printf("%s %s\n", size, f.Filename)
	}
}

func addNewFiles(files []*File, db *bolt.DB) {
	if len(files) == 0 {
		fmt.Println("warning: No file to add.")
		return
	}
	var metadatas []FileAndMeta
	for _, f := range files {
		metaPath := filepath.Join(util.METADATA, f.Filename+".json")
		fmt.Println("Create =>", metaPath)
		meta := lo.Must(util.WriteJSON(f, metaPath))
		metadatas = append(metadatas, FileAndMeta{f, meta})

		src := filepath.Join(util.INPUT, f.Filename)
		dst := filepath.Join(util.FILES, f.Filename)
		fmt.Println("Add =>", dst)
		lo.Must0(os.Rename(src, dst))
	}
	fmt.Println("Update database...")
	lo.Must0(util.AddFilesToDB(metadatas, db))
	lo.Must0(util.RebuildCTimeBucket(db))
	lo.Must0(util.AddToFileChecked(files))
	fmt.Println("OK")
}

func checkExist(files []*File, db *bolt.DB) {
	existInDB := util.FilesExistInDB(files, db)
	if len(existInDB) > 0 {
		fmt.Println("【注意！】數據庫中有同名檔案：")
		printIdAndName(existInDB)
		os.Exit(0)
	}

	var dstFiles []string
	for _, f := range files {
		dst := filepath.Join(util.FILES, f.Filename)
		meta := filepath.Join(util.METADATA, f.Filename+".json")
		dstFiles = append(dstFiles, dst, meta)
	}
	var existFiles []string
	for _, f := range dstFiles {
		if util.PathExists(f) {
			existFiles = append(existFiles, f)
		}
	}
	if len(existFiles) > 0 {
		fmt.Println("【注意！】同名檔案已存在：")
		util.PrintList(existFiles)
		os.Exit(0)
	}
}

func printIdAndName(files []*File) {
	for _, f := range files {
		fmt.Printf("%s: %s\n", f.ID, f.Filename)
	}
}
