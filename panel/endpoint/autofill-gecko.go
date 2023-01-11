package endpoint

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GeckoAutofill(autofil string) []string {
	textautofill := []string{""}

	db_autofill, err := sql.Open("sqlite3", autofil)
	if err == nil {
		defer db_autofill.Close()

		rows_autofill, err := db_autofill.Query("select fieldname, value from moz_formhistory")
		if err == nil {
			defer rows_autofill.Close()

			products := []cookie{}

			for rows_autofill.Next() {
				p := cookie{}
				err := rows_autofill.Scan(&p.name, &p.value)
				if err == nil {
					products = append(products, p)
				}
			}
			for _, p := range products {
				textautofill = append(textautofill, fmt.Sprintf("%v\n%v\n\n", p.name, p.value))
			}
		}
	}

	return textautofill
}
