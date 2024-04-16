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
		fmt.Println(RepoName)
		fmt.Println(RepoURL)
		fmt.Println("Version: 2024-04-16")
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

func ExecutableDir() string {
	return filepath.Dir(GetExePath())
}

func InitFileChecked() {
	m := make(map[int]int)
	fmt.Println("Create", FileCheckedPath)
	_ = lo.Must(WriteJSON(m, FileCheckedPath))
}

func ReadFileChecked(root string) (fcMap map[string]*FileChecked, err error) {
	fileCheckedPath := filepath.Join(root, FileCheckedPath)
	if PathNotExists(fileCheckedPath) {
		return
	}
	data := lo.Must(os.ReadFile(fileCheckedPath))
	err = json.Unmarshal(data, &fcMap)
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
	files, e1 := NamesInFiles()
	metas, e2 := namesInMetadataTrim()
	if err = WrapErrors(e1, e2); err != nil {
		return
	}
	fileOrphans, metaOrphans = lo.Difference(files, metas)
	return
}

// NewFilesFrom 把档案名 (names) 转换为 files, 此时假设档案在 folder 资料夹内。
func NewFilesFrom(names []string, folder string) (files []*File, err error) {
	for _, name := range names {
		filePath := filepath.Join(folder, name)
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
		f.Type = TypeByFilename(name)
		f.Keywords = []string{}
		f.Collections = []string{}
		f.Albums = []string{}
		files = append(files, f)
	}
	return
}

func NamesInFiles() ([]string, error) {
	return GetFilenamesBase(FILES)
}

func NamesInBuffer() ([]string, error) {
	return GetFilenamesBase(BUFFER)
}

func NamesInInput() ([]string, error) {
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
		size = fmt.Sprintf("(%s)", size)
		size = PaddingRight(size, " ", 9)
		fmt.Printf("%s\t%s %s\n", f.ID, size, f.Filename)
	}
	fmt.Println()
}

func PrintFilesMore(files []*File) {
	for _, f := range files {
		size := FileSizeToString(float64(f.Size), 0)
		size = fmt.Sprintf("(%s)", size)
		size = PaddingRight(size, " ", 9)
		fmt.Printf("%s\t%s %s\n", f.ID, size, f.Filename)
		printLike(f.Like)
		printLabel(f.Label)
		printNotes(f.Notes)
		if f.Like != 0 || f.Label+f.Notes != "" {
			fmt.Println()
		}
		printSlice(f.Keywords, "Keywords")
		printSlice(f.Collections, "Collections")
		printSlice(f.Albums, "Albums")
		fmt.Println()
	}
}

func printLike(like int) {
	for i := 0; i < like; i++ {
		fmt.Print("♥")
	}
	if like > 0 {
		fmt.Print(" ")
	}
}

func printLabel(s string) {
	if s == "" {
		return
	}
	fmt.Printf("[%s] ", s)
}

func printNotes(s string) {
	if s == "" {
		return
	}
	fmt.Printf("%s", s)
}

func printSlice(s []string, name string) {
	if len(s) == 0 {
		return
	}
	joined := strings.Join(s, ", ")
	fmt.Printf("%s: %s\n", name, joined)
}

func PaddingRight(s, char string, length int) string {
	for {
		if len(s) >= length {
			return s
		}
		s += char
	}
}

func AddToFileChecked(files []*File) error {
	fcMap, err := ReadFileChecked(".")
	if err != nil {
		return err
	}
	for _, f := range files {
		fc := &FileChecked{f.ID, f.CTime, false}
		fcMap[f.ID] = fc
	}
	_, err = WriteJSON(fcMap, FileCheckedPath)
	return err
}

func DeleteFromFileChecked(ids []string) error {
	fcMap, err := ReadFileChecked(".")
	if err != nil {
		return err
	}
	for _, id := range ids {
		delete(fcMap, id)
	}
	_, err = WriteJSON(fcMap, FileCheckedPath)
	return err
}

func DamagedOfFileChecked(fcMap map[string]*FileChecked) (ids []string) {
	for _, fc := range fcMap {
		if fc.Damaged == true {
			ids = append(ids, fc.ID)
		}
	}
	return
}
