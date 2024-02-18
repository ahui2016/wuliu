package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ahui2016/wuliu/util"
	"github.com/samber/lo"
	"os"
)

type (
	File = util.File
)

var (
	newFlag = flag.String("newjson", "", "create a JSON file for overwriting files")
	cfgPath = flag.String("json", "", "use a JSON file to overwrite files")
)

func main() {
	flag.Parse()
	util.MustInWuliu()
	db := lo.Must(util.OpenDB("."))
	defer db.Close()

	if *newFlag != "" {
		if err := newJsonFile(); err != nil {
			fmt.Println(err)
		}
		return
	}

	printFiles()
}

func newJsonFile() error {
	if util.PathExists(*newFlag) {
		return fmt.Errorf("file exists: %s", *newFlag)
	}
	*cfgFlag = ""
	files := findFiles()
	_, err := util.WriteJSON(files, *newFlag)
	return err
}

func readConfig() (cfg map[string]string) {
	data := lo.Must(os.ReadFile(*cfgPath))
	lo.Must0(json.Unmarshal(data, &cfg))
	return
}

func findFiles() map[string]string {
	if *cfgPath != "" {
		return readConfig()
	}
	names := lo.Must(util.NamesInBuffer())
	config := make(map[string]string)
	for _, name := range names {
		filetype := util.TypeByFilename(name)
		target := filetypeToTarget(filetype)
		config[name] = target
	}
	return config
}

func printFiles() {
	files := findFiles()
	if len(files) == 0 {
		if *cfgPath != "" {
			fmt.Println("在指定的 json 中未填寫檔案名稱")
		} else {
			fmt.Println("在buffer資料夾中未發現檔案")
		}
		return
	}
	for name, target := range files {
		fmt.Printf("%s <= buffer/%s\n", target, name)
	}
}

func filetypeToTarget(filetype string) string {
	if filetype == "text/json" {
		return "metadata"
	}
	return "files"
}

func printErrorExit(err error) {
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}
}
