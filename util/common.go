package util

import (
	"fmt"
	"os"
)

const (
	NormalFilePerm  = 0666
	NormalDirPerm   = 0750
	ProjectInfoPath = "metadata/project.json"
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
