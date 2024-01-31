package util

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"golang.org/x/crypto/blake2b"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// WrapErrors 把多个错误合并为一个错误.
func WrapErrors(allErrors ...error) (wrapped error) {
	for _, err := range allErrors {
		if err != nil {
			if wrapped == nil {
				wrapped = err
			} else {
				wrapped = fmt.Errorf("%w | %w", wrapped, err)
			}
		}
	}
	return
}

func GetCwd() string {
	return lo.Must(os.Getwd())
}

// GetExePath returns the path name for the executable
// that started the current process.
func GetExePath() string {
	return lo.Must1(os.Executable())
}

func DirIsEmpty(dirpath string) (ok bool, err error) {
	items, err := filepath.Glob(dirpath + Separator + "*")
	ok = len(items) == 0
	return
}

func DirIsNotEmpty(dirpath string) (ok bool, err error) {
	ok, err = DirIsEmpty(dirpath)
	return !ok, err
}

func PathNotExists(name string) (ok bool) {
	_, err := os.Lstat(name)
	if os.IsNotExist(err) {
		ok = true
		err = nil
	}
	lo.Must0(err)
	return
}

func PathExists(name string) bool {
	return !PathNotExists(name)
}

// MkdirIfNotExists 创建資料夹, 忽略 ErrExist.
// 在 Windows 里, 文件夹的只读属性不起作用, 为了统一行为, 不把資料夹设为只读.
func MkdirIfNotExists(name string) error {
	if PathExists(name) {
		return nil
	}
	return os.Mkdir(name, NormalDirPerm)
}

// WriteFile 写檔案, 使用权限 0666
func WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, NormalFilePerm)
}

// WriteJSON 把 data 转换为漂亮格式的 JSON 并写入檔案 filename 中。
func WriteJSON(data any, filename string) ([]byte, error) {
	dataJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, err
	}
	if err = WriteFile(filename, dataJSON); err != nil {
		return nil, err
	}
	return dataJSON, nil
}

func isRegularFile(name string) (ok bool, err error) {
	info, err := os.Lstat(name)
	if err != nil {
		return
	}
	return info.Mode().IsRegular(), nil
}

// GetFilenamesBase 假设 folder 里全是普通档案，没有资料夹。
func GetFilenamesBase(folder string) ([]string, error) {
	pattern := filepath.Join(folder, "/*")
	names, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	baseNames := lo.Map(names, func(name string, _ int) string {
		return filepath.Base(name)
	})
	return baseNames, nil
}

func getRegularFiles(folder string) (files []string, err error) {
	pattern := filepath.Join(folder, "*")
	items, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	for _, file := range items {
		ok, err := isRegularFile(file)
		if err != nil {
			return nil, err
		}
		if ok {
			files = append(files, file)
		}
	}
	return files, nil
}

func PrintList[T any](list []T) {
	for _, item := range list {
		fmt.Println(item)
	}
}

func PrintListWithSuffix[T any](list []T, suffix string) {
	for _, item := range list {
		fmt.Printf("%v%v", item, suffix)
	}
}

// BLAKE2b is faster than MD5, SHA-1, SHA-2, and SHA-3, on 64-bit x86-64 and ARM architectures.
// https://en.wikipedia.org/wiki/BLAKE_(hash_function)#BLAKE2
// https://blog.min.io/fast-hashing-in-golang-using-blake2/
// https://pkg.go.dev/crypto/sha256#example-New-File
func FileSum512(name string) (HexString, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := lo.Must(blake2b.New512(nil))
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	checksum := h.Sum(nil)
	return hex.EncodeToString(checksum), nil
}

// FileSizeToString 把文件大小转换为方便人类阅读的格式。
// fixed 指定小数点后几位, 设为负数表示不限制小数位。
func FileSizeToString(size float64, fixed int) string {
	format := fmt.Sprintf("%%.%df", fixed)
	if fixed < 0 {
		format = "%f"
	}
	sizeGB := size / 1024 / 1024 / 1024
	if sizeGB < 1 {
		sizeMB := sizeGB * 1024
		if sizeMB < 1 {
			sizeKB := sizeMB * 1024
			format = format + " KB"
			return fmt.Sprintf(format, sizeKB)
		}
		format = format + " MB"
		return fmt.Sprintf(format, sizeMB)
	}
	format = format + " GB"
	return fmt.Sprintf(format, sizeGB)
}

// https://github.com/gofiber/fiber/blob/master/utils/http.go (edited).
func typeByFilename(filename string) (filetype string) {
	ext := filepath.Ext(filename)
	ext = strings.ToLower(ext)
	if len(ext) == 0 {
		return MIMEOctetStream
	}
	if ext[0] == '.' {
		ext = ext[1:]
	}
	filetype = mimeExtensions[ext]
	if len(filetype) == 0 {
		filetype = MIMEOctetStream
	}

	switch ext {
	case "zip", "rar", "7z", "gz", "tar", "bz", "bz2", "xz":
		filetype = "compressed/" + ext
	case "md", "json", "xml", "html", "xhtml", "htm", "atom", "rss", "yaml",
		"js", "ts", "go", "py", "cs", "dart", "rb", "c", "h", "cpp", "rs":
		filetype = "text/" + ext
	case "doc", "docx", "ppt", "pptx", "rtf", "xls", "xlsx":
		filetype = "office/" + ext
	case "epub", "mobi", "azw", "azw3", "djvu":
		filetype = "ebook/" + ext
	}
	return filetype
}

