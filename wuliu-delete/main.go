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
)

var (
	newFlag = flag.String("newjson", "", "create a JSON file for deleting files")
	cfgPath = flag.String("json", "", "use a JSON file to delete files")
	danger  = flag.Bool("danger", false, "really do delete files")
)

func main() {
	flag.Parse()
	if *cfgPath+*newFlag == "" {
		flag.Usage()
		os.Exit(0)
	}
	if *newFlag != "" {
		newJsonFile()
		os.Exit(0)
	}

	db := lo.Must(util.OpenDB())
	defer db.Close()
	cfg := readConfig()

	if *danger {
		deleteFiles(cfg, db)
	} else {
		printConfig(cfg, db)
	}

}

func newJsonFile() {
	if util.PathExists(*newFlag) {
		log.Fatalln("file exists:", *newFlag)
	}
	v := util.FilesToDelete{IDs: []string{}, Names: []string{}}
	lo.Must(
		util.WriteJSON(v, *newFlag))
}

func printConfig(toDelete util.FilesToDelete, db *bolt.DB) {
	if len(toDelete.IDs) > 0 {
		names := lo.Must(util.IdsToNames(toDelete.IDs, db))
		printIdAndName(names)
		return
	}
	if len(toDelete.Names) > 0 {
		printIdAndName(toDelete.Names)
	}
}

func printIdAndName(names []string) {
	for _, name := range names {
		fmt.Println(util.NameToID(name), name)
	}
}

func deleteFiles(toDelete util.FilesToDelete, db *bolt.DB) {
	if len(toDelete.IDs) > 0 {
		util.DeleteFilesByID(toDelete.IDs, db)
		return
	}
	if len(toDelete.Names) > 0 {
		util.DeleteFilesByName(toDelete.Names, db)
	}
}

func readConfig() (cfg util.FilesToDelete) {
	data := lo.Must(os.ReadFile(*cfgPath))
	lo.Must0(json.Unmarshal(data, &cfg))
	if err := cfg.Check(); err != nil {
		log.Fatalln(err)
	}
	if len(cfg.IDs)+len(cfg.Names) == 0 {
		log.Fatalln(*cfgPath, "未填寫要刪除的檔案")
	}
	return
}
