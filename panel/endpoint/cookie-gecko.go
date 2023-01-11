package endpoint

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GeckoCookie(c string) []string {
	textcook := []string{}
	db_cook, err := sql.Open("sqlite3", c)
	if err == nil {
		defer db_cook.Close()
		rows_cook, err := db_cook.Query("select host, path, isSecure, expiry, name, value from moz_cookies")
		if err == nil {
			defer rows_cook.Close()
			cookies := []cookie{}

			for rows_cook.Next() {
				p := cookie{}
				err := rows_cook.Scan(&p.host, &p.path, &p.isSecure, &p.expiry, &p.name, &p.value)
				if err == nil {
					cookies = append(cookies, p)
				}
			}

			for _, p := range cookies {
				textcook = append(textcook, fmt.Sprintf("%s\tTRUE\t%s\t%s\t%d\t%s\t%s\n", p.host, p.path, p.isSecure,
					p.expiry, p.name, p.value))
			}
		}
	}

	return textcook
}

type cookie struct {
	host     string
	path     string
	isSecure string
	expiry   int64
	name     string
	value    string
}
