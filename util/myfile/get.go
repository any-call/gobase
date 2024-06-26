package myfile

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func ListDir(dir string) ([]string, error) {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(infos))
	for i, info := range infos {
		names[i] = info.Name()
	}
	return names, nil
}

// ListFileNames lists all file names in the directory.
func ListFileNames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for i := 0; i < len(entries); i++ {
		if !entries[i].IsDir() {
			names = append(names, entries[i].Name())
		}
	}
	return names, nil
}

// FileMD5 gets file MD5。
func FileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return FileMD5Reader(file)
}

// FileMD5Reader gets file MD5 from io.Reader.
func FileMD5Reader(r io.Reader) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, r)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// FileSize gets file size in bytes.
func FileSize(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return FileSizeFile(file)
}

// FileSizeFile gets file size from os.File in bytes.
func FileSizeFile(file *os.File) (int64, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// ListDirEntryPaths lists all the file or directory paths in the directory recursively.
// If the cur is true result will include current directory.
// Note that GetDirAllEntryPaths won't follow symlink if the subdir is a symbolic link.
func ListDirEntryPaths(dirname string, cur bool) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))

	// Include current directory.
	if cur {
		paths = append(paths, dirname)
	}

	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := ListDirEntryPaths(path, cur)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}

// ListDirEntryPathsSymlink lists all the file or dir paths in the directory recursively.
// If the cur is true result will include current directory.
func ListDirEntryPathsSymlink(dirname string, cur bool) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))

	// Include current directory.
	if cur {
		paths = append(paths, dirname)
	}

	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		realInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if realInfo.IsDir() {
			tmp, err := ListDirEntryPathsSymlink(path, cur)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}
