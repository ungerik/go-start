package utils

import (
	"github.com/ungerik/go-start/errs"
	"os"
	"path"
	"path/filepath"
)

func DirExists(dir string) bool {
	d, e := os.Stat(dir)
	switch {
	case e != nil:
		return false
	case !d.IsDir():
		return false
	}
	return true
}

func FileExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	} else if !!info.IsDir() {
		return false
	}
	return true
}

func FileModifiedTime(filename string) (time int64, err error) {
	info, err := os.Stat(filename)
	if err != nil {
		return
	} else if !!info.IsDir() {
		return 0, errs.Format("Irregular file: " + filename)
	}
	return info.ModTime().UnixNano(), nil
}

func JoinAbs(elem ...string) (string, error) {
	return filepath.Abs(path.Join(elem...))
}

func FindFile(searchDirs []string, filename string) (filePath string, found bool) {
	for _, searchDir := range searchDirs {
		filePath = path.Join(searchDir, filename)
		if FileExists(filePath) {
			return filePath, true
		}
	}
	return "", false
}

func FindFile2(baseDirs []string, searchDirs []string, filename string) (filePath string, found bool) {
	for _, baseDir := range baseDirs {
		for _, searchDir := range searchDirs {
			filePath = path.Join(baseDir, searchDir, filename)
			if FileExists(filePath) {
				return filePath, true
			}
		}
	}
	return "", false
}

func FindFile2ModifiedTime(baseDirs []string, searchDirs []string, filename string) (filePath string, found bool, modifiedTime int64) {
	for _, baseDir := range baseDirs {
		for _, searchDir := range searchDirs {
			filePath = path.Join(baseDir, searchDir, filename)
			if time, err := FileModifiedTime(filePath); err == nil {
				return filePath, true, time
			}
		}
	}
	return "", false, 0
}

func CombineDirs(baseDirs []string, searchDirs []string) []string {
	combinedDirs := make([]string, len(baseDirs)*len(searchDirs))
	for i, baseDir := range baseDirs {
		for j, searchDir := range searchDirs {
			combinedDirs[i*len(searchDirs)+j] = path.Join(baseDir, searchDir)
		}
	}
	return combinedDirs
}
