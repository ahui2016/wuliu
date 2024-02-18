package main

import (
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
)

type (
	File = util.File
)

var (
	fileFlag = flag.String("file", "", "specify a file ID and export the file")
	metaFlag = flag.String("meta", "", "specify a file ID and export the file's metadata(json)")
	idFlag   = flag.String("id", "", "specify a file ID and export the file and its metadata")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	printFiles()
}

func findFiles() []*File {
	names := lo.Must(util.NamesInBuffer())
	return lo.Must(util.NewFilesFrom(names, util.BUFFER))
}

func printFiles() {
	names := lo.Must(util.NamesInBuffer())
	if len(names) == 0 {
		fmt.Println("在buffer資料夾中未發現檔案")
		return
	}
	for _, name := range names {
		filetype := util.TypeByFilename(name)
		target := filetypeToTarget(filetype)
		fmt.Printf("%s <= buffer/%s\n", target, name)
	}
}

func filetypeToTarget(filetype string) string {
	if filetype == "text/json" {
		return "metadata"
	}
	return "files"
}
