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
	util.MustInWuliu()
	util.CheckNotAllowInBackup()

	if *cfgPath+*newFlag == "" {
		flag.Usage()
		return
	}
	if *newFlag != "" {
		newJsonFile()
		return
	}

	db := lo.Must(util.OpenDB("."))
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
	lo.Must(
		util.WriteJSON([]string{}, *newFlag))
}

func printConfig(ids []string, db *bolt.DB) {
	if len(ids) > 0 {
		// 如果在数据库中找不到 id, 也不会报错，
		// 因此 ids 与 names 的数量可能不一致，需要以 names 为准。
		names := lo.Must(util.IdsToNames(ids, db))
		printIdAndName(names)
		return
	}
}

func printIdAndName(names []string) {
	for _, name := range names {
		fmt.Println(util.NameToID(name), name)
	}
}

func deleteFiles(ids []string, db *bolt.DB) {
	if len(ids) == 0 {
		return
	}
	lo.Must0(util.DeleteFilesByID(ids, db))
	lo.Must0(util.RebuildCTimeBucket(db))
	lo.Must0(util.DeleteFromFileChecked(ids))
}

func readConfig() (ids []string) {
	data := lo.Must(os.ReadFile(*cfgPath))
	lo.Must0(json.Unmarshal(data, &ids))
	if len(ids) == 0 {
		log.Fatalln("未填寫要刪除的檔案", *cfgPath)
	}
	return
}
