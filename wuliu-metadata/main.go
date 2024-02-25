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
	"strings"
)

type (
	File        = util.File
	FileAndMeta = util.FileAndMeta
	EditFiles   = util.EditFiles
)

var (
	newFlag   = flag.String("newjson", "", "create a JSON file for modifying metadata")
	cfgPath   = flag.String("json", "", "use a JSON file to modify metadata")
	omitempty = flag.Bool("omitempty", true, "ignore empty values, default: true")
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
		err := overwriteMetadata(files)
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

func overwriteMetadata(files []File) error {
	fmt.Println("overwriteMetadata")
	return nil
}

func printMetadata(files []File) {
	fmt.Printf("\n檔案屬性修改預覽:\n")
	fmt.Printf("(尚未實際執行，使用參數 '-danger' 纔會實際執行)\n\n")
	for _, f := range files {
		size := util.FileSizeToString(float64(f.Size), 0)
		size = fmt.Sprintf("(%s)", size)
		size = util.PaddingRight(size, " ", 9)
		fmt.Printf("%s\t%s %s\n", f.ID, size, f.Filename)
		printLike(f.Like)
		printLabel(f.Label)
		printNotes(f.Notes)
		if f.Like != 0 || f.Label+f.Notes != "" {
			fmt.Println()
		}
		printSlice(f.Keywords, "Keywords")
		printSlice(f.Collections, "Collections")
		printSlice(f.Albums, "Albums")
		fmt.Println()
	}
}

func printLike(like int) {
	for i := 0; i < like; i++ {
		fmt.Print("♥")
	}
	if like > 0 {
		fmt.Print(" ")
	}
}

func printLabel(s string) {
	if s == "" {
		return
	}
	fmt.Printf("[%s] ", s)
}

func printNotes(s string) {
	if s == "" {
		return
	}
	fmt.Printf("%s", s)
}

func printSlice(s []string, name string) {
	if len(s) == 0 {
		return
	}
	joined := strings.Join(s, ", ")
	fmt.Printf("%s: %s\n", name, joined)
}

func readConfig(db *bolt.DB) (cfg EditFiles, files []File) {
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
		files = append(files, util.ReadFile(metaPath))
	}
	return
}

func updateFiles(cfg EditFiles, files []File) []File {
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

func updateFilesOmitEmpty(cfg EditFiles, files []File) []File {
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
