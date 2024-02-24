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
	newFlag = flag.String("newjson", "", "create a JSON file for modifying metadata")
	cfgPath = flag.String("json", "", "use a JSON file to modify metadata")
	danger  = flag.Bool("danger", false, "really do modifay metadata")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	util.CheckNotAllowInBackup()
	db := lo.Must(util.OpenDB("."))
	defer db.Close()
}

func newJsonFile() error {
	if util.PathExists(*newFlag) {
		return fmt.Errorf("file exists: %s", *newFlag)
	}
	v := util.NewEditFiles([]int{}, []int{})
	_, err := util.WriteJSON(v, *newFlag)
	return err
}

func readConfig(db *bolt.DB) (cfg EditFiles, files []*Files) {
	data := lo.Must(os.ReadFile(*cfgPath))
	err := json.Unmarshal(data, &cfg)
	util.PrintErrorExit(err)
	if len(cfg.Filenames) > 0 {
		log.Fatalln("批量修改檔案屬性時不可通過 Filenames 指定檔案")
	}
	filenames, err := util.IdsToNames(cfg.IDs, db)
	util.PrintErrorExit(err)

	for _, name := range cfg.Filenames {
		metaPath := filepath.Join(util.METADATA, name+".json")
		if util.PathNotExists(dst) {
			log.Fatalln("Warning! 找不到", metaPath)
		}
		files = append(files, util.ReadFile(metaPath))
	}
	return
}
