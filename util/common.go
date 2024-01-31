package util

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func PrintVersionExit(ok bool) {
	if ok {
		fmt.Println(DefaultWuliuInfo.RepoName)
		fmt.Println(DefaultWuliuInfo.RepoURL)
		fmt.Println("Version: 2024-01-30")
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
	_, err := WriteJSON(&info, ProjectInfoPath)
	return err
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
		f.Keywords = []string{}
		f.Collections = []string{}
		f.Albums = []string{}
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

// deleteFileByName 尝试删除档案，例如 name=abc.txt,
// 则尝试删除 files/abc.txt 和 metadata/abc.txt.json。
// 注意，这里说的删除是将档案移动到专案根目录的 recyclebin 中，
// 如果 recyclebin 里有同名档案则直接覆盖。
func deleteFileByName(name string) {
	f := FILES + "/" + name
	m := METADATA + "/" + name + ".json"
	for _, oldpath := range []string{f, m} {
		if PathNotExists(oldpath) {
			fmt.Println("NotFound =>", oldpath)
		} else {
			newpath := RECYCLEBIN + "/" + filepath.Base(oldpath)
			fmt.Println("move =>", newpath)
			if err := os.Rename(oldpath, newpath); err != nil {
				fmt.Println(err)
			}
		}
	}
}

// DeleteFilesByName 尝试删除档案，包括档案本身, metadata 以及数据库条目。
// 注意，这里说的删除是将档案移动到专案根目录的 recyclebin 中，
// 如果 recyclebin 里有同名档案则直接覆盖。
func DeleteFilesByName(names []string, db *bolt.DB) error {
	for _, name := range names {
		deleteFileByName(name)
	}
	ids := NamesToIds(names)
	return DeleteInDB(ids, db)
}

// DeleteFilesByID 尝试删除档案，包括档案本身, metadata 以及数据库条目。
// 找不到 ID 则忽略，不会报错。
// 注意，这里说的删除是将档案移动到专案根目录的 recyclebin 中，
// 如果 recyclebin 里有同名档案则直接覆盖。
func DeleteFilesByID(ids []string, db *bolt.DB) error {
	names, err := IdsToNames(ids, db)
	if err != nil {
		return err
	}
	for _, name := range names {
		deleteFileByName(name)
	}
	return DeleteInDB(ids, db)
}

func PrintFilesSimple(files []*File) {
	for _, f := range files {
		size := FileSizeToString(float64(f.Size), 0)
		fmt.Printf("%s (%s) %s\n", f.ID, size, f.Filename)
	}
}
