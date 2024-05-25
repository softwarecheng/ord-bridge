package util

import (
	"os"
	"path/filepath"
	"strings"
)

func GetDirSize(path string) (int64, error) {
	ret := int64(0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ret += info.Size()
		}
		return nil
	})
	return ret, err
}

func GetAbsolutePath(path string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(homeDir, path[2:])
	} else {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		path = filepath.Clean(absPath)
	}
	return path, nil
}

func ParseFilePath(path string) (string, string, error) {
	absPath, err := GetAbsolutePath(path)
	if err != nil {
		return "", "", err
	}
	dir, file := filepath.Split(absPath)
	return dir, file, nil
}

func ParseDirPath(path string) (string, error) {
	absPath, err := GetAbsolutePath(path)
	if err != nil {
		return "", err
	}
	if absPath[len(absPath)-1] != filepath.Separator {
		absPath += string(filepath.Separator)
	}
	return absPath, nil
}
