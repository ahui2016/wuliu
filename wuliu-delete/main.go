package main

import (
	"flag"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	"log"
	"os"
)

var (
	newFlag = flag.String("newjson", "", "create a JSON file for deleting files")
)

func main() {
	flag.Parse()
	newJsonFile()
}

func newJsonFile() {
	if *newFlag == "" {
		return
	}
	if util.PathExists(*newFlag) {
		log.Fatalln("file exists:", *newFlag)
	}
	v := util.FilesToDelete{IDs: []string{}, Names: []string{}}
	lo.Must(
		util.WriteJSON(v, *newFlag))
	os.Exit(0)
}
