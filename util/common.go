package util

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
)

type ProjectInfo struct {
	RepoName string
	RepoURL  string
}

var WuliuInfo = ProjectInfo{
	RepoName: "Wuliu File Manager",
	RepoURL:  "https://github.com/ahui2016/wuliu",
}

func PrintVersionExit(ok bool) {
	if ok {
		fmt.Println(WuliuInfo.RepoName)
		fmt.Println(WuliuInfo.RepoURL)
		fmt.Println("Version: 2024-01-04")
		os.Exit(0)
	}
}

func readProjectInfo() (info ProjectInfo) {
	data := lo.Must(os.ReadFile(ProjectInfoPath))
	lo.Must0(json.Unmarshal(data, &info))
	return
}

func MustInWuliu() {
	if PathNotExists(ProjectInfoPath) {
		log.Fatalln("找不到 project.json")
	}
	info := readProjectInfo()
	if info.RepoName != WuliuInfo.RepoName {
		log.Fatalf("RepoName (%s) != '%s'", info.RepoName, WuliuInfo.RepoName)
	}
}
