package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"time"
)

var vFlag = flag.Bool("v", false, "print the version of the command")

var Folders = []string{
	"files", "metadata", "input", "output", "webpages", "recyclebin",
}

func main() {
	flag.Parse()
	util.PrintVersionExit(*vFlag)

	checkCWD()
	makeFolders()
	writeProjectInfo()
	createDatabase()
}

func checkCWD() {
	cwd := util.GetCwd()
	if lo.Must(util.DirIsNotEmpty(cwd)) {
		log.Fatalln("當前目錄不為空:", cwd)
	}
}

func makeFolders() {
	for _, folder := range Folders {
		fmt.Println("Create folder:", folder)
		lo.Must0(os.Mkdir(folder, util.NormalDirPerm))
	}
}

func writeProjectInfo() {
	fmt.Println("Create", util.ProjectInfoPath)
	lo.Must0(util.WriteJSON(&util.WuliuInfo, util.ProjectInfoPath))
}

func createDatabase() {
	fmt.Println("Create", util.DatabasePath)
	db, err := bolt.Open(util.DatabasePath, util.NormalDirPerm, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
}
