package endpoint

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ChromiumHistory(paths string) []string {
	txt_historys := []string{}
	db_history, err := sql.Open("sqlite3", paths)
	if err == nil {
		defer db_history.Close()
		rows_history, err := db_history.Query("select url from urls")
		if err == nil {
			defer rows_history.Close()
			his := []history{}

			if err == nil {

				for rows_history.Next() {
					p := history{}
					err := rows_history.Scan(&p.url)
					if err != nil {
						geterr(err)
						continue
					}
					his = append(his, p)
				}

				for _, p := range his {
					txt_historys = append(txt_historys, string(p.url+"\n"))
				}
			}

		}
	}

	return txt_historys
}

type history struct {
	url string
}
