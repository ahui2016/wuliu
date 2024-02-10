package util

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func PrintVersionExit(ok bool) {
	if ok {
		fmt.Println(RepoName)
		fmt.Println(RepoURL)
		fmt.Println("Version: 2024-02-10")
		os.Exit(0)
	}
}

func PrintWhereExit(ok bool) {
	if ok {
		fmt.Println(GetExePath())
		os.Exit(0)
	}
}

func FolderMustEmpty(folder string) {
	if folder == "" || folder == "." {
		folder = GetCwd()
	}
	if lo.Must(DirIsNotEmpty(folder)) {
		log.Fatalln("資料夾不為空:", folder)
	}
}

func MakeFolders(verbose bool) {
	for _, folder := range Folders {
		if verbose {
			fmt.Println("Create folder:", folder)
		}
		lo.Must0(MkdirIfNotExists(folder))
	}
}

func InitFileChecked() {
	_ = lo.Must(WriteJSON([]int{}, FileCheckedPath))
}

func ReadFileChecked(root string) (fcList []*FileChecked, err error) {
	fileCheckedPath := filepath.Join(root, FileCheckedPath)
	if PathNotExists(fileCheckedPath) {
		return
	}
	data := lo.Must(os.ReadFile(fileCheckedPath))
	err = json.Unmarshal(data, &fcList)
	return
}

func ReadProjectInfo(root string) (info ProjectInfo) {
	infoPath := filepath.Join(root, ProjectInfoPath)
	data := lo.Must(os.ReadFile(infoPath))
	lo.Must0(json.Unmarshal(data, &info))
	return
}

func WriteProjectInfo(info ProjectInfo) error {
	_, err := WriteJSON(info, ProjectInfoPath)
	return err
}

func MustInWuliu() {
	if PathNotExists(ProjectInfoPath) {
		log.Fatalln("找不到 project.json")
	}
	info := ReadProjectInfo(".")
	if info.RepoName != RepoName {
		log.Fatalf("RepoName (%s) != '%s'", info.RepoName, RepoName)
	}
}

func CheckNotAllowInBackup() {
	projInfo := ReadProjectInfo(".")
	if err := notAllowInBackup(projInfo.IsBackup); err != nil {
		log.Fatalln(err)
	}
}

func notAllowInBackup(isBackup bool) error {
	if isBackup {
		return fmt.Errorf("這是備份專案, 不可使用該功能")
	}
	return nil
}

func FindOrphans() (fileOrphans, metaOrphans []string, err error) {
	files, e1 := namesInFiles()
	metas, e2 := namesInMetadataTrim()
	if err = WrapErrors(e1, e2); err != nil {
		return
	}
	fileOrphans, metaOrphans = lo.Difference(files, metas)
	return
}

// NewFilesFromInput 把档案名 (names) 转换为 File, 此时假设档案都在 input 资料夹内。
func NewFilesFromInput(names []string) (files []*File, err error) {
	for _, name := range names {
		filePath := filepath.Join(INPUT, name)
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
	f := filepath.Join(FILES, name)
	m := filepath.Join(METADATA, name+".json")
	for _, oldpath := range []string{f, m} {
		if PathNotExists(oldpath) {
			fmt.Println("NotFound =>", oldpath)
		} else {
			newpath := filepath.Join(RECYCLEBIN, filepath.Base(oldpath))
			fmt.Println("move =>", newpath)
			if err := os.Rename(oldpath, newpath); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func deleteFilesByName(names []string, db *bolt.DB) error {
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

func AddToFileChecked(files []*File) error {
	fcList, err := ReadFileChecked(".")
	if err != nil {
		return err
	}
	for _, f := range files {
		fc := &FileChecked{f.ID, f.CTime, false}
		fcList = append(fcList, fc)
	}
	_, err = WriteJSON(fcList, FileCheckedPath)
	return err
}

func DeleteFromFileChecked(ids []string) error {
	oldList, err := ReadFileChecked(".")
	if err != nil {
		return err
	}
	fcList := slices.DeleteFunc(oldList, func(fc *FileChecked) bool {
		return slices.Contains(ids, fc.ID)
	})
	_, err = WriteJSON(fcList, FileCheckedPath)
	return err
}
