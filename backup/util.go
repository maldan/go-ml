package mbackup

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type FileInfo struct {
	RelativePath string    `json:"relativePath"`
	FullPath     string    `json:"fullPath"`
	Name         string    `json:"name"`
	Ext          string    `json:"ext"`
	Dir          string    `json:"dir"`
	IsDir        bool      `json:"isDir"`
	Size         int64     `json:"size"`
	Created      time.Time `json:"created"`
}

func FSListAll(from string) ([]FileInfo, error) {
	list := make([]FileInfo, 0)

	curAbsPath, _ := filepath.Abs(from)
	curAbsPath = strings.ReplaceAll(curAbsPath, "\\", "/")

	err := filepath.Walk(from,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip dir
			if info.IsDir() {
				return nil
			}

			absPath, _ := filepath.Abs(path)
			absPath = strings.ReplaceAll(absPath, "\\", "/")

			ext := strings.Split(info.Name(), ".")

			list = append(list, FileInfo{
				FullPath:     absPath,
				RelativePath: strings.Replace(absPath, curAbsPath, "", 1),
				Name:         info.Name(),
				Ext:          ext[len(ext)-1],
				Dir:          strings.ReplaceAll(filepath.Dir(absPath), "\\", "/"),
				IsDir:        info.IsDir(),
			})

			return nil
		})
	if err != nil {
		return list, err
	}

	return list, nil
}

func Includes[T comparable](slice []T, v T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == v {
			return true
		}
	}

	return false
}

func FilterBy[T any](slice []T, filter func(*T) bool) []T {
	filtered := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			filtered = append(filtered, slice[i])
		}
	}
	return filtered
}

func Exec(args ...string) (string, string, error) {
	c := exec.Command(args[0], args[1:]...)
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	c.Stdout = stdout
	c.Stderr = stderr
	err := c.Run()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}

func SafePath(p string) string {
	return filepath.Clean(p)
}

func SafeName(p string) string {
	reg := regexp.MustCompile("[^a-zA-Z\\d._-]+")
	return reg.ReplaceAllString(p, "")
}
