package web

import (
    "net/http"
    "log"
)

const (
    username = "abc@abc.com"
    password = "abc"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

    if r.Method == "GET" {
        http.ServeFile(w,r, "web/views/login.html")
    } else if r.Method == "POST" {
        err := r.ParseForm()
        if err != nil {
            log.Println("Failed to parse form")
        }
        uname := r.Form["email"][0]
        pswd := r.Form["password"][0]
        if (uname == username) && (pswd == password) {
            http.Redirect(w,r, "/front/", http.StatusFound)
            return
        }
        http.Redirect(w,r, "/", http.StatusFound)
    }
}


func FrontPageHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "web/views/layouts/base.html")
}




