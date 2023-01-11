package main

import (
	"database/sql"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var db, _ = sql.Open("mysql", "root:@/fuck")

var store = sessions.NewCookieStore([]byte("DJKsd23dkvnDJFOFP24fsJOscgres8992mjcoQ!z!?"))

func error_page_render(w http.ResponseWriter) {
	var error_page_template = template.Must(template.ParseFiles("front/examples/error_page.html"))
	error_page_template.Execute(w, nil)
}

func home_page_render(w http.ResponseWriter) {
	var home_page_template = template.Must(template.ParseFiles("front/examples/home_page.html"))
	home_page_template.Execute(w, nil)
}

func ifpost(r *http.Request) bool {
	return r.Method == "POST"
}

func login_page_render(w http.ResponseWriter) {
	var login_page_template = template.Must(template.ParseFiles("front/examples/login_page.html"))
	login_page_template.Execute(w, nil)
}
