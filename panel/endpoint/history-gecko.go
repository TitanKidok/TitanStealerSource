package endpoint

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GeckoHistory(hist string) ([]string, []string) {
	texthistor := []string{}
	textbook := []string{}

	db_history_bookmarks, err := sql.Open("sqlite3", hist)
	if err == nil {
		defer db_history_bookmarks.Close()
		rows_history, err := db_history_bookmarks.Query("select url from moz_places")
		if err == nil {
			defer rows_history.Close()

			histor := []history{}
			for rows_history.Next() {
				p := history{}
				err := rows_history.Scan(&p.url)
				if err == nil {
					histor = append(histor, p)
				}
			}

			for _, p := range histor {
				texthistor = append(texthistor, fmt.Sprint((p.url + "\n")))
			}
		}

		rows_bookmarks, err := db_history_bookmarks.Query("select title from moz_bookmarks")
		if err == nil {
			defer rows_bookmarks.Close()
			bookmark := []history{}

			for rows_bookmarks.Next() {
				p := history{}
				err := rows_bookmarks.Scan(&p.url)
				if err == nil {
					bookmark = append(bookmark, p)
				}
			}

			for _, p := range bookmark {
				textbook = append(textbook, fmt.Sprint(p.url+"\n"))
			}
		}
	}

	return texthistor, textbook
}
