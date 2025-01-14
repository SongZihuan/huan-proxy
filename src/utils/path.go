package utils

import (
	"path/filepath"
	"runtime"
	"strings"
)

func ClearFilePathAbs(pathstr string) (string, error) {
	pathstr, err := filepath.Abs(filepath.Clean(pathstr))
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		index := strings.Index(pathstr, `:\`)
		pf := strings.ToUpper(pathstr[:index])
		ph := pathstr[index:]
		pathstr = pf + ph
	}

	return pathstr, nil
}

func FilePathEqual(path1, path2 string) bool {
	path1, err := ClearFilePathAbs(path1)
	if err != nil {
		return false
	}

	path2, err = ClearFilePathAbs(path2)
	if err != nil {
		return false
	}

	return path1 == path2
}
