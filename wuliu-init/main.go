package main

import (
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	"log"
	"os"
	"fmt"
)

const NormalDirPerm = 0750

var Folders = []string{
	"files", "metadata", "input", "output", "webpages", "recyclebin",
}

func main() {
	checkCWD()
	makeFolders()
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
		lo.Must0(os.Mkdir(folder, NormalDirPerm))
	}
}
