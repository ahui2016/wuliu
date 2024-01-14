package util

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"strings"
)

func PrintVersionExit(ok bool) {
	if ok {
		fmt.Println(DefaultWuliuInfo.RepoName)
		fmt.Println(DefaultWuliuInfo.RepoURL)
		fmt.Println("Version: 2024-01-13")
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

// NewFilesFromInput 把档案名 (names) 转换为 File, 此时假设档案都在 input 资料夹内。
func NewFilesFromInput(names []string) (files []*File, err error) {
	for _, name := range names {
		filePath := INPUT + "/" + name
		info, err := os.Lstat(filePath)
		if err != nil {
			return nil, err
		}
		checksum, err := FileSum512(filePath)
		if err != nil {
			return nil, err
		}
		f := NewFile(name)
		f.Checksum = checksum
		f.Size = info.Size()
		f.Type = typeByFilename(name)
		files = append(files, f)
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
