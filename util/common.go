package util

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"strings"
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
		fmt.Println("Version: 2024-01-07")
		os.Exit(0)
	}
}

func PrintWhereExit(ok bool) {
	if ok {
		fmt.Println(GetExePath())
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

func filesInFiles() ([]string, error) {
	return GetFilenamesBase(FILES)
}

func filesInInput() ([]string, error) {
	return GetFilenamesBase(INPUT)
}

func filesInMetadata() ([]string, error) {
	return GetFilenamesBase(METADATA)
}

func filesInMetadataTrim() ([]string, error) {
	names, err := filesInMetadata()
	if err != nil {
		return nil, err
	}
	trimmed := lo.Map(names, func(name string, _ int) string {
		return strings.TrimSuffix(name, ".json")
	})
	return trimmed, nil
}
