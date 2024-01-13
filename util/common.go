package util

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"strings"
	"time"
)

const (
	GB      = 1 << 30
	Day     = 24 * 60 * 60
	RFC3339 = "2006-01-02 15:04:05Z07:00"
)

type ProjectInfo struct {
	RepoName         string
	RepoURL          string
	OrphanLastCheck  string // 上次检查孤立档案的时间
	OrphanFilesCount int    // 孤立的档案数量
	OrphanMetaCount  int    // 孤立的 metadata 数量
}

var DefaultWuliuInfo = ProjectInfo{
	RepoName: "Wuliu File Manager",
	RepoURL:  "https://github.com/ahui2016/wuliu",
}

func PrintVersionExit(ok bool) {
	if ok {
		fmt.Println(DefaultWuliuInfo.RepoName)
		fmt.Println(DefaultWuliuInfo.RepoURL)
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

func ReadProjectInfo() (info ProjectInfo) {
	data := lo.Must(os.ReadFile(ProjectInfoPath))
	lo.Must0(json.Unmarshal(data, &info))
	return
}

func WriteProjectInfo(info ProjectInfo) error {
	return WriteJSON(&info, ProjectInfoPath)
}

func MustInWuliu() {
	if PathNotExists(ProjectInfoPath) {
		log.Fatalln("找不到 project.json")
	}
	info := ReadProjectInfo()
	if info.RepoName != DefaultWuliuInfo.RepoName {
		log.Fatalf("RepoName (%s) != '%s'", info.RepoName, DefaultWuliuInfo.RepoName)
	}
}

func FindOrphans() (fileOrphans, metaOrphans []string, err error) {
	files, e1 := namesInFiles()
	metas, e2 := namesInMetadataTrim()
	if err = WrapErrors(e1, e2); err != nil {
		return
	}
	fileOrphans, metaOrphans = lo.Difference(files, metas)
	info := ReadProjectInfo()
	info.OrphanLastCheck = Now()
	info.OrphanFilesCount = len(fileOrphans)
	info.OrphanMetaCount = len(metaOrphans)
	if err = WriteProjectInfo(info); err != nil {
		return nil, nil, err
	}
	return
}

func FindNewFiles() ([]string, error) {
	return namesInInput()
}

func namesInFiles() ([]string, error) {
	return GetFilenamesBase(FILES)
}

func namesInInput() ([]string, error) {
	return GetFilenamesBase(INPUT)
}

func namesInMetadata() ([]string, error) {
	return GetFilenamesBase(METADATA)
}

func namesInMetadataTrim() ([]string, error) {
	names, err := namesInMetadata()
	if err != nil {
		return nil, err
	}
	trimmed := lo.Map(names, func(name string, _ int) string {
		return strings.TrimSuffix(name, ".json")
	})
	return trimmed, nil
}

func Now() string {
	return time.Now().Format(RFC3339)
}
