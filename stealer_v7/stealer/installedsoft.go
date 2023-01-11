package main

import "os"

func InstalledSoftware() (software []string) {
	if files, err := os.ReadDir("C:/Program Files (x86)/"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				software = append(software, file.Name()+"\n")
			}
		}
	}

	return
}
