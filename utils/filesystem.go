package utils

import (
	"os"
	"path"
	"path/filepath"
	"time"
)

func DirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	return err == nil && info.IsDir()
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func FileModifiedTime(filename string) (time.Time, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
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

func FindFile2ModifiedTime(baseDirs []string, searchDirs []string, filename string) (filePath string, found bool, modifiedTime time.Time) {
	for _, baseDir := range baseDirs {
		for _, searchDir := range searchDirs {
			filePath = path.Join(baseDir, searchDir, filename)
			if modifiedTime, err := FileModifiedTime(filePath); err == nil {
				return filePath, true, modifiedTime
			}
		}
	}
	return "", false, time.Time{}
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