// MIME types were copied from
// https://github.com/gofiber/fiber/blob/master/utils/http.go
// https://github.com/nginx/nginx/blob/master/conf/mime.types
var mimeExtensions = map[string]string{
	"html":    "text/html",
	"htm":     "text/html",
	"shtml":   "text/html",
	"css":     "text/css",
	"xml":     "application/xml",
	"gif":     "image/gif",
	"jpeg":    "image/jpeg",
	"jpg":     "image/jpeg",
	"js":      "text/javascript",
	"atom":    "application/atom+xml",
	"rss":     "application/rss+xml",
	"mml":     "text/mathml",
	"txt":     "text/plain",
	"jad":     "text/vnd.sun.j2me.app-descriptor",
	"wml":     "text/vnd.wap.wml",
	"htc":     "text/x-component",
	"avif":    "image/avif",
	"png":     "image/png",
	"svg":     "image/svg+xml",
	"svgz":    "image/svg+xml",
	"tif":     "image/tiff",
	"tiff":    "image/tiff",
	"wbmp":    "image/vnd.wap.wbmp",
	"webp":    "image/webp",
	"ico":     "image/x-icon",
	"jng":     "image/x-jng",
	"bmp":     "image/x-ms-bmp",
	"woff":    "font/woff",
	"woff2":   "font/woff2",
	"jar":     "application/java-archive",
	"war":     "application/java-archive",
	"ear":     "application/java-archive",
	"json":    "application/json",
	"hqx":     "application/mac-binhex40",
	"doc":     "application/msword",
	"pdf":     "application/pdf",
	"ps":      "application/postscript",
	"eps":     "application/postscript",
	"ai":      "application/postscript",
	"rtf":     "application/rtf",
	"m3u8":    "application/vnd.apple.mpegurl",
	"kml":     "application/vnd.google-earth.kml+xml",
	"kmz":     "application/vnd.google-earth.kmz",
	"xls":     "application/vnd.ms-excel",
	"eot":     "application/vnd.ms-fontobject",
	"ppt":     "application/vnd.ms-powerpoint",
	"odg":     "application/vnd.oasis.opendocument.graphics",
	"odp":     "application/vnd.oasis.opendocument.presentation",
	"ods":     "application/vnd.oasis.opendocument.spreadsheet",
	"odt":     "application/vnd.oasis.opendocument.text",
	"pptx":    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"xlsx":    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"docx":    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"wmlc":    "application/vnd.wap.wmlc",
	"wasm":    "application/wasm",
	"7z":      "application/x-7z-compressed",
	"cco":     "application/x-cocoa",
	"jardiff": "application/x-java-archive-diff",
	"jnlp":    "application/x-java-jnlp-file",
	"run":     "application/x-makeself",
	"pl":      "application/x-perl",
	"pm":      "application/x-perl",
	"prc":     "application/x-pilot",
	"pdb":     "application/x-pilot",
	"rar":     "application/x-rar-compressed",
	"rpm":     "application/x-redhat-package-manager",
	"sea":     "application/x-sea",
	"swf":     "application/x-shockwave-flash",
	"sit":     "application/x-stuffit",
	"tcl":     "application/x-tcl",
	"tk":      "application/x-tcl",
	"der":     "application/x-x509-ca-cert",
	"pem":     "application/x-x509-ca-cert",
	"crt":     "application/x-x509-ca-cert",
	"xpi":     "application/x-xpinstall",
	"xhtml":   "application/xhtml+xml",
	"xspf":    "application/xspf+xml",
	"zip":     "application/zip",
	"bin":     "application/octet-stream",
	"exe":     "application/octet-stream",
	"dll":     "application/octet-stream",
	"deb":     "application/octet-stream",
	"dmg":     "application/octet-stream",
	"iso":     "application/octet-stream",
	"img":     "application/octet-stream",
	"msi":     "application/octet-stream",
	"msp":     "application/octet-stream",
	"msm":     "application/octet-stream",
	"mid":     "audio/midi",
	"midi":    "audio/midi",
	"kar":     "audio/midi",
	"mp3":     "audio/mpeg",
	"ogg":     "audio/ogg",
	"m4a":     "audio/x-m4a",
	"ra":      "audio/x-realaudio",
	"3gpp":    "video/3gpp",
	"3gp":     "video/3gpp",
	"ts":      "video/mp2t",
	"mp4":     "video/mp4",
	"mpeg":    "video/mpeg",
	"mpg":     "video/mpeg",
	"mov":     "video/quicktime",
	"webm":    "video/webm",
	"flv":     "video/x-flv",
	"m4v":     "video/x-m4v",
	"mng":     "video/x-mng",
	"asx":     "video/x-ms-asf",
	"asf":     "video/x-ms-asf",
	"wmv":     "video/x-ms-wmv",
	"avi":     "video/x-msvideo",
}
